# Execution inside Whatnow

Mode: Operate. Build path: code (explicit user choice).
Scope: extend the established Today/task-editor experience and add a focused execution route using the same visual language.

FIRST VIEWPORT: readable task title/context, selected concrete action, prominent estimate-versus-actual timer and progress indicator, checkbox completion even before starting, and one play/pause icon toggle. Keep current action, timer and play/pause sticky while scrolling on desktop and mobile, without hiding content or touch targets. Keep the task’s full step list and their actual/planned times continuous with this context; task details/editing stay accessible. Compact editable energy, size and free-text feeling sit directly beneath the task title, shared with Day. The pencil editor does not take mobile title width. Ver dia returns to the choosing surface without pausing. No second task chooser below the viewport. Without an estimate show elapsed time, never a fabricated budget. Day cards expose collapsible steps with checkboxes/play-pause and task time, compact focus entry and compact add-step entry; no oversized focus footer.

Signature interaction: choose a concrete task or pending step, start focus, pause on either client and see the authoritative state follow on the other. Reopening the application restores elapsed time from server timestamps.

Motion grammar: existing Whatnow button feedback and restrained state changes, reduced-motion respected. Timer changes do not animate layout.

Quality bar: existing Whatnow cream/yellow/purple surfaces, thick rounded dark contours, solid offset shadows, heavy headings. Preserve this pinned system; do not add a new identity, concept tournament, image mockups, reports, projects, rigid calendar, or voice.

Required states: loading, no focus, running, paused, over action estimate, no estimate, completed step awaiting explicit next action, all steps complete awaiting explicit parent completion, disconnected, failed mutation. No optimistic saved state on failure.

Confirmed correction (2026-10-25): the assistant introduced independent session duration/default15; João rejected it. Remove it from API, persistence contract and UI. Preserve accumulated time, corrections, SSE/polling and the existing day timeline. Refine Operate through Impeccable distill/clarify in code; no new interview, visual tournament, screenshot/reviewer gate or extensive browser QA. João owns real-use visual/functional acceptance.

Wireframe refinement (2026-10-25): timer must remain the main focus information, not a tiny line. User correction overrides any earlier blanket simplification of timer size. The annotated image supplies organization evidence, not a complete approved comp. No duplicate large completion buttons where a checkbox already completes the same target.

Shared-interaction round U01–U10: references govern interaction, not merely appearance. One ExecutionSteps row/list serves Day, Focus and editor; text/estimate edit directly, checkbox completes/reopens, icon toggles play/pause, grip reorders by pointer or keyboard, adder is inline, secondary split names its consequence. Visible task decomposition uses existing scoped chat and leaves manual editing available on failure. Actual running focus owns the main timer regardless of URL/view selection. Completing current chooses next stopped, never automatically starts. All compact actual/estimated pairs use MM:SS / MM:SS or —. Existing time ledger, corrections, SSE and polling persist unchanged.
