package domain

// ActionRef identifies an existing executable action. A nil step means the task itself.
type ActionRef struct {
	TaskID int64  `json:"task_id"`
	StepID *int64 `json:"step_id"`
}

// StepProgress is current state, not an inference about what the user did while timed.
type StepProgress struct {
	Step
	ActualMS int64 `json:"actual_ms"`
	Running  bool  `json:"running"`
}

// TaskProgress separates scope, pending estimates and recorded time. PendingEstimate
// is the sum of known pending action estimates, not total minus elapsed time.
type TaskProgress struct {
	TaskID           int64          `json:"task_id"`
	TotalEstimate    Estimate       `json:"total_estimate"`
	PendingEstimate  Estimate       `json:"pending_estimate"`
	CompletedSteps   int            `json:"completed_steps"`
	ActualMS         int64          `json:"actual_ms"`
	DirectActualMS   int64          `json:"direct_actual_ms"`
	ArchivedActualMS int64          `json:"archived_actual_ms"`
	Steps            []StepProgress `json:"steps"`
	Next             *ActionRef     `json:"next"`
	Running          *ActionRef     `json:"running"`
}
