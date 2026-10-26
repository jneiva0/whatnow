package service

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/jameynakama/whatnow/internal/domain"
)

const focusColumns = `id,task_id,step_id,state,segment_started_at,accumulated_ms,created_at,ended_at`

func scanFocus(row scanner) (domain.FocusSession, error) {
	var f domain.FocusSession
	var started, ended sql.NullString
	var created string
	err := row.Scan(&f.ID, &f.TaskID, &f.StepID, &f.State, &started, &f.AccumulatedMS, &created, &ended)
	if err != nil {
		return f, err
	}
	if f.CreatedAt, err = parseTime(created); err != nil {
		return f, err
	}
	if f.SegmentStartedAt, err = parseNullTime(started); err != nil {
		return f, err
	}
	f.EndedAt, err = parseNullTime(ended)
	return f, err
}
func getFocus(ctx context.Context, q querier, id int64) (domain.FocusSession, error) {
	f, err := scanFocus(q.QueryRowContext(ctx, `SELECT `+focusColumns+` FROM focus_sessions WHERE id=?`, id))
	if errors.Is(err, sql.ErrNoRows) {
		return f, fmt.Errorf("%w: focus %d", ErrNotFound, id)
	}
	return f, err
}
func saveFocus(ctx context.Context, tx *sql.Tx, f domain.FocusSession) error {
	_, err := tx.ExecContext(ctx, `UPDATE focus_sessions SET state=?,segment_started_at=?,accumulated_ms=?,ended_at=? WHERE id=?`, f.State, formatNullTime(f.SegmentStartedAt), f.AccumulatedMS, formatNullTime(f.EndedAt), f.ID)
	return err
}
func settleFocus(f *domain.FocusSession, state string, now time.Time) {
	if f.State == "running" && f.SegmentStartedAt != nil {
		ms := now.Sub(*f.SegmentStartedAt).Milliseconds()
		if ms > 0 {
			f.AccumulatedMS += ms
		}
	}
	f.State = state
	f.SegmentStartedAt = nil
	if state == "ended" {
		f.EndedAt = &now
	}
}
func pauseRunning(ctx context.Context, tx *sql.Tx, now time.Time) error {
	f, err := scanFocus(tx.QueryRowContext(ctx, `SELECT `+focusColumns+` FROM focus_sessions WHERE state='running'`))
	if errors.Is(err, sql.ErrNoRows) {
		return nil
	}
	if err != nil {
		return err
	}
	settleFocus(&f, "paused", now)
	return saveFocus(ctx, tx, f)
}
func selectFocus(ctx context.Context, tx *sql.Tx, id int64) error {
	_, err := tx.ExecContext(ctx, `INSERT INTO focus_selection(id,session_id) VALUES(1,?) ON CONFLICT(id) DO UPDATE SET session_id=excluded.session_id`, id)
	return err
}

