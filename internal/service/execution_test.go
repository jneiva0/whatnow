package service_test

import (
	"context"
	"errors"
	"fmt"
	gosync "sync"
	"testing"
	"time"

	"github.com/jameynakama/whatnow/internal/domain"
	"github.com/jameynakama/whatnow/internal/service"
)

func TestStepCRUDReorderSplitUndo(t *testing.T) {
	ctx := context.Background()
	svc, _ := newService(t)
	task := createTask(t, svc, "Ship report")
	other := createTask(t, svc, "Other")
	batch := func() int64 { return newBatch(t, svc, domain.ActorUser) }
	a, err := svc.CreateStep(ctx, batch(), task.ID, service.StepFields{Description: "Open report", EstimateMinutes: ptr(2)})
	if err != nil {
		t.Fatal(err)
	}
	b, err := svc.CreateStep(ctx, batch(), task.ID, service.StepFields{Description: "Write draft"})
	if err != nil {
		t.Fatal(err)
	}
	if _, err := svc.UpdateStep(ctx, batch(), other.ID, a.ID, service.StepFields{Description: "Foreign"}); !errors.Is(err, service.ErrNotFound) {
		t.Fatalf("foreign step: %v", err)
	}
	if _, err := svc.CreateStep(ctx, batch(), task.ID, service.StepFields{Description: " ", EstimateMinutes: ptr(-1)}); !errors.Is(err, service.ErrInvalid) {
		t.Fatal(err)
	}
	reorder := batch()
	steps, err := svc.ReorderSteps(ctx, reorder, task.ID, []int64{b.ID, a.ID})
	if err != nil || steps[0].ID != b.ID {
		t.Fatalf("reorder %v %v", steps, err)
	}
	if _, err := svc.ReorderSteps(ctx, batch(), task.ID, []int64{a.ID, a.ID}); !errors.Is(err, service.ErrInvalid) {
		t.Fatal(err)
	}
	if err := svc.Undo(ctx, reorder); err != nil {
		t.Fatal(err)
	}
	done := batch()
	if _, err := svc.UpdateStep(ctx, done, task.ID, a.ID, service.StepFields{Description: a.Description, Done: true}); err != nil {
		t.Fatal(err)
	}
	if _, err := svc.SplitStep(ctx, batch(), task.ID, a.ID, []service.StepFields{{Description: "x"}, {Description: "y"}}); !errors.Is(err, service.ErrConflict) {
		t.Fatal(err)
	}
	split := batch()
	steps, err = svc.SplitStep(ctx, split, task.ID, b.ID, []service.StepFields{{Description: "Write one sentence"}, {Description: "Read it"}})
	if err != nil {
		t.Fatal(err)
	}
	if len(steps) != 3 || !steps[0].Done || steps[1].Description != "Write one sentence" {
		t.Fatalf("lost progress/order: %+v", steps)
	}
	if err := svc.Undo(ctx, split); err != nil {
		t.Fatal(err)
	}
	steps, _ = svc.ListSteps(ctx, task.ID)
	if len(steps) != 2 || steps[1].ID != b.ID {
		t.Fatalf("split undo: %+v", steps)
	}
	if err := svc.Undo(ctx, done); err != nil {
		t.Fatal(err)
	}
	steps, _ = svc.ListSteps(ctx, task.ID)
	if steps[0].Done {
		t.Fatal("undo should reopen")
	}
	del := batch()
	if err := svc.DeleteStep(ctx, del, task.ID, a.ID); err != nil {
		t.Fatal(err)
	}
	if err := svc.Undo(ctx, del); err != nil {
		t.Fatal(err)
	}
	got, _ := svc.GetTask(ctx, task.ID)
	if got.Status != domain.StatusOpen {
		t.Fatal("steps completed parent")
	}
}
func selected(t *testing.T, st domain.ExecutionState) domain.FocusSession {
	t.Helper()
	if st.SelectedID != nil {
		for _, f := range st.Sessions {
			if f.ID == *st.SelectedID {
				return f
			}
		}
	}
	t.Fatal("no selected focus")
	return domain.FocusSession{}
}
func TestFocusClockSwitchIdempotencyAndUndo(t *testing.T) {
	ctx := context.Background()
	svc, clock := newService(t)
	task := createTask(t, svc, "Do it")
	other := createTask(t, svc, "Other")
	b := newBatch(t, svc, domain.ActorUser)
	step, err := svc.CreateStep(ctx, b, task.ID, service.StepFields{Description: "One action", EstimateMinutes: ptr(3)})
	if err != nil {
		t.Fatal(err)
	}
	in := service.StartFocusFields{TaskID: task.ID, StepID: &step.ID, RequestID: "start-first"}
	st, err := svc.StartFocus(ctx, in)
	if err != nil {
		t.Fatal(err)
	}
	first := selected(t, st)
	clock.advance(90 * time.Second)
	st, err = svc.Execution(ctx)
	if err != nil {
		t.Fatal(err)
	}
	if selected(t, st).State != "running" {
		t.Fatal("zero must not stop")
	}
	st, err = svc.StartFocus(ctx, in)
	if err != nil || len(st.Sessions) != 1 {
		t.Fatalf("retry: %v %v", st, err)
	}
	st, err = svc.StartFocus(ctx, service.StartFocusFields{TaskID: other.ID, RequestID: "start-second"})
	if err != nil {
		t.Fatal(err)
	}
	second := selected(t, st)
	if st.Sessions[1].State != "paused" || st.Sessions[1].AccumulatedMS != 90000 {
		t.Fatalf("switch lost time: %+v", st.Sessions)
	}
	st, err = svc.StartFocus(ctx, in)
	if err != nil || selected(t, st).ID != second.ID {
		t.Fatal("old retry changed selected focus", err)
	}
	clock.advance(10 * time.Second)
	st, err = svc.ChangeFocus(ctx, first.ID, "resume")
	if err != nil {
		t.Fatal(err)
	}
	if selected(t, st).AccumulatedMS != 90000 {
		t.Fatal("resume reset time")
	}
	clock.advance(15 * time.Second)
	st, err = svc.ChangeFocus(ctx, first.ID, "pause")
	if err != nil {
		t.Fatal(err)
	}
	if selected(t, st).AccumulatedMS != 105000 {
		t.Fatal("pause incorrect")
	}
	st, err = svc.ChangeFocus(ctx, first.ID, "pause")
	if err != nil || selected(t, st).AccumulatedMS != 105000 {
		t.Fatal("pause twice counted twice")
	}
	done := newBatch(t, svc, domain.ActorUser)
	st, err = svc.CompleteFocus(ctx, done, first.ID)
	if err != nil {
		t.Fatal(err)
	}
	if selected(t, st).State != "ended" {
		t.Fatal("completion did not end focus")
	}
	steps, _ := svc.ListSteps(ctx, task.ID)
	parent, _ := svc.GetTask(ctx, task.ID)
	if !steps[0].Done || parent.Status != domain.StatusOpen {
		t.Fatal("step completion must leave parent open")
	}
	if err := svc.Undo(ctx, done); err != nil {
		t.Fatal(err)
	}
	steps, _ = svc.ListSteps(ctx, task.ID)
	st, _ = svc.Execution(ctx)
	if steps[0].Done || selected(t, st).AccumulatedMS != 105000 || selected(t, st).State != "ended" {
		t.Fatal("undo rewound time or did not reopen")
	}
	if _, err := svc.ChangeFocus(ctx, first.ID, "resume"); !errors.Is(err, service.ErrConflict) {
		t.Fatal("ended focus resumed", err)
	}
	// Repeated completion cannot re-complete a reopened target.
	st, err = svc.CompleteFocus(ctx, newBatch(t, svc, domain.ActorUser), first.ID)
	if err != nil {
		t.Fatal(err)
	}
	steps, _ = svc.ListSteps(ctx, task.ID)
	if steps[0].Done {
		t.Fatal("old complete retried after undo")
	}
}
func TestFocusDirectTaskEligibilityAndConcurrentStarts(t *testing.T) {
	ctx := context.Background()
	svc, _ := newService(t)
	task := createTask(t, svc, "Direct")
	blocker := createTask(t, svc, "Wait")
	b := newBatch(t, svc, domain.ActorUser)
	if err := svc.AddBlocker(ctx, b, domain.Blocker{TaskID: task.ID, BlockedByID: blocker.ID}); err != nil {
		t.Fatal(err)
	}
	start := func(token string) (domain.ExecutionState, error) {
		return svc.StartFocus(ctx, service.StartFocusFields{TaskID: task.ID, RequestID: token})
	}
	if _, err := start("blocked-start"); !errors.Is(err, service.ErrConflict) {
		t.Fatal(err)
	}
	if _, err := svc.SetTaskStatus(ctx, newBatch(t, svc, domain.ActorUser), blocker.ID, domain.StatusDone); err != nil {
		t.Fatal(err)
	}
	var wg gosync.WaitGroup
	errs := make(chan error, 8)
	for i := 0; i < 8; i++ {
		wg.Add(1)
		go func(i int) { defer wg.Done(); _, err := start(fmt.Sprintf("request-%d", i)); errs <- err }(i)
	}
	wg.Wait()
	close(errs)
	for err := range errs {
		if err != nil {
			t.Fatal(err)
		}
	}
	st, _ := svc.Execution(ctx)
	if len(st.Sessions) != 1 {
		t.Fatalf("concurrent start duplicate: %d", len(st.Sessions))
	}
	f := selected(t, st)
	if _, err := svc.ChangeFocus(ctx, f.ID, "end"); err != nil {
		t.Fatal(err)
	}
	got, _ := svc.GetTask(ctx, task.ID)
	if got.Status != domain.StatusOpen {
		t.Fatal("ending completed task")
	}
	st, err := start("fresh-start")
	if err != nil {
		t.Fatal(err)
	}
	f = selected(t, st)
	if _, err := svc.CompleteFocus(ctx, newBatch(t, svc, domain.ActorUser), f.ID); err != nil {
		t.Fatal(err)
	}
	got, _ = svc.GetTask(ctx, task.ID)
	if got.Status != domain.StatusDone {
		t.Fatal("direct completion failed")
	}
}

