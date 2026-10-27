# Graph Report - whatnow  (2026-10-24)

## Corpus Check
- 266 files · ~133,291 words
- Verdict: corpus is large enough that graph structure adds value.
- Unclassified: 14 file(s) not represented in the graph (top: .golden 8, (none) 5, .css 1)

## Summary
- 2261 nodes · 9017 edges · 100 communities (81 shown, 19 thin omitted)
- Extraction: 89% EXTRACTED · 11% INFERRED · 0% AMBIGUOUS · INFERRED: 1004 edges (avg confidence: 0.85)
- Token cost: 0 input · 0 output

## Community Hubs (Navigation)
- Shared Backend Dependencies
- Planner Test Infrastructure
- Terminal View Tests
- HTTP Request Handlers
- Morning and Evening Rituals
- Today Task Planning
- Service and Calendar Tests
- HTTP API Tests
- Integration Selection Interface
- Live Updates and Nudges
- Day Planning Domain
- Integration Contracts
- Task Editing and State
- Application Dependency Wiring
- Transactional Writes and Undo
- Settings and Onboarding
- Integration Runner Lifecycle
- Frontend Tooling Dependencies
- Chat and Capture State
- Terminal Input Routing
- Chat and Planning Context
- Integration Runner Tests
- Application Shell and Themes
- Day Timeline and Time
- Architecture and Design Rationale
- Terminal Timeline Rendering
- Terminal Interaction Tests
- Task Mutation Domain
- Agent Tool Execution
- Terminal Screen Rendering
- Agent Tool Tests
- Calendar Feed Parsing
- Product Guide and Demo
- CLI and Embedded Web
- API Client and Forms
- Eligibility and Fallback Tests
- Toasts and Data Reset
- Deterministic Planning Facts
- Model Client Construction
- Playwright Browser Flows
- Capture Cleanup Tests
- Model Planning Decisions
- Terminal Event Handling
- Terminal Ritual Text
- Morning Review and Triage
- Calendar Parser Tests
- Structured Model Requests
- Streaming Model Contracts
- Capture Cleanup Pipeline
- Activity and Database Scanning
- Terminal Snapshot Feedback
- Terminal Integration Views
- Terminal Calendar Tests
- What Now Interface
- Structured Output Schemas
- Integration Field Persistence
- Day Closing and Summaries
- Jira Integration Tests
- Review Queue Integration Tests
- Database and Settings Loading
- Core Database Schema
- Provider Settings Configuration
- Chat Task Breakdown Example
- Model Streaming Tests
- OpenAI Chat Adapter
- Anthropic Chat Adapter
- Cleanup Result Mapping
- Demo Clock Controls
- Evening Reflection Example
- Visual Theme Examples
- Today Planning Example
- Next Action Recommendation Example
- Wrap Triage Example
- Deterministic Model Fake
- Morning Capture Example
- Morning Plan Example
- Terminal Capture Fallback Tests
- TypeScript Compiler Configuration
- Chat Dashboard Example
- Embedded Frontend Design
- Undo Batch Actors
- Homebrew Package Installation
- Chat Storage Schema
- Terminal Celebration Animation
- Day and Plan Schema
- Integration Secret Storage Design
- Reset and Testing Design
- Integration Pull Results
- Settings Storage Schema
- Integration Storage Schema
- Application Type Declarations
- Static Frontend Rendering
- SQLite Storage Design
- Notification Timing Design
- Go Module Identity
- Svelte Development Guidance
- Application Favicon
- Search Crawler Policy

**Token accounting:** actual subagent usage was unavailable; zero figures below are placeholders.

## God Nodes (most connected - your core abstractions)
1. `newFixture()` — 101 edges · [internal/tui/helpers_test.go](/Users/joao/space/whatnow/internal/tui/helpers_test.go)
2. `Model` — 84 edges · [internal/tui/model.go](/Users/joao/space/whatnow/internal/tui/model.go)
3. `Task` — 66 edges · [internal/domain/domain.go](/Users/joao/space/whatnow/internal/domain/domain.go)
4. `send()` — 62 edges · [internal/tui/view_test.go](/Users/joao/space/whatnow/internal/tui/view_test.go)
5. `newFixture()` — 60 edges · [internal/planner/helpers_test.go](/Users/joao/space/whatnow/internal/planner/helpers_test.go)

