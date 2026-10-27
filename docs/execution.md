# Steps and persistent focus

This feature extends the existing Go service, SQLite migrations, REST API, Svelte client, and SSE stream. Tasks remain the only planning units. Steps never become priority cards or nested tasks.

## Data and write rules

Migration `0009_execution.sql` adds optional task estimates, flat ordered steps, focus sessions, selected focus, and start request tokens. `0010_focus_corrections.sql` adds correction audit entries. `0011_action_time.sql` removes the rejected independent duration column while preserving every work row, ID, elapsed value, segment timestamp, correction, selection and request token. Migrations run when the application opens its chosen database.

Step writes use the existing batch/activity/undo service. Task and step estimates are optional integer minutes in 1..1440. A split replaces one pending step with at least two pending peers, retaining a tombstone for its historical sessions. Completed siblings and original task capture are preserved.

Focus writes use the same service layer and SQL transactions, but elapsed time is a separate ledger. A partial unique SQLite index permits one running session. Switching starts pauses the previous running session, settles its elapsed segment, retains its history, and selects the new session. Repeated starts use request tokens and do not duplicate or resurrect an earlier session. Pause/resume/end operate on a session ID. End does not complete a target; explicit complete changes the step or task in a batch and ends its sessions. Completing a step does not complete the parent task.

Generic undo restores task/step state without rewinding or resuming the time ledger. Unavailable targets (completed, deleted, split) close their sessions while preserving recorded time. Replanning writes plans only and does not change execution state.

## Action estimate and actual time

João rejected the assistant-introduced independent session duration during real-use QA on 2026-10-25. It is removed from the domain, API, client and current storage schema; historical migration 0009 records how the earlier preview was created, and migration 0011 upgrades it without losing recorded work. There is no default 15-minute budget.

- `estimate_minutes`: optional estimated work for a task or step; editing it never changes real recorded time.
- `accumulated_ms` plus the running segment from `segment_started_at`: actual dedicated time for a work row. Action time sums all rows matching the task and optional step exactly. Task totals include each work row once (direct work plus step work), without adding a rollup again. Historical step rows still contribute recorded time.
- The main counter is estimate minus accumulated action time: 2min minus 1m50 gives 10s left, and 2m10 gives +10s overtime. Without an estimate it shows elapsed time. Zero never completes, advances, or stops work. UI clock compensation and client-only ticks remain.

A deliberate correction supplies an absolute recorded-work total `total_ms`, a reason, and a request token. It writes a `focus_corrections` audit entry containing before/after totals, actor and server time. For a running row it sets the accumulated total and moves the segment anchor to the same server instant. A retried correction does not apply again. Corrections are independent of task/plan undo and do not modify estimates.

## Task continuity in Day and Focus

Day is the place to choose a different task. Its cards expose a visible start/focus action, actual/planned time, and a compact collapsible step list with completion/reopen and start controls. A previously finished action can be marked done without starting a timer. Focus shows the task context, accessible details/editing, selected action, a prominent time counter/progress indicator with sticky play/pause, and the complete ordered step list. Energy, size and free-text feeling are compact editable values immediately beneath the task title in both surfaces; a discrete pencil opens the full editor. Completion uses each task/step checkbox, including before any timer is started, and play/pause is a single accessible icon toggle. Completing a step leaves its parent open and offers the next action without starting it. When all steps are done the task remains available for explicit parent completion; no extra action is invented. Ver dia navigates without pausing; starting another action pauses the prior one through the unchanged service.

## REST and live updates

- `GET /api/v1/execution`: authoritative server time, steps (including tombstones for history), sessions, selection and correction audit.
- `POST /tasks/{id}/steps`, `PUT` and `DELETE /tasks/{id}/steps/{step}`; `POST /tasks/{id}/steps/reorder` and `/steps/{step}/split` (all under `/api/v1`). Step mutations return existing batch responses and publish existing batch events.
- `POST /api/v1/focus/start` with task, optional step and request ID (no duration field). `POST /api/v1/focus/{id}/{pause|resume|end|complete}` returns authoritative execution state. `POST /api/v1/focus/{id}/correct` accepts total milliseconds, reason and request ID.

Focus mutations publish `focus.changed`. Batch changes and undo also refetch execution. Opening or returning to a visible tab refetches. A visible page polls every five seconds only while SSE is disconnected and stops on reconnect. No new WebSocket or synchronization infrastructure is added.

Chat tools list/add/update/delete/reorder/split steps and read current execution. `correct_focus_time` changes actual time only on an explicit request. `update_task` accepts `estimate_minutes`; `update_step` accepts its own estimate. `What now` chooses an eligible task and its first pending step when present; both model and fallback paths apply this action selection. Timer controls work without a model.

## Local preview and verification boundary

The development preview prepared on 2026-10-25 uses this worktree's normal executable, port 7790, and `.demo/focus-preview.db`. It began as a separate test database with sample data. João now uses it with his own tasks and has configured AI; treat it as valuable user data. A read-only settings check confirmed a provider is configured without displaying secrets. Preserve its database and settings; never reseed or wipe it during preview updates. The existing service on port 7777, its database, and the existing Tailscale HTTPS proxy were not replaced. Preview URLs were verified by HTTP requests; a physical phone was not checked.

