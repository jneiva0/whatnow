import { api, ApiError, type Blocker, type CalendarEvent, type ExecutionState, type Pick, type Task, type Today } from './api';
import { available } from './execution';
import { today } from './today.svelte';
import { execution } from './execution.svelte';

export interface SuggestionSnapshot {
 day: Today;
 tasks: Task[];
 blockers: Blocker[];
 execution: ExecutionState;
 events: CalendarEvent[];
 revision: number;
}

/** Client request lifecycle only. The planner chooses; the service validates every command. */
export class SuggestionState {
 pick = $state.raw<Pick | null>(null);
 asking = $state(false);
 nothing = $state(false);
 problem = $state('');
 stale = $state(false);
 #snapshot: SuggestionSnapshot | null = null;
 #key = '';
 #generation = 0;
 #attempted = false;
 #wanted = false;
 #active = false;
 #interacting = false;
 #busy = false;
 #requiredRevision = 0;
 #captureContext = '';
 #retryRequest = {exclude: 0, context: ''};

 private request: (exclude: number, context: string) => Promise<Pick>;
 private refresh: () => Promise<void>;
 constructor(
  request = (exclude: number, context: string) => api.whatNow(exclude, '', context),
  refresh = async () => { await Promise.all([today.refresh(), execution.refresh()]); }
 ) { this.request = request; this.refresh = refresh; }

 get running() { return this.#snapshot?.execution.sessions.some((s) => s.state === 'running') ?? false; }

 sync(snapshot: SuggestionSnapshot, active: boolean, interacting = false, busy = false) {
  const dayChanged = this.#snapshot && snapshot.day.date !== this.#snapshot.day.date;
  if (dayChanged) {
   this.pick = null; this.nothing = false; this.problem = ''; this.#attempted = false;
   this.#wanted = false; this.#requiredRevision = 0; ++this.#generation;
  }
  const key = JSON.stringify([
   snapshot.day.date, snapshot.day.day?.pace, snapshot.day.plan,
   snapshot.tasks, snapshot.blockers, snapshot.execution.steps, snapshot.events,
   snapshot.execution.sessions.map((s) => [s.id, s.task_id, s.step_id, s.state]), snapshot.execution.selected_id
  ]);
  if (this.#key && key !== this.#key) {
   this.stale = true;
   if (this.asking) { ++this.#generation; this.#wanted = true; }
  }
  if (interacting && !this.#interacting && this.asking) { ++this.#generation; this.#wanted = true; }
  this.#snapshot = snapshot; this.#key = key;
  this.#active = active; this.#interacting = interacting; this.#busy = busy;
  if (this.running) {
   this.pick = null; this.#attempted = false; this.#wanted = false; ++this.#generation;
  } else if (this.pick && !this.valid(this.pick)) {
   this.pick = null; this.#wanted = true;
  }
  if (active && !this.#attempted) this.#wanted = true;
  this.#pump();
 }

 /** Events arrive before refetches. Wait for a confirmed day snapshot before asking again. */
 changed(automatic = false, captureContext = '') {
  if (automatic && captureContext) this.#captureContext = captureContext;
  this.stale = true;
  this.#requiredRevision = (this.#snapshot?.revision ?? 0) + 1;
  if (automatic || this.asking) this.#wanted = true;
  if (automatic && !this.asking) this.#retryRequest = {exclude: 0, context: this.#captureContext};
  ++this.#generation;
 }

 valid(pick: Pick) {
  const s = this.#snapshot;
  const task = s?.tasks.find((t) => t.id === pick.task_id);
  if (!s || !task || !available(task, s.tasks, s.blockers, s.day.date)) return false;
  return pick.step_id == null || s.execution.steps.some((step) => step.task_id === task.id && step.id === pick.step_id && !step.done && !step.deleted_at);
 }

 deactivate() { this.#active = false; this.#interacting = false; }

 dismiss() { this.pick = null; this.nothing = false; this.problem = ''; this.#attempted = true; this.#wanted = false; ++this.#generation; }

 #pump() {
  if (this.#wanted && this.#active && !this.#interacting && !this.#busy && !this.asking && !this.running && this.#snapshot && !!this.#snapshot.day.plan && !this.#snapshot.day.day?.wrapped_at && this.#snapshot.revision >= this.#requiredRevision) void this.ask(this.#retryRequest.exclude, this.#retryRequest.context);
 }

 async ask(exclude = 0, context = '') {
  if (this.asking || this.running || !this.#snapshot || this.#busy) return;
  this.#wanted = false; this.#attempted = true;
  if (!context && this.#captureContext) context = this.#captureContext;
  this.#captureContext = '';
  this.#retryRequest = {exclude, context};
  this.asking = true; this.problem = ''; this.nothing = false;
  const generation = ++this.#generation;
  try {
   const pick = await this.request(exclude, context);
   // The planner may name a step that arrived by SSE during the call.
   await this.refresh();
   if (generation !== this.#generation || this.running) return;
   if (!this.valid(pick)) {
    this.pick = null;
    this.problem = 'Essa ação mudou enquanto a sugestão era preparada. Atualize para escolher outra.';
    return;
   }
   this.pick = pick; this.stale = false;
   this.#retryRequest = {exclude: 0, context: ''};
  } catch (err) {
   if (generation !== this.#generation) return;
   if (err instanceof ApiError && err.status === 404) { this.pick = null; this.nothing = true; }
   else this.problem = `Não foi possível obter uma sugestão: ${err instanceof Error ? err.message : String(err)}. Tente novamente ou escolha na lista.`;
  } finally {
   this.asking = false;
   this.#pump();
  }
 }
}
export const suggestion = new SuggestionState();