## Surprising Connections (you probably didn't know these)
- `init()` --calls--> `NewClient()`  [EXTRACTED]
  cmd/whatnow/e2e.go → internal/settings/keys.go
- `open()` --references--> `App`  [EXTRACTED]
  cmd/whatnow/main.go → internal/app/app.go
- `open()` --calls--> `Open()`  [EXTRACTED]
  cmd/whatnow/main.go → internal/app/app.go
- `newServer()` --references--> `App`  [EXTRACTED]
  cmd/whatnow/main.go → internal/app/app.go
- `SvelteKit HTML shell with hover preloading` --implements--> `Static frontend build and SPA fallback`  [INFERRED]
  web/src/app.html → docs/superpowers/specs/2026-10-05-whatnow-web-ui-today-design.md

## Import Cycles
- None detected.

## Hyperedges (group relationships)
- **User-requested morning integration pulls** — readme_morning_ritual, readme_google_calendar, readme_jira, readme_review_queue [EXTRACTED 1.00]
- **Energy-sensitive personalized replanning** — demo_personalized_energy_plan, readme_about_me, readme_pace, readme_three_tier_planning [EXTRACTED 1.00]
- **Suggested report steps become tasks following an explicit user request** — docs_screenshots_chat_draft_q3_report, docs_screenshots_chat_assisted_task_decomposition, docs_screenshots_chat_user_requested_task_creation, docs_screenshots_chat_created_report_tasks, docs_screenshots_chat_undo_task_creation [EXTRACTED 1.00]
- **Three-tier daily priorities** — docs_screenshots_morning_plan_must, docs_screenshots_morning_plan_should, docs_screenshots_morning_plan_could [EXTRACTED 1.00]
- **Six visual themes share a timeline, task capture, Must/Should/Could lists, and chat panel** — docs_screenshots_themes_sticker, docs_screenshots_themes_dream_of_the_90s, docs_screenshots_themes_neon_nights, docs_screenshots_themes_terminal, docs_screenshots_themes_paper, docs_screenshots_themes_vaporwave, docs_screenshots_themes_task_planning_interface [INFERRED 0.95]
- **Each task offers Roll over, Push, and Drop** — docs_screenshots_wrap_triage_draft_q3_report, docs_screenshots_wrap_triage_renew_passport, docs_screenshots_wrap_triage_groceries, docs_screenshots_wrap_triage_roll_over, docs_screenshots_wrap_triage_push, docs_screenshots_wrap_triage_drop [EXTRACTED 1.00]
- **REST, chat, and integration changes share service batches** — docs_superpowers_specs_2026_10_05_whatnow_design_service, docs_superpowers_specs_2026_10_05_whatnow_design_httpapi, docs_superpowers_specs_2026_10_05_whatnow_design_agent, docs_superpowers_specs_2026_10_12_whatnow_integrations_design_sync_events [INFERRED 0.95]
- **Calendar and ticket integrations implement the compiled contract** — docs_superpowers_specs_2026_10_12_whatnow_integrations_design_contract, docs_superpowers_specs_2026_10_12_whatnow_integrations_design_gcal, docs_superpowers_specs_2026_10_12_whatnow_jira_design_jira, docs_superpowers_specs_2026_10_15_whatnow_review_queue_design_review_queue [INFERRED 0.95]
- **Shared application construction enables isolated model testing** — docs_superpowers_specs_2026_10_06_whatnow_e2e_design_fake_provider, docs_superpowers_specs_2026_10_06_whatnow_e2e_design_build_tag, docs_superpowers_specs_2026_10_06_whatnow_tui_design_app_open [INFERRED 0.95]

## Communities (100 total, 19 thin omitted)

