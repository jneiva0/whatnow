package domain

// Estimate is derived, never written over the user's explicit estimate or actual time.
type Estimate struct {
	Minutes *int `json:"minutes"`
	Partial bool `json:"partial"`
}

func TaskEstimate(t Task, steps []Step) Estimate {
	var sum, known, missing int
	for _, s := range steps {
		if s.TaskID != t.ID || s.DeletedAt != nil {
			continue
		}
		if s.EstimateMinutes == nil {
			missing++
		} else {
			sum += *s.EstimateMinutes
			known++
		}
	}
	e := Estimate{Partial: missing > 0}
	if t.EstimateMinutes != nil {
		n := max(*t.EstimateMinutes, sum)
		e.Minutes = &n
	} else if known > 0 {
		e.Minutes = &sum
	}
	return e
}
