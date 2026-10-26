package agent

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/jameynakama/whatnow/internal/domain"
	js "github.com/jameynakama/whatnow/internal/jsonschema"
	"github.com/jameynakama/whatnow/internal/service"
	"github.com/jameynakama/whatnow/internal/stream"
)

var stepFieldsSchema = js.Object(js.Prop("description", js.String("One concrete action, at the same level as other steps.")), js.Prop("estimate_minutes", js.Integer("Estimated action minutes, 1..1440, or 0 for unspecified.")))
var executionTools = []tool{
	{"correct_focus_time", "Correct recorded actual elapsed time ONLY when explicitly requested. This does not edit estimates, and is not reverted by generic undo. Read execution_context to find the session ID. Supply a new unique request_id; reuse it only to retry the same correction.", js.Object(js.Prop("session_id", js.Integer("Focus session ID.")), js.Prop("total_ms", js.Integer("Correct actual elapsed total in milliseconds, 0..31536000000.")), js.Prop("reason", js.String("The user's explicit correction reason.")), js.Prop("request_id", js.String("Unique correction request token, 8..200 characters."))), correctFocusTime},
	{"list_steps", "Read a task's ordered steps and completion state.", js.Object(js.Prop("task_id", taskID)), listTaskSteps},
	{"add_step", "Add a lightweight action to an existing task, only on request. Steps never become planning cards.", js.Object(js.Prop("task_id", taskID), js.Prop("description", js.String("Concrete action.")), js.Prop("estimate_minutes", js.Integer("1..1440 minutes, or 0 for unspecified."))), addTaskStep},
	{"update_step", "Edit, complete or reopen a step. Completing a step never completes its parent.", js.Object(js.Prop("task_id", taskID), js.Prop("step_id", js.Integer("Step ID.")), js.Prop("description", js.String("Complete action description.")), js.Prop("estimate_minutes", js.Integer("1..1440 minutes, or 0 for unspecified.")), js.Prop("status", js.Enum("Completion state.", "pending", "done"))), updateTaskStep},
	{"delete_step", "Remove an action, preserving its recorded focus time.", js.Object(js.Prop("task_id", taskID), js.Prop("step_id", js.Integer("Step ID."))), deleteTaskStep},
	{"reorder_steps", "Reorder all live steps of a task. Include each step ID exactly once.", js.Object(js.Prop("task_id", taskID), js.Prop("ids", js.Array(js.Integer("Step ID.")))), reorderTaskSteps},
	{"split_step", "Replace one pending action with 2..30 smaller actions at the same level. Preserve completed siblings and focus time.", js.Object(js.Prop("task_id", taskID), js.Prop("step_id", js.Integer("Pending step ID.")), js.Prop("steps", js.Array(stepFieldsSchema))), splitTaskStep},
	{"execution_context", "Read current focus, earlier sessions and steps. Starting/stopping the timer is controlled by the user.", js.Object(), executionContext},
}