### Community 0 - "Shared Backend Dependencies"
Cohesion: 0.15
Nodes (20): Delta, taskArgs, Failed, ConnectionResult, errorBody, Mutation, searchResult, queueResult (+12 more)

### Community 1 - "Planner Test Infrastructure"
Cohesion: 0.05
Nodes (82): TestACaptureReportsAFailureWhenSettingsCantBeRead(), TestOpenCleansUpCapturesInTheBackground(), TestOpenPutsTheDatabaseInTheUserConfigDir(), TestSchemaPropertiesComeInWritingOrder(), TestOnlyRawCapturesTriggerCleanup(), TestStreamDeliversWrites(), TestAPlainRequestGetsTheFixedSentence(), TestAReflectionRequestGetsTheFixedReflection() (+74 more)

### Community 2 - "Terminal View Tests"
Cohesion: 0.06
Nodes (75): TestAFailedWriteShowsTheError(), TestAToastLeavesWhenItsTimeIsUp(), TestAWriteOnAClosedDatabaseFails(), TestEscCancelsAnAdd(), TestUndoingSomethingChangedSinceShowsWhy(), TestUndoWithNothingToUndoSaysSo(), TestATooLongChatMessageShowsTheLimit(), TestChatOpensBesideTheListOnAWideTerminal() (+67 more)

### Community 3 - "HTTP Request Handlers"
Cohesion: 0.12
Nodes (18): Server, SettingsView, decode(), hostname(), isLoopback(), pathID(), writeError(), writeJSON() (+10 more)

### Community 4 - "Morning and Evening Rituals"
Cohesion: 0.05
Nodes (52): Action, Review, Today, Triage, { count, at }, Props, actions, Props (+44 more)

### Community 5 - "Today Task Planning"
Cohesion: 0.05
Nodes (40): Day, Task, Tier, { day, done, meetings }, Props, { done, meetings }, Props, { card, tier, date, listed = false } (+32 more)

### Community 6 - "Service and Calendar Tests"
Cohesion: 0.08
Nodes (60): TestAddAndRemoveBlocker(), TestAddBlockerRejectsCycles(), TestAddBlockerTwiceRecordsOnce(), TestBlockerErrors(), TestListBlockersHidesDeletedTasks(), TestPlanReasonsListsOneTriggersReasonsOldestFirst(), TestPlansAreVersionedPerDay(), TestSavePlanRejectsBadTierOrTrigger() (+52 more)

### Community 7 - "HTTP API Tests"
Cohesion: 0.11
Nodes (48): modelFunc, testAPI, TestChatFinishesAfterTheClientLeaves(), TestChatRepliesAndKeepsHistory(), TestChatWithoutAProviderSaysSo(), TestClockIsTheServersTime(), TestMoveTask(), TestPushTask() (+40 more)

### Community 8 - "Integration Selection Interface"
Cohesion: 0.07
Nodes (38): Field, Integration, IntegrationUpdate, PulledEvent, PullResult, Ticket, form, { integration = $bindable() } (+30 more)

### Community 9 - "Live Updates and Nudges"
Cohesion: 0.07
Nodes (23): Plan, Settings, celebration, open, today, cheer(), tokens, deliver() (+15 more)

### Community 10 - "Day Planning Domain"
Cohesion: 0.08
Nodes (10): wrapped, Day, Pace, Plan, Trigger, Started, started, Today (+2 more)

### Community 11 - "Integration Contracts"
Cohesion: 0.07
Nodes (13): Calendar, board, fake, Kind, Config, Field, Pulls, Source (+5 more)

### Community 12 - "Task Editing and State"
Cohesion: 0.08
Nodes (28): Blocker, Status, task, addBlocker(), adding, blockedBy, candidates, form (+20 more)

### Community 13 - "Application Dependency Wiring"
Cohesion: 0.11
Nodes (27): Deps, Agent, New(), fixture, App, Open(), open(), TestCloseClosesTheDatabase() (+19 more)

