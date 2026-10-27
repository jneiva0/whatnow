<script lang="ts">
	import { goto } from '$app/navigation';
	import { execution } from '#lib/execution.svelte.ts';
	import { type Pick, type Task } from '#lib/api.ts';

	interface Props {
		pick: Pick;
		task: Task | undefined;
		asking: boolean;
		/** Asks again without this task, with the user's optional reason. */
		onnotthis: (reason: string) => void;
		onclose: () => void;
	}
	let { pick, task, asking, onnotthis, onclose }: Props = $props();

	let explaining = $state(false);
	let reason = $state('');
	let starting = $state(false);
	let problem = $state('');

	const step = $derived(execution.view?.steps.find((s) => s.id === pick.step_id));
	async function begin() {
		if (starting || execution.busy || asking) return;
		starting = true;
		problem = '';
		const {task_id, step_id} = pick;
		try {
			if (!await execution.start(task_id, step_id ?? null)) {
				problem = execution.problem;
				return;
			}
			await goto(`/focus?task=${task_id}${step_id ? `&step=${step_id}` : '&direct=1'}`);
			onclose();
		} catch {
			problem = 'O foco iniciou, mas a tela não abriu. Abra Foco para continuar.';
		} finally {
			starting = false;
		}
	}

	function notThis(e: SubmitEvent) {
		e.preventDefault();
		onnotthis(reason.trim());
		explaining = false;
		reason = '';
	}
</script>

<section class="sticker whatnow" aria-live="polite" data-task={pick.task_id}>
	<button class="close" disabled={starting} aria-label="Close" onclick={onclose}>✕</button>
	<h2>{step?.description ?? task?.title ?? 'That task'}</h2>
	{#if step}<p class="parent">{task?.title}</p>{/if}
	<p>{pick.reason}</p>
	{#if problem}<p role="alert">{problem}</p>{/if}
	{#if explaining}
		<form onsubmit={notThis}>
			<input class="field" bind:value={reason} placeholder="Why not? (optional)" aria-label="Why not this one" />
			<button class="btn" disabled={asking}>Ask again</button>
		</form>
	{:else}
		<div class="actions">
			<button class="btn primary" disabled={starting || execution.busy || asking} onclick={begin}>{starting ? 'Iniciando…' : 'Começar'}</button>
			<button class="btn" disabled={asking || starting} onclick={() => (explaining = true)}>
				{asking ? 'Thinking…' : 'Not this one'}
			</button>
		</div>
	{/if}
</section>

<style>
	.whatnow {
		position: relative;
		background: var(--chip);
		padding: 16px 20px;
	}
	h2 {
		margin: 4px 0;
		font-size: 22px;
		font-weight: var(--heavy);
	}
	.close {
		position: absolute;
		top: 10px;
		right: 12px;
		background: none;
		border: none;
	}
	.actions,
	form {
		display: flex;
		gap: 10px;
	}
</style>
