package prompt

import (
	"fmt"
	"strings"
	"time"

	"github.com/jameynakama/whatnow/internal/domain"
	"github.com/jameynakama/whatnow/internal/facts"
)

// Facts renders the facts of the day as live context for the model.
func Facts(f facts.Facts) string {
	var b strings.Builder
	fmt.Fprintf(&b, "Now: %s (%s). The working day ends at %s.\n", f.Now.Format("Monday 2006-01-02 15:04"),
		f.Now.Format("MST"), f.DayEnd.Format("15:04"))

	b.WriteString("Free time left today:")
	if len(f.Gaps) == 0 {
		b.WriteString(" none")
	}
	for i, g := range f.Gaps {
		if i > 0 {
			b.WriteByte(',')
		}
		fmt.Fprintf(&b, " %s-%s (%s)", g.Start.Format("15:04"), g.End.Format("15:04"), facts.Duration(g.Duration()))
	}

	b.WriteString("\n\nToday's events:\n")
	if len(f.Events) == 0 {
		b.WriteString("(none)\n")
	}
	for _, e := range f.Events {
		start := e.StartsAt.In(f.Now.Location()).Format("15:04")
		if e.EndsAt != nil {
			fmt.Fprintf(&b, "- %s-%s %s\n", start, e.EndsAt.In(f.Now.Location()).Format("15:04"), e.Title)
		} else {
			fmt.Fprintf(&b, "- %s %s\n", start, e.Title)
		}
	}

	b.WriteString("\nEligible tasks, as #<id> title:\n")
	if len(f.Tasks) == 0 {
		b.WriteString("(none)\n")
	}
	for _, tf := range f.Tasks {
		t := tf.Task
		fmt.Fprintf(&b, "#%d %s", t.ID, t.Title)
		var details []string
		if tf.Estimate.Minutes != nil {
			details = append(details, fmt.Sprintf("estimated total %d min", *tf.Estimate.Minutes))
		}
		if tf.Estimate.Partial {
			details = append(details, "some steps have no estimate")
		}
		if t.Energy != nil {
			details = append(details, string(*t.Energy)+" energy")
		}
		if len(details) > 0 {
			fmt.Fprintf(&b, " [%s]", strings.Join(details, ", "))
		}
		if t.Feeling != "" {
			fmt.Fprintf(&b, " feeling: %q", t.Feeling)
		}
		b.WriteByte('\n')
		if t.Notes != "" {
			fmt.Fprintf(&b, "  Notes: %q\n", t.Notes)
		}
		if t.RawText != "" && t.RawText != t.Title && t.RawText != t.Notes {
			fmt.Fprintf(&b, "  Original capture: %q\n", t.RawText)
		}
		if tf.Progress != nil {
			b.WriteString(Execution(*tf.Progress))
		}
		for _, a := range tf.Annotations {
			fmt.Fprintf(&b, "  - %s\n", a)
		}
	}
	return b.String()
}

// Execution renders computed facts, not raw session history. Times are cumulative
// recorded focus, never a claim of activity or time spent on a particular day.
func Execution(p domain.TaskProgress) string {
	var b strings.Builder
	fmt.Fprintf(&b, "  Recorded focus (cumulative): %s; completed steps: %d/%d.\n", time.Duration(p.ActualMS)*time.Millisecond, p.CompletedSteps, len(p.Steps))
	if p.PendingEstimate.Minutes != nil {
		fmt.Fprintf(&b, "  Pending action estimates: %d min (partial=%t); not reduced by recorded time.\n", *p.PendingEstimate.Minutes, p.PendingEstimate.Partial)
	} else {
		b.WriteString("  Pending action estimates: unknown.\n")
	}
	if p.ArchivedActualMS > 0 {
		fmt.Fprintf(&b, "  Recorded focus on removed actions (included in task total): %s.\n", time.Duration(p.ArchivedActualMS)*time.Millisecond)
	}
	if p.Focus != nil {
		target := "task directly"
		if p.Focus.StepID != nil {
			target = fmt.Sprintf("step #%d", *p.Focus.StepID)
		}
		fmt.Fprintf(&b, "  Selected focus: session #%d, %s, %s.\n", p.Focus.ID, p.Focus.State, target)
	}
	if p.Running != nil {
		if p.Running.StepID == nil {
			b.WriteString("  Currently running: task directly.\n")
		} else {
			fmt.Fprintf(&b, "  Currently running: step #%d.\n", *p.Running.StepID)
		}
	}
	for _, step := range p.Steps {
		estimate := "unknown"
		if step.EstimateMinutes != nil {
			estimate = fmt.Sprintf("%d min", *step.EstimateMinutes)
		}
		fmt.Fprintf(&b, "  Step #%d (done=%t): %s; estimate %s; recorded %s.\n", step.ID, step.Done, step.Description, estimate, time.Duration(step.ActualMS)*time.Millisecond)
	}
	return b.String()
}