### Community 14 - "Transactional Writes and Undo"
Cohesion: 0.12
Nodes (22): eventMutation, Blocker, Event, waitsOn(), writeDay(), EventFields, liveEvent(), writeEvent() (+14 more)

### Community 15 - "Settings and Onboarding"
Cohesion: 0.06
Nodes (31): ConnectionResult, Provider, Props, { style = $bindable(), aboutMe = $bindable() }, permission, key, keyNote, provider (+23 more)

### Community 16 - "Integration Runner Lifecycle"
Cohesion: 0.15
Nodes (11): Integration, onList, RowView, saved, config(), Runner, View, rowSecrets() (+3 more)

### Community 17 - "Frontend Tooling Dependencies"
Cohesion: 0.07
Nodes (36): canvas-confetti, svelte-check, @sveltejs/adapter-static, @sveltejs/kit, @sveltejs/vite-plugin-svelte, @types/canvas-confetti, typescript, vite (+28 more)

### Community 18 - "Chat and Capture State"
Cohesion: 0.08
Nodes (11): BatchApplied, Message, ToolUse, chat, Draft, text, Handlers, Source (+3 more)

### Community 19 - "Terminal Input Routing"
Cohesion: 0.13
Nodes (3): newDump(), TestCycleWrapsAroundTheActions(), mode

### Community 20 - "Chat and Planning Context"
Cohesion: 0.08
Nodes (20): history, errorText(), joinText(), addDays(), foldPoint(), fixture, toLLM(), ChatFold (+12 more)

### Community 21 - "Integration Runner Tests"
Cohesion: 0.18
Nodes (26): env, PullResult, at(), event(), newEnv(), offered(), pulled(), TestAddMakesTasksInOneBatchThatOneUndoRemoves() (+18 more)

### Community 22 - "Application Shell and Themes"
Cohesion: 0.08
Nodes (14): svelte, Theme, log, { opening }, Props, text, Props, { theme = $bindable() } (+6 more)

### Community 23 - "Day Timeline and Time"
Cohesion: 0.09
Nodes (27): CalendarEvent, paceLabels, { events, now, endOfDay, startedAt }, layout, Props, canWrap, date, menuOpen (+19 more)

### Community 24 - "Architecture and Design Rationale"
Cohesion: 0.07
Nodes (27): agent: serialized streaming chat loop, whatnow local executive-function assistant architecture, cleanup: durable capture followed by structured extraction, facts: pure constraint and annotation functions, httpapi: REST and SSE under /api/v1, llm: Anthropic and OpenAI-compatible provider interface, planner: versioned plans and current-state task picks, prompt: ordered persona and live-context layers (+19 more)

### Community 25 - "Terminal Timeline Rendering"
Cohesion: 0.13
Nodes (27): TestMain(), load(), eventEnd(), floorHour(), hourLabel(), layoutStrip(), paint(), at() (+19 more)

### Community 26 - "Terminal Interaction Tests"
Cohesion: 0.16
Nodes (26): fixture, quit(), TestActionsDoNothingOnAnEmptyList(), TestAddCapturesAndCleanupSplitsIt(), TestPlusAndMinusShiftThePaceAndUUndoesIt(), TestPPushesToTomorrow(), TestSpaceMarksDoneAndUUndoesIt(), TestThePaceStopsAtItsEnds() (+18 more)

### Community 27 - "Task Mutation Domain"
Cohesion: 0.10
Nodes (19): brief, change, taskMutation, Energy, Ref, Status, Task, decode() (+11 more)

### Community 28 - "Agent Tool Execution"
Cohesion: 0.15
Nodes (24): eventBrief, tool, turn, addBlocker(), args(), blockerArgs(), briefOf(), createEvent() (+16 more)

### Community 29 - "Terminal Screen Rendering"
Cohesion: 0.13
Nodes (7): fireworks(), paceRow(), chip(), sticker(), wrapDue(), Model, toast