func init() { tools = append(tools, executionTools...) }
func estimated(n int) *int {
	if n == 0 {
		return nil
	}
	return &n
}
func listTaskSteps(ctx context.Context, t *turn, input json.RawMessage) (string, any, error) {
	in, err := args[struct {
		TaskID int64 `json:"task_id"`
	}](input)
	if err != nil {
		return "", nil, err
	}
	out, err := t.agent.svc.ListSteps(ctx, in.TaskID)
	return "Read task steps", out, err
}
func addTaskStep(ctx context.Context, t *turn, input json.RawMessage) (string, any, error) {
	in, err := args[struct {
		TaskID      int64  `json:"task_id"`
		Description string `json:"description"`
		Estimate    int    `json:"estimate_minutes"`
	}](input)
	if err != nil {
		return "", nil, err
	}
	batch, err := t.batch(ctx)
	if err != nil {
		return "", nil, err
	}
	out, err := t.agent.svc.CreateStep(ctx, batch, in.TaskID, service.StepFields{Description: in.Description, EstimateMinutes: estimated(in.Estimate)})
	return "Added a step", out, err
}
func updateTaskStep(ctx context.Context, t *turn, input json.RawMessage) (string, any, error) {
	in, err := args[struct {
		TaskID      int64  `json:"task_id"`
		StepID      int64  `json:"step_id"`
		Description string `json:"description"`
		Estimate    int    `json:"estimate_minutes"`
		Status      string `json:"status"`
	}](input)
	if err != nil {
		return "", nil, err
	}
	if in.Status != "pending" && in.Status != "done" {
		return "", nil, fmt.Errorf("%w: invalid step status", service.ErrInvalid)
	}
	batch, err := t.batch(ctx)
	if err != nil {
		return "", nil, err
	}
	out, err := t.agent.svc.UpdateStep(ctx, batch, in.TaskID, in.StepID, service.StepFields{Description: in.Description, EstimateMinutes: estimated(in.Estimate), Done: in.Status == "done"})
	return "Updated a step", out, err
}
func deleteTaskStep(ctx context.Context, t *turn, input json.RawMessage) (string, any, error) {
	in, err := args[struct {
		TaskID int64 `json:"task_id"`
		StepID int64 `json:"step_id"`
	}](input)
	if err != nil {
		return "", nil, err
	}
	batch, err := t.batch(ctx)
	if err != nil {
		return "", nil, err
	}
	err = t.agent.svc.DeleteStep(ctx, batch, in.TaskID, in.StepID)
	return "Removed a step", map[string]bool{"deleted": err == nil}, err
}
func reorderTaskSteps(ctx context.Context, t *turn, input json.RawMessage) (string, any, error) {
	in, err := args[struct {
		TaskID int64   `json:"task_id"`
		IDs    []int64 `json:"ids"`
	}](input)
	if err != nil {
		return "", nil, err
	}
	batch, err := t.batch(ctx)
	if err != nil {
		return "", nil, err
	}
	out, err := t.agent.svc.ReorderSteps(ctx, batch, in.TaskID, in.IDs)
	return "Reordered steps", out, err
}
func splitTaskStep(ctx context.Context, t *turn, input json.RawMessage) (string, any, error) {
	in, err := args[struct {
		TaskID int64 `json:"task_id"`
		StepID int64 `json:"step_id"`
		Steps  []struct {
			Description string `json:"description"`
			Estimate    int    `json:"estimate_minutes"`
		} `json:"steps"`
	}](input)
	if err != nil {
		return "", nil, err
	}
	fields := []service.StepFields{}
	for _, st := range in.Steps {
		fields = append(fields, service.StepFields{Description: st.Description, EstimateMinutes: estimated(st.Estimate)})
	}
	batch, err := t.batch(ctx)
	if err != nil {
		return "", nil, err
	}
	out, err := t.agent.svc.SplitStep(ctx, batch, in.TaskID, in.StepID, fields)
	return "Made the action smaller", out, err
}
func executionContext(ctx context.Context, t *turn, _ json.RawMessage) (string, any, error) {
	out, err := t.agent.svc.Execution(ctx)
	return "Read focus and progress", out, err
}

func correctFocusTime(ctx context.Context, t *turn, input json.RawMessage) (string, any, error) {
	in, err := args[struct {
		SessionID int64  `json:"session_id"`
		TotalMS   int64  `json:"total_ms"`
		Reason    string `json:"reason"`
		RequestID string `json:"request_id"`
	}](input)
	if err != nil {
		return "", nil, err
	}
	out, err := t.agent.svc.CorrectFocus(ctx, in.SessionID, in.TotalMS, in.Reason, in.RequestID, domain.ActorAgent)
	if err == nil {
		t.agent.hub.Publish(stream.Message{Type: "focus.changed", Data: struct{}{}})
	}
	return "Corrected recorded focus time", out, err
}
