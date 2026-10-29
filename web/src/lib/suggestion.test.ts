import { afterEach, describe, expect, it, vi } from 'vitest';
import { ApiError, type Pick, type Task } from './api';
import { SuggestionState, type SuggestionSnapshot } from './suggestion.svelte';
const task = (id: number): Task => ({id, title: `Task ${id}`, notes: '', raw_text: '', size: null, energy: null, feeling: '', due_on: null, not_before: null, status: 'open', refs: [], created_at: '', updated_at: ''});
const snapshot = (): SuggestionSnapshot => ({day: {date: '2026-10-29', day: null, plan: {id: 1, entries: [], opening: '', fallback: false}, unplanned: [1,2]}, tasks: [task(1), task(2)], blockers: [], events: [], revision: 1, execution: {server_now: '', steps: [], sessions: [], selected_id: null}});
const chosen: Pick = {task_id: 1, step_id: null, reason: 'A concrete start'};
const settle = async () => { for (let i = 0; i < 10; i++) await Promise.resolve(); };
afterEach(() => vi.restoreAllMocks());
describe('suggestion lifecycle', () => {
 it('asks once on entering Actions, preserves across views and ignores clock-only reads', async () => {
  const request = vi.fn().mockResolvedValue(chosen); const state = new SuggestionState(request, async () => {}); const input = snapshot();
  state.sync(input, true); await settle(); state.sync({...input, revision: 2, execution: {...input.execution, server_now: 'later'}}, false); state.sync(input, true); await settle();
  expect(request).toHaveBeenCalledTimes(1); expect(state.pick).toEqual(chosen); expect(state.stale).toBe(false);
 });
 it('re-evaluates captures after a confirmed read and after interaction ends', async () => {
  const request = vi.fn().mockResolvedValue(chosen); const state = new SuggestionState(request, async () => {}); const input = snapshot();
  state.sync(input, true); await settle(); state.changed(true); state.changed(true); state.sync(input, true); await settle(); expect(request).toHaveBeenCalledTimes(1);
  state.sync({...input, revision: 2, tasks: [...input.tasks, task(3)]}, true, true); await settle(); expect(request).toHaveBeenCalledTimes(1);
  state.sync({...input, revision: 2, tasks: [...input.tasks, task(3)]}, true); await settle(); expect(request).toHaveBeenCalledTimes(2);
 });
 it('ignores old responses and coalesces into one new request', async () => {
  const old = Promise.withResolvers<Pick>(); const request = vi.fn().mockReturnValueOnce(old.promise).mockResolvedValue({...chosen, task_id: 2}); const state = new SuggestionState(request, async () => {}); const input = snapshot();
  state.sync(input, true); state.changed(true); state.changed(true); state.sync({...input, revision: 2}, true); old.resolve(chosen); await settle();
  expect(request).toHaveBeenCalledTimes(2); expect(state.pick?.task_id).toBe(2);
 });
 it('removes a completed target and rejects results after another device starts focus', async () => {
  const pending = Promise.withResolvers<Pick>(); const request = vi.fn().mockResolvedValueOnce(chosen).mockReturnValueOnce(pending.promise); const state = new SuggestionState(request, async () => {}); const input = snapshot();
  state.sync(input, true); await settle(); state.sync({...input, tasks: [{...task(1), status: 'done'}, task(2)]}, true); expect(state.pick).toBeNull();
  state.sync({...input, execution: {...input.execution, sessions: [{id: 1, task_id: 2, step_id: null, state: 'running', segment_started_at: '', accumulated_ms: 0, created_at: '', ended_at: null}]}}, true);
  pending.resolve({...chosen, task_id: 2}); await settle(); expect(state.pick).toBeNull(); expect(request).toHaveBeenCalledTimes(2);
 });
 it('keeps valid choices after metadata edits with explicit update', async () => {
  const request = vi.fn().mockResolvedValue(chosen); const state = new SuggestionState(request, async () => {}); const input = snapshot(); state.sync(input, true); await settle();
  state.sync({...input, tasks: [{...task(1), estimate_minutes: 20}, task(2)]}, true); await settle(); expect(state.pick).toEqual(chosen); expect(state.stale).toBe(true); expect(request).toHaveBeenCalledTimes(1);
 });
 it('defers a result if editing begins while the request is in flight', async () => {
  const old = Promise.withResolvers<Pick>(); const request = vi.fn().mockReturnValueOnce(old.promise).mockResolvedValue(chosen); const state = new SuggestionState(request, async () => {}); const input = snapshot();
  state.sync(input, true); state.sync(input, true, true); old.resolve(chosen); await settle(); expect(state.pick).toBeNull(); expect(request).toHaveBeenCalledTimes(1);
  state.sync(input, true); await settle(); expect(state.pick).toEqual(chosen); expect(request).toHaveBeenCalledTimes(2);
 });
 it('passes exclusion and situational context', async () => {
  const request = vi.fn().mockResolvedValue(chosen); const state = new SuggestionState(request, async () => {}); state.sync(snapshot(), false); await state.ask(2, 'Tired right now'); expect(request).toHaveBeenCalledWith(2, 'Tired right now'); expect(state.pick).toEqual(chosen);
 });
 it('rejects missing steps without retry loops', async () => {
  const request = vi.fn().mockResolvedValue({...chosen, step_id: 99}); const state = new SuggestionState(request, async () => {}); state.sync(snapshot(), true); await settle(); expect(state.pick).toBeNull(); expect(state.problem).toContain('mudou'); state.sync(snapshot(), true); await settle(); expect(request).toHaveBeenCalledTimes(1);
 });
 it('preserves a refusal context when a request becomes outdated', async () => {
  const old = Promise.withResolvers<Pick>(); const request = vi.fn().mockReturnValueOnce(old.promise).mockResolvedValue(chosen); const state = new SuggestionState(request, async () => {}); const input = snapshot();
  state.sync(input, false); const asking = state.ask(2, 'tired now'); state.changed(); state.sync({...input, revision: 2}, true); old.resolve(chosen); await asking; await settle();
  expect(request).toHaveBeenLastCalledWith(2, 'tired now'); expect(request).toHaveBeenCalledTimes(2);
 });
 it('separates empty state and transport failure without automatic retries', async () => {
  const request = vi.fn().mockRejectedValueOnce(new ApiError(404, 'none')).mockRejectedValueOnce(new Error('offline')); const state = new SuggestionState(request, async () => {}); state.sync(snapshot(), true); await settle(); expect(state.nothing).toBe(true); expect(state.problem).toBe(''); await state.ask(); expect(state.problem).toContain('offline'); expect(state.nothing).toBe(false); state.sync(snapshot(), true); await settle(); expect(request).toHaveBeenCalledTimes(2);
 });
});