// stopFocusTargets stops live/paused sessions for an unavailable target, preserving all time.
func stopFocusTargets(ctx context.Context, tx *sql.Tx, taskID int64, stepID *int64, now time.Time) error {
	query := `SELECT ` + focusColumns + ` FROM focus_sessions WHERE state!='ended' AND task_id=?`
	var arg any = taskID
	if stepID != nil {
		query = `SELECT ` + focusColumns + ` FROM focus_sessions WHERE state!='ended' AND step_id=?`
		arg = *stepID
	}
	rows, err := tx.QueryContext(ctx, query, arg)
	if err != nil {
		return err
	}
	list := []domain.FocusSession{}
	for rows.Next() {
		f, err := scanFocus(rows)
		if err != nil {
			rows.Close()
			return err
		}
		list = append(list, f)
	}
	rows.Close()
	if err := rows.Err(); err != nil {
		return err
	}
	for _, f := range list {
		settleFocus(&f, "ended", now)
		if err := saveFocus(ctx, tx, f); err != nil {
			return err
		}
	}
	return nil
}
func eligibleTarget(ctx context.Context, q querier, taskID int64, stepID *int64, now time.Time) error {
	t, err := liveTask(ctx, q, taskID)
	if err != nil {
		return err
	}
	if t.Status != domain.StatusOpen {
		return fmt.Errorf("%w: task is not open", ErrConflict)
	}
	if t.NotBefore != nil && string(*t.NotBefore) > now.Format(time.DateOnly) {
		return fmt.Errorf("%w: task is put off until %s", ErrConflict, *t.NotBefore)
	}
	var blockers int
	err = q.QueryRowContext(ctx, `SELECT COUNT(*) FROM blockers b JOIN tasks t ON t.id=b.blocked_by_id WHERE b.task_id=? AND t.deleted_at IS NULL AND t.status='open'`, taskID).Scan(&blockers)
	if err != nil {
		return err
	}
	if blockers > 0 {
		return fmt.Errorf("%w: task is blocked", ErrConflict)
	}
	if stepID != nil {
		st, err := liveStep(ctx, q, taskID, *stepID)
		if err != nil {
			return err
		}
		if st.Done {
			return fmt.Errorf("%w: step is already complete", ErrConflict)
		}
	}
	return nil
}
func executionState(ctx context.Context, tx *sql.Tx, now time.Time) (domain.ExecutionState, error) {
	out := domain.ExecutionState{ServerNow: now.UTC(), Sessions: []domain.FocusSession{}}
	var err error
	out.Steps, err = listSteps(ctx, tx, 0, true)
	if err != nil {
		return out, err
	}
	rows, err := tx.QueryContext(ctx, `SELECT `+focusColumns+` FROM focus_sessions ORDER BY id DESC`)
	if err != nil {
		return out, err
	}
	for rows.Next() {
		f, err := scanFocus(rows)
		if err != nil {
			rows.Close()
			return out, err
		}
		out.Sessions = append(out.Sessions, f)
	}
	rows.Close()
	if err := rows.Err(); err != nil {
		return out, err
	}
	out.Corrections = []domain.FocusCorrection{}
	corrections, err := tx.QueryContext(ctx, `SELECT id,session_id,before_ms,after_ms,reason,actor,created_at FROM focus_corrections ORDER BY id`)
	if err != nil {
		return out, err
	}
	for corrections.Next() {
		var c domain.FocusCorrection
		var created string
		if err := corrections.Scan(&c.ID, &c.SessionID, &c.BeforeMS, &c.AfterMS, &c.Reason, &c.Actor, &created); err != nil {
			corrections.Close()
			return out, err
		}
		c.CreatedAt, err = parseTime(created)
		if err != nil {
			corrections.Close()
			return out, err
		}
		out.Corrections = append(out.Corrections, c)
	}
	corrections.Close()
	if err := corrections.Err(); err != nil {
		return out, err
	}
	err = tx.QueryRowContext(ctx, `SELECT session_id FROM focus_selection WHERE id=1`).Scan(&out.SelectedID)
	if errors.Is(err, sql.ErrNoRows) {
		err = nil
	}
	return out, err
}
func (s *Service) Execution(ctx context.Context) (domain.ExecutionState, error) {
	var out domain.ExecutionState
	err := s.inTx(ctx, func(tx *sql.Tx) error { var err error; out, err = executionState(ctx, tx, s.now()); return err })
	return out, err
}

type StartFocusFields struct {
	TaskID    int64  `json:"task_id"`
	StepID    *int64 `json:"step_id"`
	RequestID string `json:"request_id"`
}

func sameTarget(f domain.FocusSession, in StartFocusFields) bool {
	return f.TaskID == in.TaskID && ((f.StepID == nil && in.StepID == nil) || (f.StepID != nil && in.StepID != nil && *f.StepID == *in.StepID))
}