func TestCorrectFocusLedgerWhileRunningAndUndo(t *testing.T) {
	ctx := context.Background()
	svc, clock := newService(t)
	task := createTask(t, svc, "Keep time")
	st, err := svc.StartFocus(ctx, service.StartFocusFields{TaskID: task.ID, RequestID: "correction-start"})
	if err != nil {
		t.Fatal(err)
	}
	id := selected(t, st).ID
	clock.advance(2 * time.Minute)
	st, err = svc.CorrectFocus(ctx, id, 30000, "Timer left on", "correction-first", domain.ActorUser)
	if err != nil {
		t.Fatal(err)
	}
	f := selected(t, st)
	if f.State != "running" || f.AccumulatedMS != 30000 || !f.SegmentStartedAt.Equal(clock.now()) || len(st.Corrections) != 1 || st.Corrections[0].BeforeMS != 120000 {
		t.Fatalf("correction: %+v", st)
	}
	clock.advance(10 * time.Second)
	st, err = svc.CorrectFocus(ctx, id, 30000, "Timer left on", "correction-first", domain.ActorUser)
	if err != nil {
		t.Fatal(err)
	}
	if len(st.Corrections) != 1 || !selected(t, st).SegmentStartedAt.Equal(clock.now().Add(-10*time.Second)) {
		t.Fatal("retry reapplied correction")
	}
	st, err = svc.ChangeFocus(ctx, id, "pause")
	if err != nil {
		t.Fatal(err)
	}
	if selected(t, st).AccumulatedMS != 40000 {
		t.Fatal("correction double counted segment")
	}
	batch := newBatch(t, svc, domain.ActorUser)
	if _, err := svc.UpdateTask(ctx, batch, task.ID, service.TaskFields{Title: task.Title, EstimateMinutes: ptr(99)}); err != nil {
		t.Fatal(err)
	}
	if err := svc.Undo(ctx, batch); err != nil {
		t.Fatal(err)
	}
	st, _ = svc.Execution(ctx)
	if selected(t, st).AccumulatedMS != 40000 || len(st.Corrections) != 1 {
		t.Fatal("task undo reset corrected time")
	}
	st, err = svc.CorrectFocus(ctx, id, 0, "Clear mistaken session", "correction-zero", domain.ActorAgent)
	if err != nil {
		t.Fatal(err)
	}
	if selected(t, st).AccumulatedMS != 0 || st.Corrections[1].Actor != domain.ActorAgent {
		t.Fatal("agent correction not recorded")
	}
	for _, ms := range []int64{-1, 31536000001} {
		if _, err := svc.CorrectFocus(ctx, id, ms, "invalid", "invalid-correction", domain.ActorUser); !errors.Is(err, service.ErrInvalid) {
			t.Fatal(err)
		}
	}
	if _, err := svc.CorrectFocus(ctx, id, 5, " ", "empty-reason", domain.ActorUser); !errors.Is(err, service.ErrInvalid) {
		t.Fatal(err)
	}
	if _, err := svc.CorrectFocus(ctx, id, 5, "different", "correction-zero", domain.ActorUser); !errors.Is(err, service.ErrConflict) {
		t.Fatal("token collision accepted", err)
	}
}