### Community 30 - "Agent Tool Tests"
Cohesion: 0.15
Nodes (30): TestAFailedToolGoesBackToTheModelAndEarlierWritesStay(), TestAnErrorAfterWritesKeepsThem(), TestAReplyWhoseWritesAllFailedHasNoBatch(), TestAReplyWithNoTextIsNotReplayed(), TestAReplyWithoutWritesHasNoBatch(), TestATurnsWritesShareOneBatchAndStream(), TestMessagesAreChecked(), TestOneTurnAtATime() (+22 more)

### Community 31 - "Calendar Feed Parsing"
Cohesion: 0.13
Nodes (14): parser, Clock, email(), inWindow(), isDate(), occurrenceKey(), owner(), parse() (+6 more)

### Community 32 - "Product Guide and Demo"
Cohesion: 0.08
Nodes (26): Report avoidance and subtask creation, A day with whatnow demo script, Isolated demo database and port, Manually advanced demo clock, Small task before standup, Wrap-up reflection and next-morning reconciliation, Passport rescheduling and undo, About me planning context (+18 more)

### Community 33 - "CLI and Embedded Web"
Cohesion: 0.10
Nodes (17): command(), listenUnlessTaken(), main(), newServer(), runTUI(), serve(), TestCommand(), TestListenUnlessTakenListensOnAFreeAddress() (+9 more)

### Community 34 - "API Client and Forms"
Cohesion: 0.10
Nodes (22): Activity, ApiError, configured(), Energy, FieldKind, IntegrationRow, KeySource, NudgeToggles (+14 more)

### Community 35 - "Eligibility and Fallback Tests"
Cohesion: 0.19
Nodes (25): span, Size, annotations(), at(), date(), deleted(), due(), event() (+17 more)

### Community 36 - "Toasts and Data Reset"
Cohesion: 0.14
Nodes (9): Mutation, Source, deleting, typed, cleanupText(), count(), Timer, Toast (+1 more)

### Community 37 - "Deterministic Planning Facts"
Cohesion: 0.14
Nodes (21): Input, TaskFacts, Compute(), dayEnd(), daysBetween(), due(), Duration(), end() (+13 more)

### Community 38 - "Model Client Construction"
Cohesion: 0.14
Nodes (16): init(), NewAnthropic(), anthropicClient(), openAIClient(), TestChatProviderErrorIsReturned(), anthropicReply(), request(), TestAnthropicFallbackOptIn() (+8 more)

### Community 39 - "Playwright Browser Flows"
Cohesion: 0.22
Nodes (10): @playwright/test, fake, ok(), setUp(), startDay(), test, feed(), icalTime() (+2 more)

### Community 40 - "Capture Cleanup Tests"
Cohesion: 0.27
Nodes (17): fixture, newFixture(), reply(), task(), TestABadReplyChangesNothing(), TestAnnouncesTheWorkAroundTheModelCall(), TestANoteAboutAnOpenTicketJoinsItsTask(), TestASplitOutTicketJoinsItsTask() (+9 more)

### Community 41 - "Model Planning Decisions"
Cohesion: 0.15
Nodes (12): PlanEntry, Tier, checkPlan(), fallbackPlan(), ask(), parse(), checkPick(), Pick (+4 more)

### Community 42 - "Terminal Event Handling"
Cohesion: 0.21
Nodes (5): Message, friendly(), TestFriendlyErrors(), tick(), hubMsg

### Community 43 - "Terminal Ritual Text"
Cohesion: 0.12
Nodes (13): PulledEvent, eventLine(), count(), TestWinLine(), winLine(), names(), ticketLine(), pushLabel() (+5 more)

### Community 44 - "Morning Review and Triage"
Cohesion: 0.15
Nodes (14): Morning, checkSuggestions(), Action, Review, payload(), TestPayloadSetsUntilOnlyForAPush(), TestTriageRows(), triageRows() (+6 more)

### Community 45 - "Calendar Parser Tests"
Cohesion: 0.31
Nodes (16): New(), day(), feed(), mustParse(), mustZone(), summary(), TestParseDropsWhatTheOwnerDeclined(), TestParseExpandsASeriesAroundItsExceptions() (+8 more)

