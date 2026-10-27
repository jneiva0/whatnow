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