// StartFocus serializes switching and persists request tokens. A retry never resurrects a paused
// or ended session, and two clients cannot leave two running timers.
func (s *Service) StartFocus(ctx context.Context, in StartFocusFields) (domain.ExecutionState, error) {
	if len(in.RequestID) < 8 || len(in.RequestID) > 200 {
		return domain.ExecutionState{}, fmt.Errorf("%w: request_id (8..200 characters) is required", ErrInvalid)
	}
	var out domain.ExecutionState
	err := s.inTx(ctx, func(tx *sql.Tx) error {
		now := s.now()
		var existing int64
		err := tx.QueryRowContext(ctx, `SELECT session_id FROM focus_requests WHERE request_id=?`, in.RequestID).Scan(&existing)
		if err == nil {
			f, err := getFocus(ctx, tx, existing)
			if err != nil {
				return err
			}
			if !sameTarget(f, in) {
				return fmt.Errorf("%w: request_id already used for another focus", ErrConflict)
			}
			out, err = executionState(ctx, tx, now)
			return err
		}
		if !errors.Is(err, sql.ErrNoRows) {
			return err
		}
		if err := eligibleTarget(ctx, tx, in.TaskID, in.StepID, now.In(s.loc)); err != nil {
			return err
		}
		f, err := scanFocus(tx.QueryRowContext(ctx, `SELECT `+focusColumns+` FROM focus_sessions WHERE state='running'`))
		if err != nil && !errors.Is(err, sql.ErrNoRows) {
			return err
		}
		if err != nil || !sameTarget(f, in) {
			if err := pauseRunning(ctx, tx, now); err != nil {
				return err
			}
			res, err := tx.ExecContext(ctx, `INSERT INTO focus_sessions(task_id,step_id,state,segment_started_at,created_at) VALUES(?,?,'running',?,?)`, in.TaskID, in.StepID, formatTime(now), formatTime(now))
			if err != nil {
				return err
			}
			f.ID, err = res.LastInsertId()
			if err != nil {
				return err
			}
		}
		if err := selectFocus(ctx, tx, f.ID); err != nil {
			return err
		}
		if _, err := tx.ExecContext(ctx, `INSERT INTO focus_requests VALUES(?,?)`, in.RequestID, f.ID); err != nil {
			return err
		}
		out, err = executionState(ctx, tx, now)
		return err
	})
	return out, err
}
func (s *Service) ChangeFocus(ctx context.Context, id int64, action string) (domain.ExecutionState, error) {
	var out domain.ExecutionState
	err := s.inTx(ctx, func(tx *sql.Tx) error {
		f, err := getFocus(ctx, tx, id)
		if err != nil {
			return err
		}
		now := s.now()
		switch action {
		case "pause":
			if f.State == "running" {
				settleFocus(&f, "paused", now)
			}
		case "end":
			if f.State != "ended" {
				settleFocus(&f, "ended", now)
			}
		case "resume":
			if f.State == "ended" {
				return fmt.Errorf("%w: focus has ended; start a new session", ErrConflict)
			}
			if err := eligibleTarget(ctx, tx, f.TaskID, f.StepID, now.In(s.loc)); err != nil {
				return err
			}
			if f.State != "running" {
				if err := pauseRunning(ctx, tx, now); err != nil {
					return err
				}
				f.State = "running"
				f.SegmentStartedAt = &now
			}
			if err := selectFocus(ctx, tx, id); err != nil {
				return err
			}
		default:
			return fmt.Errorf("%w: unknown focus action", ErrInvalid)
		}
		if err := saveFocus(ctx, tx, f); err != nil {
			return err
		}
		out, err = executionState(ctx, tx, now)
		return err
	})
	return out, err
}

