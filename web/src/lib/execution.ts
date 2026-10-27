import type { Estimate, FocusSession, Step, Task, Blocker } from './api';

/** Derived from server timestamps, so sleeping tabs need no missed ticks replayed. */
export function elapsedMs(session: FocusSession, serverNow: number): number {
 const segment = session.state === 'running' && session.segment_started_at
  ? Math.max(0, serverNow - Date.parse(session.segment_started_at)) : 0;
 return Math.max(0, session.accumulated_ms + segment);
}

export function timerText(ms: number): string {
 const seconds = Math.max(0, Math.floor(ms / 1000));
 return `${Math.floor(seconds / 60).toString().padStart(2, '0')}:${(seconds % 60).toString().padStart(2, '0')}`;
}

export function orderedSteps(steps: Step[], taskID: number): Step[] {
 return steps.filter((s) => s.task_id === taskID && !s.deleted_at).sort((a, b) => a.position - b.position || a.id - b.id);
}

export function available(task: Task, tasks: Task[], blockers: Blocker[], date: string): boolean {
 return task.status === 'open' && (!task.not_before || task.not_before <= date)
  && !blockers.some((b) => b.task_id === task.id && tasks.some((t) => t.id === b.blocked_by_id && t.status === 'open'));
}

/** randomUUID requires a secure context; private-network HTTP still has getRandomValues. */
export function requestID(): string {
 if (typeof crypto.randomUUID === 'function') return crypto.randomUUID();
 const bytes = crypto.getRandomValues(new Uint8Array(16));
 bytes[6] = (bytes[6] & 15) | 64;
 bytes[8] = (bytes[8] & 63) | 128;
 const hex = Array.from(bytes, (b) => b.toString(16).padStart(2, '0')).join('');
 return `${hex.slice(0, 8)}-${hex.slice(8, 12)}-${hex.slice(12, 16)}-${hex.slice(16, 20)}-${hex.slice(20)}`;
}

/** Exact action totals include every segment, including older ended sessions. */
export function actionActual(sessions: FocusSession[], taskID: number, stepID: number | null, now: number): number {
 return sessions.filter((s) => s.task_id === taskID && s.step_id === stepID).reduce((sum, s) => sum + elapsedMs(s, now), 0);
}

/** Each stored segment belongs to the parent exactly once, whether direct or on a step. */
export function taskActual(sessions: FocusSession[], taskID: number, now: number): number {
 return sessions.filter((s) => s.task_id === taskID).reduce((sum, s) => sum + elapsedMs(s, now), 0);
}

/** Preview only: an unsaved editor draft cannot use the server's saved total yet.
 * Saved task surfaces read execution.estimate instead; Go validates on save. */
export function draftTaskEstimate(task: Task, steps: Step[]): Estimate {
 const live = orderedSteps(steps, task.id);
 const estimated = live.filter((s) => s.estimate_minutes != null);
 const sum = estimated.reduce((sum, s) => sum + s.estimate_minutes!, 0);
 return {minutes: task.estimate_minutes != null ? Math.max(task.estimate_minutes, sum) : estimated.length ? sum : null, partial: estimated.length < live.length};
}

export function actionClock(estimateMinutes: number | null | undefined, actual: number): {text: string; label: string} {
 if (estimateMinutes == null) return {text: timerText(actual), label: 'Tempo dedicado'};
 const remaining = estimateMinutes * 60000 - actual;
 return {text: `${remaining < 0 ? '+' : ''}${timerText(Math.abs(remaining))}`, label: remaining < 0 ? 'Tempo excedido' : 'Tempo restante'};
}

export interface ExecutionTarget { taskID: number; stepID: number | null; ready: boolean }
/** Running always owns the timer. Completed selected targets advance only to a stopped presentation. */
export function focusTarget(view: import('./api').ExecutionState | null, requested: {taskID: number; stepID: number | null} | null): ExecutionTarget | null {
 if (!view) return requested ? {...requested, ready: false} : null;
 const running = view.sessions.find((s) => s.state === 'running');
 if (running) return {taskID: running.task_id, stepID: running.step_id, ready: false};
 const selected = view.sessions.find((s) => s.id === view.selected_id);
 if ((!requested || requested.taskID === selected?.task_id) && selected?.step_id != null && view.steps.some((s) => s.id === selected.step_id && (s.done || s.deleted_at))) {
  const next = orderedSteps(view.steps, selected.task_id).find((s) => !s.done);
  if (next) return {taskID: selected.task_id, stepID: next.id, ready: true};
  return {taskID: selected.task_id, stepID: null, ready: true};
 }
 if (requested) {
  const completed = requested.stepID != null && view.steps.some((s) => s.id === requested.stepID && (s.done || s.deleted_at));
  if (completed) return {taskID: requested.taskID, stepID: orderedSteps(view.steps, requested.taskID).find((s) => !s.done)?.id ?? null, ready: true};
  return {...requested, ready: false};
 }
 return selected ? {taskID: selected.task_id, stepID: selected.step_id, ready: false} : null;
}


/** Returns no write for a canceled/same/foreign drop; otherwise moves one identity among siblings. */
export function reorderStepIDs(ids: number[], id: number, destination: number): number[] | null {
 const from = ids.indexOf(id); const to = ids.indexOf(destination);
 if (from < 0 || to < 0 || from === to) return null;
 const result = [...ids]; result.splice(from, 1); result.splice(to, 0, id); return result;
}