### Community 46 - "Structured Model Requests"
Cohesion: 0.16
Nodes (6): Request, gate, hang, leaving, noChat, script

### Community 47 - "Streaming Model Contracts"
Cohesion: 0.21
Nodes (7): blocking, model, replyFunc, Client, chatModel, ChatRequest, Reply

### Community 48 - "Capture Cleanup Pipeline"
Cohesion: 0.24
Nodes (4): ClientFunc, tickets, Runner, resolve()

### Community 49 - "Activity and Database Scanning"
Cohesion: 0.24
Nodes (9): Activity, scanDay(), scanEvent(), batchActivity(), parseNullTime(), parseTime(), scanTask(), querier (+1 more)

### Community 50 - "Terminal Snapshot Feedback"
Cohesion: 0.19
Nodes (9): BatchApplied, Source, cheerFor(), applied(), cheerSnap(), TestCheerFor(), cheer, snapshot (+1 more)

### Community 51 - "Terminal Integration Views"
Cohesion: 0.17
Nodes (12): calendar, integrationsMsg, pickTicket, pulledFrom, pulledMsg, pullFailedMsg, removedMsg, tickets (+4 more)

### Community 52 - "Terminal Calendar Tests"
Cohesion: 0.34
Nodes (10): fixture, TestABusyCalendarScrollsToTheCursor(), TestAFailedPullStaysOnTheStep(), TestCalendarView(), TestEscCancelsTheSkipQuestion(), TestJustTodayLeavesTheSeries(), TestTabSkipsThePull(), TestTheCursorStaysOnTheListAndRPullsAgain() (+2 more)

### Community 53 - "What Now Interface"
Cohesion: 0.15
Nodes (8): api, Pick, explaining, { pick, task, asking, onnotthis, onclose }, Props, reason, TodayState, day

### Community 54 - "Structured Output Schemas"
Cohesion: 0.22
Nodes (8): taskFields(), Enum(), Integer(), Object(), Prop(), TestObjectKeepsFieldOrderAndRequiresEveryField(), Field, fields

### Community 55 - "Integration Field Persistence"
Cohesion: 0.18
Nodes (9): Row, RowUpdate, Skip, Update, field(), listRows(), DeleteSecrets(), newRowID() (+1 more)

### Community 56 - "Day Closing and Summaries"
Cohesion: 0.27
Nodes (4): checkStale(), Triage, writeTitles(), fallbackNote()

### Community 57 - "Jira Integration Tests"
Cohesion: 0.40
Nodes (9): New(), cfg(), pull(), server(), TestPullBlamesTheNetworkWhenJiraIsOutOfReach(), TestPullExplainsWhatWentWrongWithoutTheToken(), TestPullGivesUpOnASlowJira(), TestPullOffersTicketsMineFirst() (+1 more)

### Community 58 - "Review Queue Integration Tests"
Cohesion: 0.49
Nodes (10): New(), cfg(), pull(), server(), TestABlankMinimumKeepsNegativeScores(), TestALongQueueIsCappedAndCountsTheRest(), TestATimeoutBlamesGitHubNotReviewQueue(), TestMinScoreCutsTheLowerPulls() (+2 more)

### Community 59 - "Database and Settings Loading"
Cohesion: 0.25
Nodes (11): Defaults(), Load(), TestLoadDefaultsThenSaveRoundTrip(), TestLoadFillsFieldsAddedSinceTheSave(), TestSettingsSavedBeforeEndOfDayExistedGetTheDefault(), TestValidate(), Open(), TestForeignKeysEnforced() (+3 more)

### Community 60 - "Core Database Schema"
Cohesion: 0.27
Nodes (9): activity, activity_batch, activity_entity, batches, blockers, events, events_starts_at, tasks (+1 more)

### Community 61 - "Provider Settings Configuration"
Cohesion: 0.33
Nodes (7): Options, fakeClient(), ProviderConfig, Settings, Save(), Nudges, Theme

