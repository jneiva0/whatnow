package domain

import "testing"

func TestTaskEstimatePreservesExplicitAndIncompleteWork(t *testing.T) {
	small, large, fifteen := 5, 60, 15
	task := Task{ID: 1, EstimateMinutes: &small}
	steps := []Step{{TaskID: 1, EstimateMinutes: &fifteen, Done: true}, {TaskID: 1, EstimateMinutes: &fifteen}, {TaskID: 1, EstimateMinutes: &fifteen}, {TaskID: 2, EstimateMinutes: &large}}
	got := TaskEstimate(task, steps)
	if got.Minutes == nil || *got.Minutes != 45 || got.Partial {
		t.Fatalf("got %+v", got)
	}
	task.EstimateMinutes = &large
	if got := TaskEstimate(task, steps); *got.Minutes != 60 {
		t.Fatalf("explicit estimate lost: %+v", got)
	}
	steps[1].EstimateMinutes = nil
	if got := TaskEstimate(task, steps); !got.Partial || *got.Minutes != 60 {
		t.Fatalf("unknown not represented: %+v", got)
	}
	task.EstimateMinutes = nil
	if got := TaskEstimate(task, []Step{{TaskID: 1}}); got.Minutes != nil || !got.Partial {
		t.Fatalf("unknown became zero: %+v", got)
	}
}
