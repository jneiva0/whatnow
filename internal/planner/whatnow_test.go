package planner_test

import (
	"context"
	"errors"
	"reflect"
	"strings"
	"testing"

	"github.com/jameynakama/whatnow/internal/domain"
	"github.com/jameynakama/whatnow/internal/planner"
	"github.com/jameynakama/whatnow/internal/service"
)

func pick(id int64, reason string) map[string]any {
	return map[string]any{"reason": reason, "task_id": id}
}

func TestWhatNowPicksFromTheModel(t *testing.T) {
	f := newFixture(t)
	a := f.task(t, service.TaskFields{Title: "a"})
	f.model.replies = []any{pick(a.ID, "Ten minutes before standup: perfect for this.")}
	f.start(t, 2, "")

	got, err := f.planner.WhatNow(context.Background(), 0)
	if err != nil {
		t.Fatal(err)
	}
	if got != (planner.Pick{TaskID: a.ID, Reason: "Ten minutes before standup: perfect for this."}) {
		t.Errorf("pick = %+v", got)
	}
	req := f.model.sent[0]
	if req.User != "What should I do right now?" || !req.Fallback || req.Model != "claude-sonnet-5-5" {
		t.Errorf("request: %+v", req)
	}
}

func TestNextActionContextAppliesOnlyToThisRequest(t *testing.T) {
	f := newFixture(t)
	task := f.task(t, service.TaskFields{Title: "Prepare presentation"})
	f.model.replies = []any{pick(task.ID, "Start here"), pick(task.ID, "Continue")}
	note := "I am stuck opening the document; help me find the first movement."
	if _, err := f.planner.NextAction(context.Background(), planner.ActionRequest{Context: note}); err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(f.model.sent[0].User, note) {
		t.Fatalf("planner lost current need: %q", f.model.sent[0].User)
	}
	if _, err := f.planner.WhatNow(context.Background(), 0); err != nil {
		t.Fatal(err)
	}
	if strings.Contains(f.model.sent[1].User, note) {
		t.Fatalf("a situational request leaked into a later pick: %q", f.model.sent[1].User)
	}
}

func TestWhatNowWorksBeforeTheDayStarts(t *testing.T) {
	f := newFixture(t)
	a := f.task(t, service.TaskFields{Title: "a"})
	f.model.replies = []any{pick(a.ID, "go")}

	if got, err := f.planner.WhatNow(context.Background(), 0); err != nil || got.TaskID != a.ID {
		t.Errorf("pick %+v, err %v", got, err)
	}
	if strings.Contains(system(f.model.sent[0]), "Pace today") {
		t.Error("no day, so no pace should be claimed")
	}
}

func TestWhatNowRetriesABlockedOrExcludedPick(t *testing.T) {
	f := newFixture(t)
	blocked := f.task(t, service.TaskFields{Title: "blocked"})
	first := f.task(t, service.TaskFields{Title: "first"})
	passed := f.task(t, service.TaskFields{Title: "passed on"})
	f.block(t, blocked.ID, first.ID)

	cases := []struct {
		bad     int64
		problem string
	}{{blocked.ID, "#1 is blocked"}, {passed.ID, "#3 is the task the user just passed on"}, {42, "#42 is not an eligible task"}}
	for _, c := range cases {
		f.model.sent = nil
		f.model.replies = []any{pick(c.bad, "x"), pick(first.ID, "fine")}
		got, err := f.planner.WhatNow(context.Background(), passed.ID)
		if err != nil {
			t.Fatal(err)
		}
		if got.TaskID != first.ID || got.Fallback {
			t.Errorf("pick = %+v", got)
		}
		if !strings.Contains(f.model.sent[0].User, "Not #3") || !strings.Contains(f.model.sent[1].User, c.problem) {
			t.Errorf("requests: %q then %q", f.model.sent[0].User, f.model.sent[1].User)
		}
	}
}

func TestWhatNowFallsBackToTheFirstFreeTaskInFallbackOrder(t *testing.T) {
	f := newFixture(t, errors.New("rate limited"))
	f.task(t, service.TaskFields{Title: "undated"})
	dueToday := f.task(t, service.TaskFields{Title: "due today", DueOn: date("2026-10-05")})
	overdue := f.task(t, service.TaskFields{Title: "overdue", DueOn: date("2026-10-01")})
	research := f.task(t, service.TaskFields{Title: "research", DueOn: date("2026-10-07")})
	f.block(t, overdue.ID, research.ID)

	// Overdue is first but blocked; due today is the user's pass; research is next by due date.
	got, err := f.planner.WhatNow(context.Background(), dueToday.ID)
	if err != nil {
		t.Fatal(err)
	}
	if got.TaskID != research.ID || !got.Fallback {
		t.Errorf("pick = %+v", got)
	}
}

func TestWhatNowWithNothingFree(t *testing.T) {
	f := newFixture(t)
	a := f.task(t, service.TaskFields{Title: "a"})
	f.task(t, service.TaskFields{Title: "later", NotBefore: date("2026-10-09")})
	if _, err := f.planner.WhatNow(context.Background(), a.ID); !errors.Is(err, planner.ErrNothingToDo) {
		t.Errorf("err = %v, want ErrNothingToDo", err)
	}
	if len(f.model.sent) != 0 {
		t.Error("the model was asked with nothing to pick")
	}
}