### Community 62 - "Chat Task Breakdown Example"
Cohesion: 0.22
Nodes (9): Chat confirms tasks #8–#11: open report, make heading list, paste scraps, write three Wins bullets, Draft Q3 report: Must priority, M size, high energy, and expressed dread, Groceries: Could priority, S size, low energy, Must, Should, Could, and Unplanned task groups, Renew passport before November trip: Could priority, M size, medium energy, Replan control, Review Davey's API doc: Should priority, M size, medium energy, Undo control on chat task-creation confirmation (+1 more)

### Community 63 - "Model Streaming Tests"
Cohesion: 0.31
Nodes (9): anthropicEvent(), anthropicStream(), chatRequest(), collect(), openAIChunk(), TestAnthropicChatStreamsTextAndToolCalls(), TestChatStopReasonsBecomeErrors(), TestOpenAIChatStreamsTextAndToolCalls() (+1 more)

### Community 64 - "OpenAI Chat Adapter"
Cohesion: 0.22
Nodes (5): system(), Block, openAIMessages(), openAISystem(), openAIClient

### Community 65 - "Anthropic Chat Adapter"
Cohesion: 0.29
Nodes (3): anthropicMessage(), anthropicSystem(), anthropicClient

### Community 66 - "Cleanup Result Mapping"
Cohesion: 0.29
Nodes (6): blocker, event, newTask, original, result, taskOut

### Community 67 - "Demo Clock Controls"
Cohesion: 0.39
Nodes (7): demoClock(), readDemoClock(), TestDemoClockFallsBackToTheRealTimeWithoutAFile(), TestReadDemoClockRejectsATimeItCantParse(), TestReadDemoClockRunsOnFromTheTimeSet(), writeClock(), open()

### Community 68 - "Evening Reflection Example"
Cohesion: 0.29
Nodes (7): Completed work: Dana reply, checkout test fix, and Davey’s API document review, End-of-day narrative reflection, Goodnight continuation button, Rolled-over tasks: report draft, passport renewal, and groceries, whatnow Goodnight screen, Tomorrow’s plan: coffee, Wins bullets, report first, and nothing scary after 3, whatnow — executive function for the rest of us

### Community 69 - "Visual Theme Examples"
Cohesion: 0.25
Nodes (8): Dream of the 90s, Neon Nights, Paper, Sticker, Task planning and chat interface, Terminal, Whatnow theme comparison, Vaporwave

### Community 70 - "Today Planning Example"
Cohesion: 0.36
Nodes (7): Daily meeting timeline, Draft Q3 report, Must, Should, Could priorities, Quick capture and add task, Task energy and size, Today planning view, Unplanned tasks and replan

### Community 71 - "Next Action Recommendation Example"
Cohesion: 0.36
Nodes (7): Conversational day guidance, Whatnow daily dashboard, Reply to Dana recommendation, Must and Should task groups, Quick capture and task creation, Daily meeting timeline, Task energy and emotional context

### Community 72 - "Wrap Triage Example"
Cohesion: 0.25
Nodes (7): Draft Q3 report, Drop, Groceries, Push, Renew passport, Roll over, Unfinished-task wrap-up triage

### Community 73 - "Deterministic Model Fake"
Cohesion: 0.32
Nodes (7): fact, cleanup(), eligible(), pick(), plan(), suggestions(), task()

### Community 74 - "Morning Capture Example"
Cohesion: 0.38
Nodes (5): Mixed commitments, errands, emotional context, and someday aspirations, Skip and Next navigation controls, Three-position progress indicator with first position active, Whatnow morning dump screen, Whatnow executive function assistant

### Community 75 - "Morning Plan Example"
Cohesion: 0.40
Nodes (5): Could: passport renewal and groceries, Must: draft Q3 report, Morning plan screen, Should: test, API review, conference reply, Let's go action

### Community 76 - "Terminal Capture Fallback Tests"
Cohesion: 0.33
Nodes (3): fixture, newFixtureWith(), TestTheDumpIsPlannedRawWhenCleanupFails()

