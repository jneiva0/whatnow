import { api, type Estimate, type ExecutionState, type FocusSession } from './api';
import { elapsedMs, FocusEstimateAlerts, type FocusEstimateNotice, orderedSteps, requestID } from './execution';
import { today } from './today.svelte';

export class Execution {
 view = $state.raw<ExecutionState | null>(null);
 now = $state(Date.now());
 #readProblem = $state('');
 #writeProblem = $state('');
 busy = $state(false);
 #skew = 0;
 #generation = 0;
 #connected = true;
 #refreshPending = false;
 #estimateAlerts = new FocusEstimateAlerts();
 #onEstimateReached: ((notice: FocusEstimateNotice) => void) | undefined;
 #poll: ReturnType<typeof setInterval> | undefined;
 #correctionRetry: {key: string; request: string} | undefined;
 #retry: {key: string; request: string} | undefined;

 get selected(): FocusSession | undefined {
  return this.view?.sessions.find((s) => s.id === this.view?.selected_id);
 }
 /** A successful background read cannot confirm a command whose response was lost. */
 get problem(): string { return this.#writeProblem || this.#readProblem; }
 steps(taskID: number) { return orderedSteps(this.view?.steps ?? [], taskID); }
 /** Saved totals are authoritative server facts, not a second client calculation. */
 estimate(taskID: number): Estimate { return this.view?.estimates?.[taskID] ?? {minutes: null, partial: false}; }
 elapsed(session: FocusSession): number { return elapsedMs(session, this.now); }

 #accept(view: ExecutionState, sent: number) {
  this.#skew = Date.parse(view.server_now) - (sent + Date.now()) / 2;
  this.now = Date.now() + this.#skew;
  this.view = view;
  this.#readProblem = '';
  this.#checkEstimate();
 }
 #checkEstimate() {
  if (!this.#onEstimateReached) return;
  const notice = this.#estimateAlerts.observe(this.view, this.now);
  if (notice) this.#onEstimateReached(notice);
 }
 async refresh(): Promise<void> {
  if (this.busy) { this.#refreshPending = true; return; }
  this.#refreshPending = false;
  const generation = ++this.#generation;
  const sent = Date.now();
  try {
   const view = await api.execution();
   if (generation === this.#generation) this.#accept(view, sent);
  } catch (err) {
   if (generation === this.#generation) this.#readProblem = `Não foi possível atualizar o foco: ${err instanceof Error ? err.message : String(err)}. Tente novamente.`;
  }
 }
 async start(taskID: number, stepID: number | null): Promise<boolean> {
  const key = `${taskID}:${stepID}`;
  if (!this.#retry || this.#retry.key !== key) this.#retry = {key, request: requestID()};
  const request = this.#retry.request;
  return this.#write(async () => {
   const result = await api.startFocus(taskID, stepID, request);
   this.#retry = undefined;
   return result;
  });
 }
 async correct(id: number, totalMS: number, reason: string): Promise<boolean> {
  const key = `${id}:${totalMS}:${reason}`;
  if (!this.#correctionRetry || this.#correctionRetry.key !== key) this.#correctionRetry = {key, request: requestID()};
  const request = this.#correctionRetry.request;
  return this.#write(async () => {
   const result = await api.correctFocus(id, totalMS, reason, request);
   this.#correctionRetry = undefined;
   return result;
  });
 }
 async action(id: number, action: 'pause' | 'resume' | 'end' | 'complete'): Promise<boolean> {
  return this.#write(() => api.focusAction(id, action));
 }
 async #write(op: () => Promise<ExecutionState>): Promise<boolean> {
  if (this.busy) return false;
  this.busy = true;
  this.#writeProblem = '';
  ++this.#generation;
  const sent = Date.now();
  try {
   const result = await op();
   ++this.#generation;
   this.#accept(result, sent);
   today.changed();
   return true;
  } catch (err) {
   this.#writeProblem = `A ação não foi confirmada: ${err instanceof Error ? err.message : String(err)}. Tente novamente.`;
   return false;
  } finally {
   this.busy = false;
   // Coalesce updates received during the write, including this client's own SSE
   // and changes from another device that happened after the mutation snapshot.
   if (this.#refreshPending) void this.refresh();
  }
 }
 connected(value: boolean) {
  this.#connected = value;
  this.#polling();
 }
 #polling() {
  clearInterval(this.#poll);
  this.#poll = undefined;
  if (!this.#connected && document.visibilityState === 'visible') {
   this.#poll = setInterval(() => {
    if (!this.busy) { void this.refresh(); today.changed(); }
   }, 5000);
  }
 }
 startLive(onEstimateReached: (notice: FocusEstimateNotice) => void = () => {}): () => void {
  this.#onEstimateReached = onEstimateReached;
  this.#estimateAlerts = new FocusEstimateAlerts();
  void this.refresh();
  const tick = setInterval(() => { this.now = Date.now() + this.#skew; this.#checkEstimate(); }, 1000);
  const caughtUp = () => { if (document.visibilityState === 'visible') void this.refresh(); this.#polling(); };
  document.addEventListener('visibilitychange', caughtUp);
  window.addEventListener('focus', caughtUp);
  window.addEventListener('online', caughtUp);
  return () => {
   this.#onEstimateReached = undefined;
   clearInterval(tick); clearInterval(this.#poll);
   document.removeEventListener('visibilitychange', caughtUp);
   window.removeEventListener('focus', caughtUp);
   window.removeEventListener('online', caughtUp);
  };
 }
}
export const execution = new Execution();
