package facts_test

import (
	"testing"
	"time"

	"github.com/jameynakama/whatnow/internal/domain"
	"github.com/jameynakama/whatnow/internal/facts"
)

func TestExecutionSeparatesScopePendingWorkAndRecordedTime(t *testing.T) {
	now := at(5, 10, 0)
	started := now.Add(-30 * time.Second)
	removed := now.Add(-time.Hour)
	five, ten, fifty, manual := 5, 10, 50, 4
	first, active, archived := int64(1), int64(2), int64(4)
	selectedID := int64(3)
	task := domain.Task{ID: 7, Status: domain.StatusOpen, EstimateMinutes: &manual}
	steps := []domain.Step{
		{ID: first, TaskID: 7, Position: 0, Done: true, EstimateMinutes: &five},
		{ID: active, TaskID: 7, Position: 1, EstimateMinutes: &ten},
		{ID: 3, TaskID: 7, Position: 2},
		{ID: archived, TaskID: 7, EstimateMinutes: &fifty, DeletedAt: &removed},
	}
	sessions := []domain.FocusSession{
		{TaskID: 7, AccumulatedMS: 2000, State: "ended"},
		{TaskID: 7, StepID: &first, AccumulatedMS: 3000, State: "ended"},
		{ID: selectedID, TaskID: 7, StepID: &active, AccumulatedMS: 5000, State: "running", SegmentStartedAt: &started},
		{TaskID: 7, StepID: &archived, AccumulatedMS: 7000, State: "ended"},
		{TaskID: 99, AccumulatedMS: 900000, State: "ended"},
	}
	p := facts.Execution(now, []domain.Task{task}, steps, sessions, &selectedID)[7]
	if p.ActualMS != 47000 || p.DirectActualMS != 2000 || p.ArchivedActualMS != 7000 || p.Steps[1].ActualMS != 35000 {
		t.Fatalf("recorded time counted incorrectly: %+v", p)
	}
	if p.TotalEstimate.Minutes == nil || *p.TotalEstimate.Minutes != 15 || !p.TotalEstimate.Partial ||
		p.PendingEstimate.Minutes == nil || *p.PendingEstimate.Minutes != 10 || !p.PendingEstimate.Partial {
		t.Fatalf("scope/pending estimates: %+v", p)
	}
	if len(p.Steps) != 3 || p.CompletedSteps != 1 || p.Next == nil || *p.Next.StepID != active || p.Running == nil || !p.Steps[1].Running {
		t.Fatalf("action references/progress: %+v", p)
	}
	// Computing facts never edits source state or turns elapsed time into completion.
	if steps[1].Done || sessions[2].AccumulatedMS != 5000 {
		t.Fatal("facts mutated execution")
	}
	// Pausing is still a selected focus, not loss of the current action.
	sessions[2].State, sessions[2].SegmentStartedAt = "paused", nil
	p = facts.Execution(now, []domain.Task{task}, steps, sessions, &selectedID)[7]
	if p.Focus == nil || p.Focus.ID != selectedID || p.Focus.State != "paused" || p.Running != nil {
		t.Fatalf("paused focus disappeared from context: %+v", p)
	}
}

func TestExecutionUnknownPendingEstimateDiffersFromNothingPending(t *testing.T) {
	now := at(5, 10, 0)
	task := domain.Task{ID: 1, Status: domain.StatusOpen}
	step := domain.Step{ID: 2, TaskID: 1}
	p := facts.Execution(now, []domain.Task{task}, []domain.Step{step}, nil, nil)[1]
	if p.PendingEstimate.Minutes != nil || !p.PendingEstimate.Partial {
		t.Fatalf("unknown became zero: %+v", p)
	}
	step.Done = true
	p = facts.Execution(now, []domain.Task{task}, []domain.Step{step}, nil, nil)[1]
	if p.Next != nil || p.PendingEstimate.Minutes == nil || *p.PendingEstimate.Minutes != 0 || p.PendingEstimate.Partial {
		t.Fatalf("all steps done should have no pending action: %+v", p)
	}
	p = facts.Execution(now, []domain.Task{task}, nil, nil, nil)[1]
	if p.Next == nil || p.Next.StepID != nil || p.PendingEstimate.Minutes != nil {
		t.Fatalf("direct task: %+v", p)
	}
}
