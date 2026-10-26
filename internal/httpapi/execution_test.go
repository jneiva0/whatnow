package httpapi_test

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"testing"

	"github.com/jameynakama/whatnow/internal/domain"
	"github.com/jameynakama/whatnow/internal/httpapi"
)

func TestExecutionRESTAndEvents(t *testing.T) {
	a := newAPI(t)
	msgs, unsub := a.hub.Subscribe("focus.changed")
	defer unsub()
	taskResponse := a.do("POST", "/api/v1/tasks", map[string]any{"title": "Write", "raw_text": "Original capture", "estimate_minutes": 20, "size": "M", "energy": "low", "feeling": "uncertain"})
	if taskResponse.Code != 201 {
		t.Fatal(taskResponse.Body.String())
	}
	var taskMutation struct {
		Data domain.Task `json:"data"`
	}
	json.Unmarshal(taskResponse.Body.Bytes(), &taskMutation)
	task := taskMutation.Data
	path := fmt.Sprintf("/api/v1/tasks/%d/steps", task.ID)
	// Unknown nesting fields are rejected rather than silently creating a second level.
	nested := a.do("POST", path, map[string]any{"description": "child", "parent_step_id": 1})
	if nested.Code != 400 {
		t.Fatal(nested.Code, nested.Body.String())
	}
	made := a.do("POST", path, map[string]any{"description": "Open document", "estimate_minutes": 2})
	if made.Code != 201 {
		t.Fatal(made.Body.String())
	}
	var mutation struct {
		BatchID int64       `json:"batch_id"`
		Data    domain.Step `json:"data"`
	}
	json.Unmarshal(made.Body.Bytes(), &mutation)
	start := map[string]any{"task_id": task.ID, "step_id": mutation.Data.ID, "request_id": "request-http-1"}
	res := a.do("POST", "/api/v1/focus/start", start)
	if res.Code != 200 {
		t.Fatal(res.Code, res.Body.String())
	}
	var st domain.ExecutionState
	json.Unmarshal(res.Body.Bytes(), &st)
	if st.SelectedID == nil || len(st.Sessions) != 1 {
		t.Fatal(st)
	}
	if strings.Contains(res.Body.String(), "duration_minutes") {
		t.Fatal("API exposed rejected independent duration")
	}
	select {
	case <-msgs:
	default:
		t.Fatal("missing focus SSE event")
	}
	for _, action := range []string{"pause", "resume", "complete", "complete"} {
		res = a.do("POST", fmt.Sprintf("/api/v1/focus/%d/%s", *st.SelectedID, action), map[string]any{})
		if res.Code != 200 {
			t.Fatal(action, res.Code, res.Body.String())
		}
	}
	tasks := a.do("GET", "/api/v1/tasks", nil)
	var all []domain.Task
	json.Unmarshal(tasks.Body.Bytes(), &all)
	if all[0].Status != domain.StatusOpen || all[0].RawText != "Original capture" || all[0].EstimateMinutes == nil || *all[0].EstimateMinutes != 20 {
		t.Fatalf("task changed: %+v", all)
	}
	execution := a.do("GET", "/api/v1/execution", nil)
	if execution.Code != 200 {
		t.Fatal(execution.Body.String())
	}
	json.Unmarshal(execution.Body.Bytes(), &st)
	if !st.Steps[0].Done || st.Sessions[0].State != "ended" {
		t.Fatal(st)
	}
}
func TestExecutionStepUndoAndValidationREST(t *testing.T) {
	a := newAPI(t)
	task := a.do("POST", "/api/v1/tasks", map[string]any{"title": "Direct"})
	var m struct {
		Data domain.Task `json:"data"`
	}
	json.Unmarshal(task.Body.Bytes(), &m)
	path := fmt.Sprintf("/api/v1/tasks/%d/steps", m.Data.ID)
	made := a.do("POST", path, map[string]any{"description": "draft"})
	var mutation httpapi.Mutation
	json.Unmarshal(made.Body.Bytes(), &mutation)
	undo := a.do("POST", fmt.Sprintf("/api/v1/batches/%d/undo", mutation.BatchID), map[string]any{})
	if undo.Code != 204 {
		t.Fatal(undo.Body.String())
	}
	for _, body := range []map[string]any{{"task_id": m.Data.ID, "duration_minutes": 15, "request_id": "removed-budget"}, {"task_id": m.Data.ID, "request_id": "x"}} {
		res := a.do("POST", "/api/v1/focus/start", body)
		if res.Code != http.StatusBadRequest {
			t.Fatal(res.Code, res.Body.String())
		}
	}
}

func TestTaskDecompositionRESTKeepsManualFallback(t *testing.T) {
	a := newAPI(t)
	made := a.do("POST", "/api/v1/tasks", map[string]any{"title": "Prepare"})
	var created struct {
		Data domain.Task `json:"data"`
	}
	if err := json.Unmarshal(made.Body.Bytes(), &created); err != nil {
		t.Fatal(err)
	}
	// No provider is configured in this fixture: the existing chat contract reports a message error.
	res := a.do("POST", fmt.Sprintf("/api/v1/tasks/%d/decompose", created.Data.ID), map[string]any{})
	if res.Code != 200 {
		t.Fatal(res.Code, res.Body.String())
	}
	var reply domain.Message
	if err := json.Unmarshal(res.Body.Bytes(), &reply); err != nil {
		t.Fatal(err)
	}
	if reply.Error == "" {
		t.Fatal("AI failure presented as success")
	}
	manual := a.do("POST", fmt.Sprintf("/api/v1/tasks/%d/steps", created.Data.ID), map[string]any{"description": "Get water"})
	if manual.Code != 201 {
		t.Fatal(manual.Body.String())
	}
	missing := a.do("POST", "/api/v1/tasks/999/decompose", map[string]any{})
	if missing.Code != 404 {
		t.Fatal(missing.Code, missing.Body.String())
	}
}
