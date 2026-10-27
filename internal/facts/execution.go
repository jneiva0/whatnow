package facts

import (
	"sort"
	"time"

	"github.com/jameynakama/whatnow/internal/domain"
)

// Execution derives the shared view used by API, planner, capture and chat. Live
// steps are current scope; time on removed steps still belongs to their parent.
func Execution(now time.Time, tasks []domain.Task, steps []domain.Step, sessions []domain.FocusSession) map[int64]domain.TaskProgress {
	byTask := map[int64][]domain.Step{}
	for _, step := range steps {
		if step.DeletedAt == nil {
			byTask[step.TaskID] = append(byTask[step.TaskID], step)
		}
	}
	actual := map[int64]int64{}
	direct := map[int64]int64{}
	stepActual := map[int64]int64{}
	running := map[int64]*domain.ActionRef{}
	for _, session := range sessions {
		elapsed := max(int64(0), session.AccumulatedMS)
		if session.State == "running" && session.SegmentStartedAt != nil {
			elapsed += max(int64(0), now.Sub(*session.SegmentStartedAt).Milliseconds())
			running[session.TaskID] = &domain.ActionRef{TaskID: session.TaskID, StepID: session.StepID}
		}
		actual[session.TaskID] += elapsed
		if session.StepID == nil {
			direct[session.TaskID] += elapsed
		} else {
			stepActual[*session.StepID] += elapsed
		}
	}
	out := map[int64]domain.TaskProgress{}
	for _, task := range tasks {
		if task.DeletedAt != nil {
			continue
		}
		live := byTask[task.ID]
		sort.SliceStable(live, func(i, j int) bool {
			return live[i].Position < live[j].Position || (live[i].Position == live[j].Position && live[i].ID < live[j].ID)
		})
		p := domain.TaskProgress{TaskID: task.ID, TotalEstimate: domain.TaskEstimate(task, live),
			ActualMS: actual[task.ID], DirectActualMS: direct[task.ID], Running: running[task.ID], Steps: []domain.StepProgress{}}
		var known, missing, pendingMinutes int
		liveActual := int64(0)
		for _, step := range live {
			sp := domain.StepProgress{Step: step, ActualMS: stepActual[step.ID], Running: p.Running != nil && p.Running.StepID != nil && *p.Running.StepID == step.ID}
			p.Steps = append(p.Steps, sp)
			liveActual += sp.ActualMS
			if step.Done {
				p.CompletedSteps++
				continue
			}
			if p.Next == nil {
				id := step.ID
				p.Next = &domain.ActionRef{TaskID: task.ID, StepID: &id}
			}
			if step.EstimateMinutes == nil {
				missing++
			} else {
				known++
				pendingMinutes += *step.EstimateMinutes
			}
		}
		p.ArchivedActualMS = max(int64(0), p.ActualMS-p.DirectActualMS-liveActual)
		if len(live) == 0 {
			p.PendingEstimate = p.TotalEstimate
			if task.Status == domain.StatusOpen {
				p.Next = &domain.ActionRef{TaskID: task.ID}
			}
		} else {
			p.PendingEstimate.Partial = missing > 0
			if known > 0 || missing == 0 {
				p.PendingEstimate.Minutes = &pendingMinutes
			}
		}
		if task.Status != domain.StatusOpen {
			p.Next = nil
		}
		out[task.ID] = p
	}
	return out
}
