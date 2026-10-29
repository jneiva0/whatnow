package httpapi

import (
	"context"
	"github.com/jameynakama/whatnow/internal/domain"
	"github.com/jameynakama/whatnow/internal/service"
	"github.com/jameynakama/whatnow/internal/stream"
	"net/http"
)

func (s *Server) execution(w http.ResponseWriter, r *http.Request) {
	out, err := s.Service.Execution(r.Context())
	if err != nil {
		writeError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, out)
}
func (s *Server) stepWrite(w http.ResponseWriter, r *http.Request) {
	taskID, err := pathID(r, "id")
	if err != nil {
		writeError(w, err)
		return
	}
	var id int64
	if r.PathValue("step") != "" {
		id, err = pathID(r, "step")
		if err != nil {
			writeError(w, err)
			return
		}
	}
	var f service.StepFields
	if r.Method != "DELETE" {
		if err := decode(r, &f); err != nil {
			writeError(w, err)
			return
		}
	}
	status := http.StatusOK
	if r.Method == "POST" {
		status = http.StatusCreated
	}
	s.write(w, r, status, func(batch int64) (any, error) {
		switch r.Method {
		case "POST":
			return s.Service.CreateStep(r.Context(), batch, taskID, f)
		case "PUT":
			return s.Service.UpdateStep(r.Context(), batch, taskID, id, f)
		default:
			return nil, s.Service.DeleteStep(r.Context(), batch, taskID, id)
		}
	})
}
func (s *Server) reorderSteps(w http.ResponseWriter, r *http.Request) {
	id, err := pathID(r, "id")
	if err != nil {
		writeError(w, err)
		return
	}
	var body struct {
		IDs []int64 `json:"ids"`
	}
	if err := decode(r, &body); err != nil {
		writeError(w, err)
		return
	}
	s.write(w, r, http.StatusOK, func(batch int64) (any, error) { return s.Service.ReorderSteps(r.Context(), batch, id, body.IDs) })
}
func (s *Server) splitStep(w http.ResponseWriter, r *http.Request) {
	id, err := pathID(r, "id")
	if err != nil {
		writeError(w, err)
		return
	}
	step, err := pathID(r, "step")
	if err != nil {
		writeError(w, err)
		return
	}
	var body struct {
		Steps []service.StepFields `json:"steps"`
	}
	if err := decode(r, &body); err != nil {
		writeError(w, err)
		return
	}
	s.write(w, r, http.StatusOK, func(batch int64) (any, error) { return s.Service.SplitStep(r.Context(), batch, id, step, body.Steps) })
}
func (s *Server) focusChanged() {
	s.Hub.Publish(stream.Message{Type: "focus.changed", Data: struct{}{}})
}
func (s *Server) startFocus(w http.ResponseWriter, r *http.Request) {
	var in service.StartFocusFields
	if err := decode(r, &in); err != nil {
		writeError(w, err)
		return
	}
	out, err := s.Service.StartFocus(r.Context(), in)
	if err != nil {
		writeError(w, err)
		return
	}
	s.focusChanged()
	writeJSON(w, http.StatusOK, out)
}
func (s *Server) changeFocus(w http.ResponseWriter, r *http.Request) {
	id, err := pathID(r, "id")
	if err != nil {
		writeError(w, err)
		return
	}
	if r.PathValue("action") == "correct" {
		var correction struct {
			TotalMS   int64  `json:"total_ms"`
			Reason    string `json:"reason"`
			RequestID string `json:"request_id"`
		}
		if err := decode(r, &correction); err != nil {
			writeError(w, err)
			return
		}
		out, err := s.Service.CorrectFocus(r.Context(), id, correction.TotalMS, correction.Reason, correction.RequestID, domain.ActorUser)
		if err != nil {
			writeError(w, err)
			return
		}
		s.focusChanged()
		writeJSON(w, http.StatusOK, out)
		return
	}
	var body struct{}
	if err := decode(r, &body); err != nil {
		writeError(w, err)
		return
	}
	var out domain.ExecutionState
	if r.PathValue("action") == "complete" {
		batch, e := s.Service.NewBatch(r.Context(), domain.ActorUser)
		if e != nil {
			writeError(w, e)
			return
		}
		out, err = s.Service.CompleteFocus(r.Context(), batch, id)
		if err == nil {
			s.publishBatch(r.Context(), batch)
		}
	} else {
		out, err = s.Service.ChangeFocus(r.Context(), id, r.PathValue("action"))
	}
	if err != nil {
		writeError(w, err)
		return
	}
	s.focusChanged()
	writeJSON(w, http.StatusOK, out)
}

// decomposeTask is a contextual entry to the existing chat agent, restricted to this task.
func (s *Server) decomposeTask(w http.ResponseWriter, r *http.Request) {
	id, err := pathID(r, "id")
	if err != nil {
		writeError(w, err)
		return
	}
	var body struct{}
	if err := decode(r, &body); err != nil {
		writeError(w, err)
		return
	}
	reply, err := s.Agent.DecomposeTask(context.WithoutCancel(r.Context()), id)
	if err != nil {
		writeError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, reply)
}
