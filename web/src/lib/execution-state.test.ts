import { afterEach, describe, expect, it, vi } from 'vitest';
import { api, type ExecutionState } from './api';
import { Execution } from './execution.svelte';
import { today } from './today.svelte';

const view = (): ExecutionState => ({server_now: new Date(Date.now() + 120000).toISOString(), steps: [], sessions: [], selected_id: null});
afterEach(() => {vi.restoreAllMocks(); vi.unstubAllGlobals(); vi.useRealTimers();});

describe('execution authority and recovery', () => {
 it('compensates server clock skew and retains the last confirmed state after a failure', async () => {
  const snapshot = {...view(), estimates: {2: {minutes: 45, partial: true}}};
  vi.spyOn(api, 'execution').mockResolvedValue(snapshot);
  const state = new Execution();
  await state.refresh();
  expect(Math.abs(state.now - Date.now() - 120000)).toBeLessThan(100);
  // Render the server total even while the separate step list is empty/stale.
  expect(state.estimate(2)).toEqual({minutes: 45, partial: true});
  const confirmed = state.view;
  vi.spyOn(api, 'focusAction').mockRejectedValue(new Error('offline'));
  expect(await state.action(1, 'pause')).toBe(false);
  expect(state.view).toBe(confirmed);
  expect(state.estimate(2)).toEqual({minutes: 45, partial: true});
  expect(state.problem).toContain('não foi confirmada');
 });
 it('reuses the same request token after an uncertain start and clears it after confirmation', async () => {
  const start = vi.spyOn(api, 'startFocus').mockRejectedValueOnce(new Error('response lost')).mockResolvedValue(view());
  vi.spyOn(today, 'changed').mockImplementation(() => {});
  const state = new Execution();
  await state.start(2, null);
  await state.start(2, null);
  expect(start.mock.calls[0][2]).toBe(start.mock.calls[1][2]);
  await state.start(2, null);
  expect(start.mock.calls[2][2]).not.toBe(start.mock.calls[1][2]);
 });
 it('catches up once when stream updates arrive while a local focus command is saving', async () => {
  const pending = Promise.withResolvers<ExecutionState>();
  const paused = {...view(), selected_id: 1};
  const newer = {...view(), selected_id: 2};
  vi.spyOn(api, 'focusAction').mockReturnValue(pending.promise);
  const read = vi.spyOn(api, 'execution').mockResolvedValue(newer);
  vi.spyOn(today, 'changed').mockImplementation(() => {});
  const state = new Execution();
  const saving = state.action(1, 'pause');
  await state.refresh();
  await state.refresh();
  expect(read).not.toHaveBeenCalled();
  pending.resolve(paused);
  expect(await saving).toBe(true);
  await Promise.resolve();
  expect(read).toHaveBeenCalledTimes(1);
  expect(state.view).toBe(newer);
 });
 it('keeps an unconfirmed command visible even when the catch-up read succeeds', async () => {
  const pending = Promise.withResolvers<ExecutionState>();
  vi.spyOn(api, 'focusAction').mockReturnValue(pending.promise);
  const latest = view();
  vi.spyOn(api, 'execution').mockResolvedValue(latest);
  const state = new Execution();
  const saving = state.action(1, 'pause');
  await state.refresh();
  pending.reject(new Error('response lost'));
  expect(await saving).toBe(false);
  await Promise.resolve();
  expect(state.view).toBe(latest);
  expect(state.problem).toContain('não foi confirmada');
 });
 it('announces the watched action estimate after a suspended tab catches up without a timer write', async () => {
  vi.useFakeTimers();
  vi.stubGlobal('document', Object.assign(new EventTarget(), {visibilityState: 'visible'}));
  vi.stubGlobal('window', new EventTarget());
  const start = new Date(Date.now() - 59000).toISOString();
  const snapshot: ExecutionState = {server_now: new Date().toISOString(), selected_id: 1, steps: [],
   estimates: {2: {minutes: 1, partial: false}}, sessions: [{id: 1, task_id: 2, step_id: null, state: 'running', segment_started_at: start, accumulated_ms: 0, created_at: start, ended_at: null}]};
  vi.spyOn(api, 'execution').mockResolvedValue(snapshot);
  const write = vi.spyOn(api, 'focusAction');
  const notify = vi.fn();
  const state = new Execution();
  const stop = state.startLive(notify);
  await vi.advanceTimersByTimeAsync(0);
  vi.setSystemTime(Date.now() + 120000);
  await vi.advanceTimersByTimeAsync(1000);
  expect(notify).toHaveBeenCalledOnce();
  expect(notify).toHaveBeenCalledWith(expect.objectContaining({task_id: 2, step_id: null, minutes: 1}));
  await vi.advanceTimersByTimeAsync(3000);
  expect(notify).toHaveBeenCalledOnce();
  expect(write).not.toHaveBeenCalled();
  expect(state.view?.sessions[0].state).toBe('running');
  stop();
 });
 it('polls every five seconds only while visible and SSE is disconnected', async () => {
  vi.useFakeTimers();
  const doc = Object.assign(new EventTarget(), {visibilityState: 'visible'});
  vi.stubGlobal('document', doc);
  vi.stubGlobal('window', new EventTarget());
  const read = vi.spyOn(api, 'execution').mockImplementation(async () => view());
  vi.spyOn(today, 'changed').mockImplementation(() => {});
  const state = new Execution();
  const stop = state.startLive();
  await vi.advanceTimersByTimeAsync(10000);
  expect(read).toHaveBeenCalledTimes(1);
  state.connected(false);
  await vi.advanceTimersByTimeAsync(5000);
  expect(read).toHaveBeenCalledTimes(2);
  doc.visibilityState = 'hidden'; doc.dispatchEvent(new Event('visibilitychange'));
  await vi.advanceTimersByTimeAsync(10000);
  expect(read).toHaveBeenCalledTimes(2);
  doc.visibilityState = 'visible'; doc.dispatchEvent(new Event('visibilitychange'));
  await vi.advanceTimersByTimeAsync(5000);
  expect(read).toHaveBeenCalledTimes(4);
  state.connected(true);
  await vi.advanceTimersByTimeAsync(10000);
  expect(read).toHaveBeenCalledTimes(4);
  stop();
 });
});