func TestWhatNowNextActionAndReplanPreserveExecution(t *testing.T) {
	f := newFixture(t)
	ctx := context.Background()
	task := f.task(t, service.TaskFields{Title: "Write"})
	f.configured = false
	batch, err := f.svc.NewBatch(ctx, domain.ActorUser)
	if err != nil {
		t.Fatal(err)
	}
	done, err := f.svc.CreateStep(ctx, batch, task.ID, service.StepFields{Description: "Open", Done: true})
	if err != nil {
		t.Fatal(err)
	}
	next, err := f.svc.CreateStep(ctx, batch, task.ID, service.StepFields{Description: "Write a sentence"})
	if err != nil {
		t.Fatal(err)
	}
	pick, err := f.planner.WhatNow(ctx, 0)
	if err != nil || pick.StepID == nil || *pick.StepID != next.ID || !pick.Fallback {
		t.Fatalf("pick=%+v err=%v", pick, err)
	}
	focus, err := f.svc.StartFocus(ctx, service.StartFocusFields{TaskID: task.ID, StepID: &next.ID, RequestID: "plan-focus-start"})
	if err != nil {
		t.Fatal(err)
	}
	if _, err := f.svc.CorrectFocus(ctx, *focus.SelectedID, 70000, "Already worked", "plan-correct-time", domain.ActorUser); err != nil {
		t.Fatal(err)
	}
	f.start(t, 3, "")
	if _, err := f.planner.Replan(ctx, domain.TriggerUser, "Keep progress"); err != nil {
		t.Fatal(err)
	}
	after, err := f.svc.Execution(ctx)
	if err != nil {
		t.Fatal(err)
	}
	if *after.SelectedID != *focus.SelectedID || after.Sessions[0].AccumulatedMS != 70000 || after.Sessions[0].State != "running" || len(after.Steps) != 2 || after.Steps[0].ID != done.ID || !after.Steps[0].Done {
		t.Fatalf("replan changed execution: %+v", after)
	}
}

func TestWhatNowCanChoosePendingStepAndRejectCompletedOrForeign(t *testing.T) {
	f := newFixture(t)
	ctx := context.Background()
	task := f.task(t, service.TaskFields{Title: "Prepare to work"})
	other := f.task(t, service.TaskFields{Title: "Other task"})
	batch, _ := f.svc.NewBatch(ctx, domain.ActorUser)
	done, err := f.svc.CreateStep(ctx, batch, task.ID, service.StepFields{Description: "Already done", Done: true})
	if err != nil {
		t.Fatal(err)
	}
	first, err := f.svc.CreateStep(ctx, batch, task.ID, service.StepFields{Description: "First"})
	if err != nil {
		t.Fatal(err)
	}
	chosen, err := f.svc.CreateStep(ctx, batch, task.ID, service.StepFields{Description: "Get water"})
	if err != nil {
		t.Fatal(err)
	}
	foreign, err := f.svc.CreateStep(ctx, batch, other.ID, service.StepFields{Description: "Unrelated"})
	if err != nil {
		t.Fatal(err)
	}
	for _, invalid := range []int64{done.ID, foreign.ID, 9999} {
		f.model.replies = []any{map[string]any{"task_id": task.ID, "step_id": invalid, "reason": "bad"}, map[string]any{"task_id": task.ID, "step_id": chosen.ID, "reason": "An easy start"}}
		got, err := f.planner.WhatNow(ctx, 0)
		if err != nil || got.StepID == nil || *got.StepID != chosen.ID || got.Fallback {
			t.Fatalf("got %+v err %v (first %d)", got, err, first.ID)
		}
	}
}

// A suggestion is a read even with absent metadata and a situational rejection.
func TestNextActionWithUnknownMetadataIsReadOnly(t *testing.T) {
	f := newFixture(t)
	ctx := context.Background()
	a := f.task(t, service.TaskFields{Title: "Drink water", RawText: "A short ritual to get started"})
	b := f.task(t, service.TaskFields{Title: "Prepare report", Notes: "Delivery matters today"})
	beforeTasks, _ := f.svc.ListTasks(ctx)
	beforeFocus, _ := f.svc.Execution(ctx)
	f.model.replies = []any{pick(a.ID, "Start with your concrete preparation"), pick(b.ID, "Try the report")}
	first, err := f.planner.NextAction(ctx, planner.ActionRequest{})
	if err != nil || first.TaskID != a.ID {
		t.Fatalf("pick=%+v err=%v", first, err)
	}
	note := "I am tired right now; another action please"
	second, err := f.planner.NextAction(ctx, planner.ActionRequest{ExcludeTaskID: a.ID, Context: note})
	if err != nil || second.TaskID != b.ID {
		t.Fatalf("pick=%+v err=%v", second, err)
	}
	if !strings.Contains(system(f.model.sent[0]), a.RawText) || !strings.Contains(f.model.sent[1].User, note) {
		t.Fatal("lost original capture or current request")
	}
	afterTasks, _ := f.svc.ListTasks(ctx)
	afterFocus, _ := f.svc.Execution(ctx)
	if !reflect.DeepEqual(beforeTasks, afterTasks) || !reflect.DeepEqual(beforeFocus.Steps, afterFocus.Steps) || !reflect.DeepEqual(beforeFocus.Sessions, afterFocus.Sessions) {
		t.Fatal("suggestion mutated task or execution")
	}
}
