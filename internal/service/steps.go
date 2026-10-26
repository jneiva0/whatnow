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

type StepFields struct {
	Description     string `json:"description"`
	EstimateMinutes *int   `json:"estimate_minutes"`
	Done            bool   `json:"done"`
}

func validEstimate(v *int) error {
	if v != nil && (*v < 1 || *v > 1440) {
		return fmt.Errorf("%w: estimate_minutes must be between 1 and 1440", ErrInvalid)
	}
	return nil
}
func (f *StepFields) normalize() error {
	f.Description = strings.TrimSpace(f.Description)
	if f.Description == "" || len([]rune(f.Description)) > 2000 {
		return fmt.Errorf("%w: description must contain 1 to 2000 characters", ErrInvalid)
	}
	return validEstimate(f.EstimateMinutes)
}

const stepColumns = `id, task_id, description, estimate_minutes, position, done, deleted_at`

func scanStep(row scanner) (domain.Step, error) {
	var st domain.Step
	var deleted sql.NullString
	err := row.Scan(&st.ID, &st.TaskID, &st.Description, &st.EstimateMinutes, &st.Position, &st.Done, &deleted)
	if err != nil {
		return st, err
	}
	st.DeletedAt, err = parseNullTime(deleted)
	return st, err
}
func liveStep(ctx context.Context, q querier, taskID, id int64) (domain.Step, error) {
	st, err := scanStep(q.QueryRowContext(ctx, `SELECT `+stepColumns+` FROM steps WHERE id = ? AND task_id = ? AND deleted_at IS NULL`, id, taskID))
	if errors.Is(err, sql.ErrNoRows) {
		return st, fmt.Errorf("%w: step %d of task %d", ErrNotFound, id, taskID)
	}
	return st, err
}
func writeStep(ctx context.Context, tx *sql.Tx, st domain.Step) error {
	_, err := tx.ExecContext(ctx, `UPDATE steps SET description=?, estimate_minutes=?, position=?, done=?, deleted_at=? WHERE id=?`, st.Description, st.EstimateMinutes, st.Position, st.Done, formatNullTime(st.DeletedAt), st.ID)
	return err
}
func insertStep(ctx context.Context, tx *sql.Tx, taskID int64, position int, f StepFields) (domain.Step, error) {
	st := domain.Step{TaskID: taskID, Description: f.Description, EstimateMinutes: f.EstimateMinutes, Position: position, Done: f.Done}
	res, err := tx.ExecContext(ctx, `INSERT INTO steps (task_id,description,estimate_minutes,position,done) VALUES (?,?,?,?,?)`, taskID, f.Description, f.EstimateMinutes, position, f.Done)
	if err == nil {
		st.ID, err = res.LastInsertId()
	}
	return st, err
}
func listSteps(ctx context.Context, q querier, taskID int64, includeDeleted bool) ([]domain.Step, error) {
	query := `SELECT ` + stepColumns + ` FROM steps WHERE (?=0 OR task_id=?)`
	if !includeDeleted {
		query += ` AND deleted_at IS NULL`
	}
	query += ` ORDER BY task_id,position,id`
	rows, err := q.QueryContext(ctx, query, taskID, taskID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	out := []domain.Step{}
	for rows.Next() {
		st, err := scanStep(rows)
		if err != nil {
			return nil, err
		}
		out = append(out, st)
	}
	return out, rows.Err()
}
func (s *Service) ListSteps(ctx context.Context, taskID int64) ([]domain.Step, error) {
	return listSteps(ctx, s.db, taskID, false)
}
func (s *Service) CreateStep(ctx context.Context, batch, taskID int64, f StepFields) (domain.Step, error) {
	if err := f.normalize(); err != nil {
		return domain.Step{}, err
	}
	var st domain.Step
	err := s.inTx(ctx, func(tx *sql.Tx) error {
		if _, err := batchActor(ctx, tx, batch); err != nil {
			return err
		}
		if _, err := liveTask(ctx, tx, taskID); err != nil {
			return err
		}
		var pos int
		if err := tx.QueryRowContext(ctx, `SELECT COALESCE(MAX(position),-1)+1 FROM steps WHERE task_id=? AND deleted_at IS NULL`, taskID).Scan(&pos); err != nil {
			return err
		}
		var err error
		st, err = insertStep(ctx, tx, taskID, pos, f)
		if err != nil {
			return err
		}
		return record(ctx, tx, batch, actionCreate, "step", st.ID, nil, st)
	})
	return st, err
}
func (s *Service) UpdateStep(ctx context.Context, batch, taskID, id int64, f StepFields) (domain.Step, error) {
	if err := f.normalize(); err != nil {
		return domain.Step{}, err
	}
	var st domain.Step
	err := s.inTx(ctx, func(tx *sql.Tx) error {
		if _, err := batchActor(ctx, tx, batch); err != nil {
			return err
		}
		if _, err := liveTask(ctx, tx, taskID); err != nil {
			return err
		}
		before, err := liveStep(ctx, tx, taskID, id)
		if err != nil {
			return err
		}
		st = before
		st.Description = f.Description
		st.EstimateMinutes = f.EstimateMinutes
		st.Done = f.Done
		if st.Done {
			if err := stopStepFocus(ctx, tx, id, s.now()); err != nil {
				return err
			}
		}
		if err := writeStep(ctx, tx, st); err != nil {
			return err
		}
		return record(ctx, tx, batch, actionUpdate, "step", id, before, st)
	})
	return st, err
}
func (s *Service) DeleteStep(ctx context.Context, batch, taskID, id int64) error {
	return s.inTx(ctx, func(tx *sql.Tx) error {
		if _, err := batchActor(ctx, tx, batch); err != nil {
			return err
		}
		if _, err := liveTask(ctx, tx, taskID); err != nil {
			return err
		}
		before, err := liveStep(ctx, tx, taskID, id)
		if err != nil {
			return err
		}
		st := before
		now := s.now()
		st.DeletedAt = &now
		if err := stopStepFocus(ctx, tx, id, now); err != nil {
			return err
		}
		if err := writeStep(ctx, tx, st); err != nil {
			return err
		}
		return record(ctx, tx, batch, actionDelete, "step", id, before, st)
	})
}
func reorderSteps(ctx context.Context, tx *sql.Tx, batch int64, steps []domain.Step, ids []int64) error {
	if len(ids) != len(steps) {
		return fmt.Errorf("%w: ids must contain every live step exactly once", ErrInvalid)
	}
	byID := map[int64]domain.Step{}
	for _, st := range steps {
		byID[st.ID] = st
	}
	for pos, id := range ids {
		before, ok := byID[id]
		if !ok {
			return fmt.Errorf("%w: duplicate or foreign step %d", ErrInvalid, id)
		}
		delete(byID, id)
		st := before
		st.Position = pos
		if st.Position != before.Position {
			if err := writeStep(ctx, tx, st); err != nil {
				return err
			}
			if err := record(ctx, tx, batch, actionUpdate, "step", id, before, st); err != nil {
				return err
			}
		}
	}
	return nil
}
func (s *Service) ReorderSteps(ctx context.Context, batch, taskID int64, ids []int64) ([]domain.Step, error) {
	var out []domain.Step
	err := s.inTx(ctx, func(tx *sql.Tx) error {
		if _, err := batchActor(ctx, tx, batch); err != nil {
			return err
		}
		if _, err := liveTask(ctx, tx, taskID); err != nil {
			return err
		}
		steps, err := listSteps(ctx, tx, taskID, false)
		if err != nil {
			return err
		}
		if err := reorderSteps(ctx, tx, batch, steps, ids); err != nil {
			return err
		}
		out, err = listSteps(ctx, tx, taskID, false)
		return err
	})
	return out, err
}

// SplitStep replaces only one pending action, at the same level. Completed siblings and all
// focus history survive. The old action is retained as a tombstone for its recorded sessions.
func (s *Service) SplitStep(ctx context.Context, batch, taskID, id int64, fields []StepFields) ([]domain.Step, error) {
	if len(fields) < 2 || len(fields) > 30 {
		return nil, fmt.Errorf("%w: split needs 2 to 30 actions", ErrInvalid)
	}
	for i := range fields {
		if fields[i].Done {
			return nil, fmt.Errorf("%w: new actions must be pending", ErrInvalid)
		}
		if err := fields[i].normalize(); err != nil {
			return nil, err
		}
	}
	var out []domain.Step
	err := s.inTx(ctx, func(tx *sql.Tx) error {
		if _, err := batchActor(ctx, tx, batch); err != nil {
			return err
		}
		if _, err := liveTask(ctx, tx, taskID); err != nil {
			return err
		}
		before, err := liveStep(ctx, tx, taskID, id)
		if err != nil {
			return err
		}
		if before.Done {
			return fmt.Errorf("%w: a completed action cannot be split", ErrConflict)
		}
		siblings, err := listSteps(ctx, tx, taskID, false)
		if err != nil {
			return err
		}
		now := s.now()
		st := before
		st.DeletedAt = &now
		if err := stopStepFocus(ctx, tx, id, now); err != nil {
			return err
		}
		if err := writeStep(ctx, tx, st); err != nil {
			return err
		}
		if err := record(ctx, tx, batch, actionDelete, "step", id, before, st); err != nil {
			return err
		}
		ids := []int64{}
		for _, sib := range siblings {
			if sib.ID != id {
				ids = append(ids, sib.ID)
				continue
			}
			for _, f := range fields {
				created, err := insertStep(ctx, tx, taskID, len(siblings)+len(ids), f)
				if err != nil {
					return err
				}
				if err := record(ctx, tx, batch, actionCreate, "step", created.ID, nil, created); err != nil {
					return err
				}
				ids = append(ids, created.ID)
			}
		}
		all, err := listSteps(ctx, tx, taskID, false)
		if err != nil {
			return err
		}
		if err := reorderSteps(ctx, tx, batch, all, ids); err != nil {
			return err
		}
		out, err = listSteps(ctx, tx, taskID, false)
		return err
	})
	return out, err
}
func stopStepFocus(ctx context.Context, tx *sql.Tx, id int64, now time.Time) error {
	return stopFocusTargets(ctx, tx, 0, &id, now)
}
