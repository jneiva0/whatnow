package agent_test

import (
	"context"
	"strings"
	"testing"

	"github.com/jameynakama/whatnow/internal/domain"
	"github.com/jameynakama/whatnow/internal/service"
)

func TestContextualDecompositionOnlyOrganizesStepsOfItsTask(t *testing.T) {
	f := newFixture(t)
	ctx := context.Background()
	task := f.task(t, service.TaskFields{Title: "Prepare for the day"})
	other := f.task(t, service.TaskFields{Title: "Different task"})
	batch, err := f.svc.NewBatch(ctx, domain.ActorUser)
	if err != nil {
		t.Fatal(err)
	}
	completed, err := f.svc.CreateStep(ctx, batch, task.ID, service.StepFields{Description: "Go upstairs", Done: true})
	if err != nil {
		t.Fatal(err)
	}
	f.model.replies = []any{calls(
		call("wrong-task", "organize_steps", map[string]any{"task_id": other.ID, "steps": []map[string]any{{"description": "Wrong action", "estimate_minutes": 2}}}),
		call("wrong-tool", "complete_task", map[string]any{"task_id": task.ID}),
		call("right-task", "organize_steps", map[string]any{"task_id": task.ID, "steps": []map[string]any{{"step_id": completed.ID}, {"description": "Get water", "estimate_minutes": 2}}}),
	), text("Added the next action.")}
	reply, err := f.agent.DecomposeTask(ctx, task.ID)
	if err != nil {
		t.Fatal(err)
	}
	if reply.Error != "" || len(reply.Tools) != 3 || reply.Tools[0].Error == "" || reply.Tools[1].Error == "" || reply.Tools[2].Error != "" {
		t.Fatalf("scoped reply: %+v", reply)
	}
	steps, err := f.svc.ListSteps(ctx, task.ID)
	if err != nil {
		t.Fatal(err)
	}
	if len(steps) != 2 || steps[0].ID != completed.ID || !steps[0].Done || steps[1].Description != "Get water" {
		t.Fatalf("progress: %+v", steps)
	}
	use := reply.Tools[2]
	if use.Outcome != "organized" || use.Action == nil || use.Action.TaskID != task.ID || use.Action.StepID == nil || *use.Action.StepID != steps[1].ID {
		t.Fatalf("missing applied execution target: %+v", use)
	}
	if reply.Tools[0].Action != nil || reply.Tools[1].Action != nil {
		t.Fatal("failed calls exposed an execution target")
	}
	foreign, _ := f.svc.ListSteps(ctx, other.ID)
	if len(foreign) != 0 {
		t.Fatal("decomposition modified another task")
	}
	parent, _ := f.svc.GetTask(ctx, task.ID)
	if parent.Status != domain.StatusOpen {
		t.Fatal("decomposition completed parent")
	}
	st, _ := f.svc.Execution(ctx)
	if len(st.Sessions) != 0 {
		t.Fatal("decomposition invented focus")
	}
	for _, tl := range f.model.chats[0].Tools {
		if tl.Name != "list_steps" && tl.Name != "organize_steps" && tl.Name != "split_step" {
			t.Fatalf("unexpected contextual tool: %s", tl.Name)
		}
	}
	history, err := f.agent.History(ctx)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(history[0].Text, "#1") || !strings.Contains(history[0].Text, task.Title) {
		t.Fatal("request lost task context")
	}
	if history[1].Tools[2].Action == nil || *history[1].Tools[2].Action.StepID != steps[1].ID {
		t.Fatal("action reference was lost when reloading chat history")
	}
}

func TestContextualDecompositionFailureLeavesManualStepsUsable(t *testing.T) {
	f := newFixture(t)
	ctx := context.Background()
	task := f.task(t, service.TaskFields{Title: "Manual task"})
	f.configured = false
	reply, err := f.agent.DecomposeTask(ctx, task.ID)
	if err != nil {
		t.Fatal(err)
	}
	if reply.Error == "" {
		t.Fatal("missing AI failure")
	}
	batch, err := f.svc.NewBatch(ctx, domain.ActorUser)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := f.svc.CreateStep(ctx, batch, task.ID, service.StepFields{Description: "Manual action"}); err != nil {
		t.Fatal(err)
	}
	steps, _ := f.svc.ListSteps(ctx, task.ID)
	if len(steps) != 1 {
		t.Fatal("manual creation failed after AI failure")
	}
}
