import { goto } from '$app/navigation';
import type { ActionRef } from './api';
import { execution } from './execution.svelte';

/** Web-only orchestration. The service owns switching, time and eligibility.
 * null means navigation succeeded; callers own how an error is presented. */
export async function startAndOpenFocus(action: ActionRef, fromActions = false): Promise<string | null> {
	if (execution.busy) return 'Aguarde a ação em andamento e tente novamente.';
	if (!await execution.start(action.task_id, action.step_id)) return execution.problem;
	const query = new URLSearchParams({ task: String(action.task_id) });
	if (action.step_id === null) query.set('direct', '1');
	else query.set('step', String(action.step_id));
	if (fromActions) query.set('from', 'actions');
	try {
		await goto(`/focus?${query}`);
		return null;
	} catch {
		return 'O foco iniciou, mas a tela não abriu. Abra Foco para continuar.';
	}
}
