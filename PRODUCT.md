# Whatnow

<!-- impeccable:product-schema 1 -->

## Platform

web — the same responsive application serves desktop and mobile.

## Users

João uses Whatnow on his computer and phone to capture, plan, and decide what to do. Starting a concrete action is an additional job the product must support.

## Product Purpose

Keep capture, planning, chat, and replanning together with a lightweight execution experience: optional flat task steps and a persistent focus timer. A simple task can start directly.

## Operating Context

Go, SQLite, REST, SSE, and Svelte are established constraints. Mobile focus is essential. Simple synchronization between open clients is sufficient; there is no Passos integration or stack migration.

## Capabilities and Constraints

- Tasks remain the planning unit, under Must / Should / Could / Unplanned. Preserve S/M/L, energy, feelings, original capture, blockers, events, chat, and undo.
- Steps have descriptions, optional minute estimates, order, and explicit completion. One level only; decomposition is on demand, manually or through chat.
- One focus session runs at a time. Starting another pauses the previous session and preserves its time. Selected focus is distinct from previous paused sessions.
- An action has an optional estimate and accumulated actual time across all its work segments. Remaining time is estimate minus actual time, with explicit overtime after zero; without an estimate show elapsed time. There is no independent focus duration or default 15-minute budget. João explicitly rejected that assistant-introduced concept during QA on 2026-10-25.
- A step can be completed or reopened directly without starting a timer. Completing a step never completes its parent. Completing the task is explicit. Pausing, resuming, ending, and replanning preserve actual time.
- Actual focus time can be corrected explicitly through the UI or chat, with its own audit entry independent of generic undo. Changing an estimate does not change real time; a correction during a running segment must not double-count it. The agent can edit both task and step estimates.
- Day cards expose task time, collapsible steps, direct completion/reopen and start controls. Starting a task opens its first pending step when present; completing the current step presents the next pending step stopped, with no automatic start or invented action after the last step. Focus keeps readable task context, details, editing and the full step list in one continuous flow. Time is the primary focus information on desktop and mobile: a large estimate/actual counter with a progress indicator, plus a sticky current-action/timer/play-pause region while scrolling. Ver dia navigates without pausing; starting another action pauses the prior one.
- The server owns focus state. Refresh, returning to a suspended tab, SSE updates, and disconnected visible-page polling recover that state. Manual execution works without AI.

## Shared interaction — 2026-10-25

One task-step interaction serves Day, Focus and the editor: checkbox completion/reopen without a timer; text edits inline; estimated time at the right edits inline; play/pause icon; grip drag by mouse/touch with keyboard-arrow alternative; compact inline creation and discrete contextual remove/split. A separate parallel step editor is not part of the product. Compact time reads actual / estimated in MM:SS; absent estimate is a dash, not infinity.

A visible contextual “Dividir em passos com IA” operates on the selected task through the existing agent, preserving completed/existing steps. It differs from manual add and “Dividir este passo”. Provider/tool failures remain visible and do not block manual editing. Running focus, viewed content and next-ready action are different states: the actual running action owns the main timer, and text editing another row must never mask it. Completing the current action presents the next stopped; completing another row must leave the active focus in place.

## Confirmed UI refinements — 2026-10-25

Energy, size and free-text feeling belong immediately below the task title in both Day and Focus. Their displayed value is also the edit control: dropdowns for energy/size, free-text editing for feeling. Do not move these into a large lower information form. A discrete pencil may open the full existing editor without taking title width on mobile.

Play/pause is a single icon toggle with an accessible name and touch target. Task and step checkboxes already provide completion; avoid duplicated large completion buttons. Day uses a compact recognizable focus entry and compact add-step affordance. Focus retains task/step context while emphasizing time; the earlier problem was the wrong counter and lost context, not simply a large timer.

The annotated wireframe was input for these refinements, not wholesale approval of its layout. João authorized implementation after giving these corrections; no new image approval or visual tournament is needed. Visual/functional QA remains his real-use flow.

## Brand Commitments

Preserve Whatnow's existing visual identity and language. This work extends the running product; it is not a full redesign. Passos and Sunsama inform execution behavior, not a mandatory replacement palette.

## Delivery

User explicitly chose a navigable implementation directly in Whatnow (code-led), confirmed Impeccable initialization, and authorized necessary technical decisions. No production deployment is authorized. Real-model verification must be distinguished from fake-model tests; viewport testing must be distinguished from a physical phone.

## Evidence

Confirmed implementation brief delegated on 2026-10-25. Existing repository supplies the established implementation and visual system. No new audience, commercial claim, clinical claim, or deployment claim is inferred.
