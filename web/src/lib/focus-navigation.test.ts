import { afterEach, describe, expect, it, vi } from 'vitest';
import { goto } from '$app/navigation';
import { api } from './api';
import { execution } from './execution.svelte';
import { today } from './today.svelte';
import { startAndOpenFocus } from './focus-navigation';

vi.mock('$app/navigation', () => ({ goto: vi.fn() }));
afterEach(() => vi.restoreAllMocks());

describe('web focus entry', () => {
	it('does not navigate or replace confirmed state when starting fails', async () => {
		const confirmed = execution.view;
		vi.spyOn(api, 'startFocus').mockRejectedValue(new Error('offline'));
		expect(await startAndOpenFocus({task_id: 2, step_id: 9})).toContain('não foi confirmada');
		expect(goto).not.toHaveBeenCalled();
		expect(execution.view).toBe(confirmed);
	});
	it('waits for server confirmation before navigating and preserves direct execution', async () => {
		const pending = Promise.withResolvers<Awaited<ReturnType<typeof api.startFocus>>>();
		vi.spyOn(api, 'startFocus').mockReturnValue(pending.promise);
		vi.spyOn(today, 'changed').mockImplementation(() => {});
		const opening = startAndOpenFocus({task_id: 2, step_id: null}, true);
		expect(goto).not.toHaveBeenCalled();
		pending.resolve({server_now: new Date().toISOString(), steps: [], sessions: [], selected_id: null});
		expect(await opening).toBeNull();
		expect(goto).toHaveBeenCalledWith('/focus?task=2&direct=1&from=actions');
	});
	it('reports a navigation failure without retrying or undoing the confirmed start', async () => {
		vi.mocked(goto).mockRejectedValueOnce(new Error('navigation failed'));
		const start = vi.spyOn(api, 'startFocus').mockResolvedValue({server_now: new Date().toISOString(), steps: [], sessions: [], selected_id: 7});
		vi.spyOn(today, 'changed').mockImplementation(() => {});
		expect(await startAndOpenFocus({task_id: 2, step_id: 9})).toContain('O foco iniciou');
		expect(goto).toHaveBeenCalledWith('/focus?task=2&step=9');
		expect(start).toHaveBeenCalledTimes(1);
		expect(execution.view?.selected_id).toBe(7);
	});
});
