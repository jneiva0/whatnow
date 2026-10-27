<script lang="ts">
	import { goto } from '$app/navigation';
	import type { ActionRef } from '#lib/api.ts';
	import { chat } from '#lib/chat.svelte.ts';
	import { execution } from '#lib/execution.svelte.ts';
	import { available } from '#lib/execution.ts';
	import { today } from '#lib/today.svelte.ts';
	import { toasts } from '#lib/toasts.svelte.ts';
	import { Button } from '#lib/components/ui/button/index.js';
	import PlayIcon from '@lucide/svelte/icons/play';

	let { action }: { action: ActionRef } = $props();
	let starting = $state(false);
	const task = $derived(today.tasks.find((t) => t.id === action.task_id));
	const step = $derived(execution.view?.steps.find((s) => s.task_id === action.task_id && s.id === action.step_id && !s.deleted_at && !s.done));
	const ready = $derived(!!task && available(task, today.tasks, today.blockers, today.today?.date ?? '') && (action.step_id === null || !!step));
	const label = $derived(step?.description ?? task?.title ?? 'ação');

	async function begin() {
		if (!ready || starting || execution.busy || chat.busy) return;
		starting = true;
		try {
			if (!await execution.start(action.task_id, action.step_id)) {
				toasts.error(execution.problem);
				return;
			}
			await goto(`/focus?task=${action.task_id}${action.step_id === null ? '&direct=1' : `&step=${action.step_id}`}`);
			chat.open = false;
		} catch {
			toasts.error('O foco iniciou, mas a tela não abriu. Abra Foco para continuar.');
		} finally {
			starting = false;
		}
	}
</script>

{#if ready}
	<Button variant="outline" disabled={starting || execution.busy || chat.busy} onclick={begin} aria-label={`Começar: ${label}`} title={label}>
		<PlayIcon data-icon="inline-start" aria-hidden="true" />
		{starting ? 'Iniciando…' : 'Começar esta ação'}
	</Button>
{/if}
