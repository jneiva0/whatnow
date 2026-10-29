package service_test

import (
	"context"
	"errors"
	"reflect"
	"testing"
	"time"

	"github.com/jameynakama/whatnow/internal/domain"
	"github.com/jameynakama/whatnow/internal/service"
)

func TestOrganizePreparationPreservesActionsFocusAndUndo(t *testing.T) {
	ctx := context.Background()
	svc, clock := newService(t)
	task := createTask(t, svc, "Organize room")
	initial := newBatch(t, svc, domain.ActorUser)
	done, _ := svc.CreateStep(ctx, initial, task.ID, service.StepFields{Description: "Drink water", Done: true, EstimateMinutes: ptr(1)})
	meditate, _ := svc.CreateStep(ctx, initial, task.ID, service.StepFields{Description: "Meditate", EstimateMinutes: ptr(3)})
	clean, _ := svc.CreateStep(ctx, initial, task.ID, service.StepFields{Description: "Pick up five objects", EstimateMinutes: ptr(4)})
	before, _ := svc.ListSteps(ctx, task.ID)
	state, err := svc.StartFocus(ctx, service.StartFocusFields{TaskID: task.ID, StepID: &clean.ID, RequestID: "organize-focus-test"})
	if err != nil {
		t.Fatal(err)
	}
	active := selected(t, state)
	clock.advance(20 * time.Second)
	batch := newBatch(t, svc, domain.ActorAgent)
	progress, err := svc.OrganizeSteps(ctx, batch, task.ID, []service.StepPlacement{
		{StepID: done.ID},
		{Description: "Drink water", Estimate: 1, Repeat: true},
		{Description: "Stretch briefly", Estimate: 2},
		{StepID: clean.ID},
	})
	if err != nil {
		t.Fatal(err)
	}
	if len(progress.Steps) != 5 || progress.Steps[0].ID != done.ID || !progress.Steps[0].Done ||
		progress.Steps[1].Done || progress.Steps[1].ID == done.ID || progress.Steps[4].ID != meditate.ID ||
		progress.Steps[3].ActualMS != 20000 || progress.Next == nil || *progress.Next.StepID != progress.Steps[1].ID {
		t.Fatalf("lost identity, order, repetition or time: %+v", progress)
	}
	if progress.Running == nil || *progress.Running.StepID != clean.ID || *progress.TotalEstimate.Minutes != 11 {
		t.Fatalf("organization changed focus or estimates: %+v", progress)
	}
	clock.advance(10 * time.Second)
	if err := svc.Undo(ctx, batch); err != nil {
		t.Fatal(err)
	}
	after, _ := svc.ListSteps(ctx, task.ID)
	state, _ = svc.Execution(ctx)
	if !reflect.DeepEqual(before, after) || selected(t, state).ID != active.ID || state.Progress[task.ID].ActualMS != 30000 {
		t.Fatalf("undo rewound execution or failed to restore organization: %+v %+v", after, state)
	}
}

func TestOrganizeRejectsInvalidSequencesAtomically(t *testing.T) {
	ctx := context.Background()
	svc, _ := newService(t)
	task := createTask(t, svc, "Task")
	other := createTask(t, svc, "Other")
	initial := newBatch(t, svc, domain.ActorUser)
	done, _ := svc.CreateStep(ctx, initial, task.ID, service.StepFields{Description: "Drink water", Done: true})
	foreign, _ := svc.CreateStep(ctx, initial, other.ID, service.StepFields{Description: "Foreign"})
	before, _ := svc.ListSteps(ctx, task.ID)
	for name, invalid := range map[string]service.StepPlacement{
		"invalid estimate": {Description: "Bad estimate", Estimate: -1},
		"foreign ID":       {StepID: foreign.ID},
		"implicit repeat":  {Description: " drink   water "},
		"disguised edit":   {StepID: done.ID, Description: "Do something unrelated"},
	} {
		t.Run(name, func(t *testing.T) {
			batch := newBatch(t, svc, domain.ActorAgent)
			_, err := svc.OrganizeSteps(ctx, batch, task.ID, []service.StepPlacement{{Description: "Valid first action"}, invalid})
			if !errors.Is(err, service.ErrInvalid) && !errors.Is(err, service.ErrConflict) {
				t.Fatalf("accepted invalid organization: %v", err)
			}
			after, _ := svc.ListSteps(ctx, task.ID)
			activity, _ := svc.BatchActivity(ctx, batch)
			if !reflect.DeepEqual(before, after) || len(activity) != 0 {
				t.Fatalf("partial changes survived failure: %+v %+v", after, activity)
			}
		})
	}
}