// CompleteFocus records only the explicit target's status in undo. Reverting that status never
// rewinds or resumes the elapsed-time ledger.
func (s *Service) CompleteFocus(ctx context.Context, batch, id int64) (domain.ExecutionState, error) {
	var out domain.ExecutionState
	err := s.inTx(ctx, func(tx *sql.Tx) error {
		if _, err := batchActor(ctx, tx, batch); err != nil {
			return err
		}
		f, err := getFocus(ctx, tx, id)
		if err != nil {
			return err
		}
		now := s.now()
		// Repeated completion is a no-op, including after an undo: an ended session cannot complete twice.
		if f.State == "ended" {
			out, err = executionState(ctx, tx, now)
			return err
		}
		if err := eligibleTarget(ctx, tx, f.TaskID, f.StepID, now.In(s.loc)); err != nil {
			return err
		}
		if f.StepID != nil {
			before, err := liveStep(ctx, tx, f.TaskID, *f.StepID)
			if err != nil {
				return err
			}
			st := before
			st.Done = true
			if err := writeStep(ctx, tx, st); err != nil {
				return err
			}
			if err := record(ctx, tx, batch, actionStatus, "step", st.ID, before, st); err != nil {
				return err
			}
			if err := stopStepFocus(ctx, tx, st.ID, now); err != nil {
				return err
			}
		} else {
			before, err := liveTask(ctx, tx, f.TaskID)
			if err != nil {
				return err
			}
			t := before
			t.Status = domain.StatusDone
			t.UpdatedAt = now
			if err := writeTask(ctx, tx, t); err != nil {
				return err
			}
			if err := record(ctx, tx, batch, actionStatus, entityTask, t.ID, before, t); err != nil {
				return err
			}
			if err := stopFocusTargets(ctx, tx, t.ID, nil, now); err != nil {
				return err
			}
		}
		out, err = executionState(ctx, tx, now)
		return err
	})
	return out, err
}

// CorrectFocus deliberately replaces the elapsed total at the server's current instant. The
// segment anchor moves with it if running, so elapsed time before correction is never counted twice.
func (s *Service) CorrectFocus(ctx context.Context, id, totalMS int64, reason, requestID string, actor domain.Actor) (domain.ExecutionState, error) {
	reason = strings.TrimSpace(reason)
	if totalMS < 0 || totalMS > 365*24*60*60*1000 || reason == "" || len([]rune(reason)) > 2000 || !actor.Valid() || len(requestID) < 8 || len(requestID) > 200 {
		return domain.ExecutionState{}, fmt.Errorf("%w: total_ms must be 0..31536000000, with a reason and request_id (8..200 chars)", ErrInvalid)
	}
	var out domain.ExecutionState
	err := s.inTx(ctx, func(tx *sql.Tx) error {
		f, err := getFocus(ctx, tx, id)
		if err != nil {
			return err
		}
		now := s.now()
		var previousID, previousMS int64
		var previousReason string
		err = tx.QueryRowContext(ctx, `SELECT session_id,after_ms,reason FROM focus_corrections WHERE request_id=?`, requestID).Scan(&previousID, &previousMS, &previousReason)
		if err == nil {
			if previousID != id || previousMS != totalMS || previousReason != reason {
				return fmt.Errorf("%w: correction request_id already used", ErrConflict)
			}
			out, err = executionState(ctx, tx, now)
			return err
		}
		if !errors.Is(err, sql.ErrNoRows) {
			return err
		}
		before := f.AccumulatedMS
		if f.State == "running" && f.SegmentStartedAt != nil {
			before += max(int64(0), now.Sub(*f.SegmentStartedAt).Milliseconds())
			f.SegmentStartedAt = &now
		}
		f.AccumulatedMS = totalMS
		if err := saveFocus(ctx, tx, f); err != nil {
			return err
		}
		_, err = tx.ExecContext(ctx, `INSERT INTO focus_corrections(session_id,before_ms,after_ms,reason,actor,request_id,created_at) VALUES(?,?,?,?,?,?,?)`, id, before, totalMS, reason, actor, requestID, formatTime(now))
		if err != nil {
			return err
		}
		out, err = executionState(ctx, tx, now)
		return err
	})
	return out, err
}