### Community 77 - "TypeScript Compiler Configuration"
Cohesion: 0.33
Nodes (5): $app/tsconfig, compilerOptions, strict, extends, include

### Community 78 - "Chat Dashboard Example"
Cohesion: 0.40
Nodes (5): Daily schedule showing Standup at 9:30 AM and 1:1 with Priya at 2:00 PM, Whatnow daily planner with chat sidebar, Collapsed progress summary: two things done and one meeting survived, Capture anything input and Add Task button, What now? action beside steady pace and next-meeting countdown

### Community 79 - "Embedded Frontend Design"
Cohesion: 0.50
Nodes (4): Single Go binary with embedded static SvelteKit frontend, Static frontend build and SPA fallback, SvelteKit HTML shell with hover preloading, Pre-paint persisted theme bootstrap

### Community 82 - "Chat Storage Schema"
Cohesion: 0.67
Nodes (3): chat_folds, messages, messages_date

## Ambiguous Edges - Review These
- `Separate planner summary and user reflection` → `Web morning and evening ritual pages`  [AMBIGUOUS]
  docs/superpowers/specs/2026-10-06-whatnow-web-ui-rituals-design.md · relation: references

## Knowledge Gaps
- **64 isolated node(s):** `github.com/jameynakama/whatnow`, `fixture`, `fixture`, `$app/tsconfig`, `Google Calendar secret iCal import` (+59 more)
  These have ≤1 connection - possible missing edges or undocumented components. (Counts symbols only; 466 node(s) total have ≤1 connection when file, concept and rationale nodes are included.)
- **19 thin communities (<3 nodes) omitted from report** — run `graphify query` to explore isolated nodes.

## Suggested Questions
_Questions this graph is uniquely positioned to answer:_

- **How does the terminal Model connect input handling, planning, shared services, and screen updates?**
  _High betweenness centrality (0.041) - this node is a cross-community bridge._
- **Are the 97 inferred relationships involving `newFixture()` (e.g. with `TestActionsDoNothingOnAnEmptyList()` and `TestAddCapturesAndCleanupSplitsIt()`) actually correct?**
  _`newFixture()` has 97 INFERRED edges - extractor-inferred connections that need verification._
- **What connects `github.com/jameynakama/whatnow`, `fixture`, `fixture` to the rest of the system?**
  _64 weakly-connected nodes found - possible documentation gaps or missing edges._
- **Should `Planner Test Infrastructure` be split into smaller, more focused modules?**
  _Cohesion score 0.04968789013732834 - nodes in this community are weakly interconnected._

### Low-confidence Hints
_AMBIGUOUS edges — the extractor was unsure. Verify before acting on these._

- **What is the exact relationship between `Separate planner summary and user reflection` and `Web morning and evening ritual pages`?**
  _Edge tagged AMBIGUOUS (relation: references) - confidence is low._

## Graph health warning

836 dangling-endpoint edges; 2 self-loop edges; 283 collapsed undirected edges. The map is usable for orientation but may omit or merge relationships. Raw relationships are preserved in extraction.json.


## Extraction coverage and limitations

- Mapped the current working tree, including uncommitted edits. This is static analysis, not runtime validation.
- All 241 detected code files have file nodes; 12 yielded no extracted symbols. Constants, configuration, and test-only files can have no symbols under the extractor's model.
- SQL support was installed and extraction rerun to include all eight detected SQL files.
- Sensitive-name filtering excluded web/.npmrc and web/src/lib/tokens.css; the latter appears to be a filename false positive and remains excluded.
- Semantic extraction used the host session's subagents. Their tool results do not expose token usage. Any 0 token figures above and in cost.json are placeholders, not measured cost or free usage.
- graphify-out/ was written into the working tree and is currently untracked. Add it to .gitignore or commit only outputs the team should share: https://docs.graphify.com/guides/team-workflows.
- See coverage.json and graph-health.json for machine-readable audit details.
