package service

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	"github.com/jameynakama/whatnow/internal/domain"
)

// StepPlacement either reuses an action by ID or adds a new one. Reuse preserves
// description, estimate, completion and recorded time; edits use UpdateStep.
type StepPlacement struct {
	StepID      int64  `json:"step_id"`
	Description string `json:"description"`
	Estimate    int    `json:"estimate_minutes"`
	Repeat      bool   `json:"repeat"`
}

func stepKey(s string) string { return strings.ToLower(strings.Join(strings.Fields(s), " ")) }

// OrganizeSteps applies a proposed sequence atomically. Unmentioned actions stay
// after the sequence in their previous relative order. Nothing is deleted,
// reopened or completed, and the running focus is never switched. A preparation
// can therefore precede work without pretending to replace/decompose that work.
func (s *Service) OrganizeSteps(ctx context.Context, batch, taskID int64, placements []StepPlacement) (domain.TaskProgress, error) {
	if len(placements) > 60 {
		return domain.TaskProgress{}, fmt.Errorf("%w: at most 60 placements per organization", ErrInvalid)
	}
	var out domain.TaskProgress
	err := s.inTx(ctx, func(tx *sql.Tx) error {
		if _, err := batchActor(ctx, tx, batch); err != nil {
			return err
		}
		task, err := liveTask(ctx, tx, taskID)
		if err != nil {
			return err
		}
		if task.Status != domain.StatusOpen {
			return fmt.Errorf("%w: only an open task can be organized", ErrConflict)
		}
		existing, err := listSteps(ctx, tx, taskID, false)
		if err != nil {
			return err
		}
		byID := map[int64]domain.Step{}
		for _, st := range existing {
			byID[st.ID] = st
		}
		used := map[int64]bool{}
		ids := make([]int64, 0, len(existing)+len(placements))
		steps := append([]domain.Step{}, existing...)
		for _, p := range placements {
			if p.StepID != 0 {
				st, ok := byID[p.StepID]
				if !ok || used[p.StepID] {
					return fmt.Errorf("%w: duplicate, removed or foreign step #%d; read current steps", ErrConflict, p.StepID)
				}
				// Empty/zero mean reuse. Reject disguised replacements even for pending
				// actions: inserting preparation is different from editing their meaning.
				if (strings.TrimSpace(p.Description) != "" && p.Description != st.Description) ||
					(p.Estimate != 0 && (st.EstimateMinutes == nil || p.Estimate != *st.EstimateMinutes)) || p.Repeat {
					return fmt.Errorf("%w: reuse step #%d unchanged; use update_step for an explicit edit", ErrInvalid, p.StepID)
				}
				used[p.StepID] = true
				ids = append(ids, p.StepID)
				continue
			}
			var estimate *int
			if p.Estimate != 0 {
				estimate = &p.Estimate
			}
			f := StepFields{Description: p.Description, EstimateMinutes: estimate}
			if err := f.normalize(); err != nil {
				return err
			}
			for _, st := range steps {
				if stepKey(st.Description) == stepKey(f.Description) && !p.Repeat {
					return fmt.Errorf("%w: action already exists as step #%d (done=%t); reuse its ID, or set repeat only for a deliberate new occurrence", ErrConflict, st.ID, st.Done)
				}
			}
			st, err := insertStep(ctx, tx, taskID, len(steps), f)
			if err != nil {
				return err
			}
			if err := record(ctx, tx, batch, actionCreate, "step", st.ID, nil, st); err != nil {
				return err
			}
			steps = append(steps, st)
			ids = append(ids, st.ID)
		}
		for _, st := range existing {
			if !used[st.ID] {
				ids = append(ids, st.ID)
			}
		}
		if err := reorderSteps(ctx, tx, batch, steps, ids); err != nil {
			return err
		}
		state, err := executionState(ctx, tx, s.now())
		out = state.Progress[taskID]
		return err
	})
	return out, err
}
