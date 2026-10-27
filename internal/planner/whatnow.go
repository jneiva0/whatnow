package planner

import (
	"context"
	"errors"
	"fmt"

	"github.com/jameynakama/whatnow/internal/facts"
	js "github.com/jameynakama/whatnow/internal/jsonschema"
)

// ErrNothingToDo means no eligible task is free to start.
var ErrNothingToDo = errors.New("no task is free to start")

var whatNowSchema = js.Object(
	js.Prop("reason", js.String("One or two sentences to the user on why this task, now.")),
	js.Prop("task_id", js.Integer("The task's number, from #<id>.")),
	js.Prop("step_id", js.Integer("An existing pending step belonging to this task, or 0 to use its next pending step / execute a concrete task directly.")),
)

// Pick is the answer to "what now?". Fallback means code picked because the model couldn't.
type Pick struct {
	TaskID   int64  `json:"task_id"`
	StepID   *int64 `json:"step_id"`
	Reason   string `json:"reason"`
	Fallback bool   `json:"fallback"`
}

// WhatNow picks one task to start now from every eligible task, planned or not. exclude, when not
// zero, is a task the user just passed on; it is left out of this answer only.
func (p *Planner) WhatNow(ctx context.Context, exclude int64) (Pick, error) {
	st, err := p.gather(ctx, true)
	if err != nil {
		return Pick{}, err
	}
	var first int64
	for _, tf := range st.facts.Tasks {
		if len(tf.BlockedBy) == 0 && tf.Task.ID != exclude {
			first = tf.Task.ID
			break
		}
	}
	if first == 0 {
		return Pick{}, ErrNothingToDo
	}

	pick, err := p.askPick(ctx, st, exclude)
	if err != nil {
		return p.pickAction(ctx, Pick{TaskID: first, Reason: fallbackNote(err), Fallback: true})
	}
	return p.pickAction(ctx, pick)
}

func (p *Planner) askPick(ctx context.Context, st state, exclude int64) (Pick, error) {
	if st.model == nil {
		return Pick{}, st.modelErr
	}
	ctx, cancel := context.WithTimeout(ctx, modelTimeout)
	defer cancel()
	user := "What should I do right now?"
	if exclude != 0 {
		user += fmt.Sprintf(" Not #%d, I just passed on that one.", exclude)
	}
	return ask(ctx, st.model, st.request(whatNowInstructions, user, whatNowSchema, "pick"), func(pk Pick) error {
		if err := checkPick(pk, st.facts, exclude); err != nil {
			return err
		}
		if pk.StepID != nil && *pk.StepID != 0 {
			steps, err := p.svc.ListSteps(ctx, pk.TaskID)
			if err != nil {
				return err
			}
			for _, step := range steps {
				if step.ID == *pk.StepID && !step.Done {
					return nil
				}
			}
			return fmt.Errorf("step #%d is not a pending action of task #%d", *pk.StepID, pk.TaskID)
		}
		return nil
	})
}

func checkPick(pk Pick, f facts.Facts, exclude int64) error {
	tf, ok := f.Task(pk.TaskID)
	switch {
	case !ok:
		return fmt.Errorf("#%d is not an eligible task", pk.TaskID)
	case len(tf.BlockedBy) > 0:
		return fmt.Errorf("#%d is blocked", pk.TaskID)
	case pk.TaskID == exclude:
		return fmt.Errorf("#%d is the task the user just passed on", pk.TaskID)
	}
	return nil
}

// Honor a validated model action; otherwise use the first pending action for a reliable fallback.
func (p *Planner) pickAction(ctx context.Context, pick Pick) (Pick, error) {
	steps, err := p.svc.ListSteps(ctx, pick.TaskID)
	if err != nil {
		return pick, err
	}
	for _, st := range steps {
		if pick.StepID != nil && st.ID == *pick.StepID && !st.Done {
			return pick, nil
		}
	}
	pick.StepID = nil
	for _, st := range steps {
		if !st.Done {
			pick.StepID = &st.ID
			break
		}
	}
	return pick, nil
}