func TestCompleteStepWithoutStartingAndKeepActionTimeAcrossSegments(t *testing.T) {
	ctx := context.Background()
	svc, clock := newService(t)
	task := createTask(t, svc, "Known task")
	batch := func() int64 { return newBatch(t, svc, domain.ActorUser) }
	step, err := svc.CreateStep(ctx, batch(), task.ID, service.StepFields{Description: "Already done", EstimateMinutes: ptr(2)})
	if err != nil {
		t.Fatal(err)
	}
	step, err = svc.UpdateStep(ctx, batch(), task.ID, step.ID, service.StepFields{Description: step.Description, EstimateMinutes: step.EstimateMinutes, Done: true})
	if err != nil {
		t.Fatal(err)
	}
	st, err := svc.Execution(ctx)
	if err != nil {
		t.Fatal(err)
	}
	if len(st.Sessions) != 0 {
		t.Fatal("completion invented focus time")
	}
	parent, _ := svc.GetTask(ctx, task.ID)
	if parent.Status != domain.StatusOpen {
		t.Fatal("completing step completed parent")
	}
	if _, err := svc.UpdateStep(ctx, batch(), task.ID, step.ID, service.StepFields{Description: step.Description, EstimateMinutes: step.EstimateMinutes}); err != nil {
		t.Fatal(err)
	}
	start := func(token string) domain.ExecutionState {
		t.Helper()
		st, err := svc.StartFocus(ctx, service.StartFocusFields{TaskID: task.ID, StepID: &step.ID, RequestID: token})
		if err != nil {
			t.Fatal(err)
		}
		return st
	}
	st = start("first-segment")
	first := *st.SelectedID
	clock.advance(50 * time.Second)
	if _, err := svc.ChangeFocus(ctx, first, "end"); err != nil {
		t.Fatal(err)
	}
	st = start("second-segment")
	clock.advance(time.Minute)
	st, err = svc.ChangeFocus(ctx, *st.SelectedID, "pause")
	if err != nil {
		t.Fatal(err)
	}
	var total int64
	for _, session := range st.Sessions {
		if session.TaskID == task.ID && session.StepID != nil && *session.StepID == step.ID {
			total += session.AccumulatedMS
		}
	}
	if total != 110000 {
		t.Fatalf("action time=%d want 110000", total)
	}
	// The next segment may go beyond the action estimate, and must keep recording.
	st, err = svc.ChangeFocus(ctx, *st.SelectedID, "resume")
	if err != nil {
		t.Fatal(err)
	}
	clock.advance(30 * time.Second)
	st, err = svc.Execution(ctx)
	if err != nil {
		t.Fatal(err)
	}
	if selected(t, st).State != "running" {
		t.Fatal("estimate expiry must not end work")
	}
	// The same checkbox path used in Day/Focus must stop active work, even without
	// routing completion through a session-specific timer button.
	if _, err := svc.UpdateStep(ctx, batch(), task.ID, step.ID, service.StepFields{Description: step.Description, EstimateMinutes: step.EstimateMinutes, Done: true}); err != nil {
		t.Fatal(err)
	}
	st, err = svc.Execution(ctx)
	if err != nil {
		t.Fatal(err)
	}
	if selected(t, st).State != "ended" {
		t.Fatal("direct step completion left work running")
	}
	total = 0
	for _, session := range st.Sessions {
		total += session.AccumulatedMS
	}
	if total != 140000 {
		t.Fatalf("direct completion lost or duplicated work: %d", total)
	}
	parent, _ = svc.GetTask(ctx, task.ID)
	if parent.Status != domain.StatusOpen {
		t.Fatal("direct step completion completed parent")
	}
}