To restart the same preview after rebuilding:

```sh
cd web
bun run build
cd ..
cp -R web/build/. internal/web/dist/
go build -buildvcs=false -o whatnow ./cmd/whatnow
TZ=America/Sao_Paulo ./whatnow serve -addr 0.0.0.0:7790 -db .demo/focus-preview.db
```

The preview database allows the already configured private Tailscale hostname/IP as hosts. It was initialized separately from the service on 7777, and now contains João’s preview work. No production deployment was performed.

The initial delivery’s Go tests and vet passed, as did Svelte check and 157 frontend unit tests. In the current correction round Go tests/vet passed, Svelte check has zero errors/warnings, 159 frontend unit tests passed and the frontend build passed. The previous browser results below apply only to the earlier UI. Focused browser tests ran before João took over QA: the final run was 5/6 passing. The remaining test reached step removal/undo and failed because the toast was behind the editor; its z-index was fixed but browser QA was not repeated after the user's instruction. Two-context Chromium checks covered desktop/mobile viewport with touch, SSE pause, reload, correction, replan preservation, explicit step/parent completion, API failure/retry, and What now navigation. These are viewport checks and fake-model tests, not physical-device or live-provider validation.

João explicitly chose to perform visual and functional QA through real use and later confirmed desktop/mobile pause-resume synchronization and useful real brain-dump decomposition. Those are his reports, distinct from the earlier automated fake-model tests. The current correction round uses focused logic/build verification and a bounded acceptance check, not extensive browser QA or a visual review gate. The Svelte MCP accepted changed components; its stream TypeScript-module analysis encountered an internal tool error, while normal compilation and tests for that code passed.


## Current correction-round verification

A focused migration test creates an existing running work row plus selection, request token and correction audit, upgrades it, checks its time/anchor/state and references survived, and checks the rejected duration column is absent. Service tests cover completion without starting a timer and accumulated time across separate work rows; existing pause/resume/switch/idempotency/correction/undo tests remain. Relevant frontend timing and task-flow tests are updated to the action-time contract. See `docs/feedback.md` for per-item implementation evidence and the real-use acceptance still owned by João.


The refined preview was updated on 2026-10-25 at 14:31 America/Sao_Paulo using `.demo/whatnow-preview-next`, the same `.demo/focus-preview.db`, port 7790 and existing AI settings. A fresh consistent SQLite backup was made before restart. Read-only comparison with it confirmed all surviving columns/rows of tasks, steps, focus work, corrections, start tokens, selection and settings were identical after migration. SQLite quick_check and foreign_key_check passed. The rejected duration column and API field are absent. Local and Tailscale HTTP200 pages exactly matched the final build200 SHA256 `14dc1ff4889e497b6a619a64872f73a88bb5a220dd363a2393e650baf29de2f8`. The service on7777 retained its existing process. No browser or visual review was repeated; sticky behavior and the refined layout still need João’s real-use acceptance.


## Shared step interaction and focus representation (U01–U10)

`ExecutionSteps.svelte` now serves Day, Focus and TaskEditor; the parallel StepsEditor is removed. It supports inline text/estimate updates, direct done/reopen, one play/pause icon, compact actual/estimated time with dash for unknown, inline creation, contextual pending-step split and removal. Pointer grip reordering shares the existing reorder API with the keyboard-arrow alternative and accessible result announcement. Timed execution and editing a row are separate interactions.

`focusTarget` resolves presentation from authoritative execution state: an actually running session has priority over a different requested URL target. Completed current action presents the next live pending sibling stopped, and completing another row leaves the active target intact. A completed last step hides obsolete execution controls and leaves parent completion explicit. The clock still sums all appropriate segments and does not invent a session budget.

`POST /api/v1/tasks/{id}/decompose` is a contextual entry to the existing chat turn. It validates the open task and restricts both advertised and dispatched tools to list_steps/add_step for that exact ID, preserving existing/completed work and preventing unrelated mutations. The response is the existing chat Message, including error on provider failure; clients must inspect error/tool errors and refetch even after partial failure. Manual creation/editing remains available. Tests use the existing fake/scripted fixtures; no claim of live-provider validation is made.


The shared-interaction preview was updated on 2026-10-25 at 15:44 America/Sao_Paulo using `.demo/whatnow-preview-shared` on7790 with the same database and AI configuration. A new consistent backup preceded restart. Every row of tasks/steps/focus work/corrections/request tokens/selection/settings/messages compared identical afterward; integrity and foreign-key checks passed. Local and Tailscale HTTP200 HTML matched final SHA256 `0f71fd56c7ef74ae76200275abbf47468a2b9b351641e4941ad40396e329864f`. A disposable-database executable HTTP smoke covered step create/edit/reorder, completion without timer, timer preserved during inline-field edits, explicit step completion leaving parent open, and AI failure with manual operations still usable.164 frontend unit tests, Svelte check without errors/warnings, agent/API tests and vet passed. No browser/touch or live-provider validation was performed in this round; physical drag and UX remain real-use acceptance.
