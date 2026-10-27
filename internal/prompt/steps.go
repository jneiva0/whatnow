package prompt

import js "github.com/jameynakama/whatnow/internal/jsonschema"

// StepSequence is shared by capture and chat. This is arrangement, not a
// replacement list: omitted actions survive, and existing identities are reused.
var StepSequence = js.Array(js.Object(
	js.Prop("step_id", js.Integer("Existing live action ID ONLY when it belongs in the requested starting sequence, or 0 for a NEW action. Do not list unrelated old actions to preserve them: the service preserves omissions automatically.")),
	js.Prop("description", js.String("Concrete new action. Empty when reusing a step_id. Preparation before work is an addition, not a replacement of that work.")),
	js.Prop("estimate_minutes", js.Integer("Rough minutes for a new action, 1..1440. For familiar concrete actions, give a reasonable rough estimate; 0 only when genuinely unclear, or when reusing a step_id.")),
	js.Prop("repeat", map[string]any{"type": "boolean", "description": "True only for a deliberate NEW occurrence of an already represented action, e.g. drinking water again. Never reopen the completed occurrence."}),
))

const StepOrganization = `Return only the requested starting sequence. Do not copy old steps just to preserve them: the service already keeps every omission. Reuse a pending step ID when it is part of that sequence. A deliberate repeat such as drinking water AGAIN is a new action with repeat=true, not a reopening. Put the requested starting sequence BEFORE other pending work, not merely appended somewhere in the task. Insert preparation without replacing unrelated work. Preserve completed progress and recorded focus. For familiar concrete actions, provide realistic rough minute estimates, not a fixed budget or zero merely because the user did not dictate a number. Unknown is valid when the work is genuinely unclear. Organize only what helps them begin now, not an exhaustive checklist.`

const Organization = `When a request describes an actionable sequence or asks for help preparing work, organize that sequence in the existing task with organize_steps before replying. Merely reading steps or giving instructions in chat does not register an action. ` + StepOrganization + ` Use split_step only for a genuine decomposition of the SAME pending action into smaller parts of that action, never for preparation or a different sequence. Do not turn every feeling into tasks: conversation, reflection and an explicit request not to organize must remain conversation. If the intended task is genuinely ambiguous, ask one short question. Trust the applied result and its next reference; don't claim changes that failed.`
