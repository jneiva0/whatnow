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

	const step = $derived(execution.view?.steps.find((s) => s.id === pick.step_id));
	function begin() {
		const destination = `/focus?task=${pick.task_id}${pick.step_id ? `&step=${pick.step_id}` : ''}`;
		onclose();
		goto(destination);
	}

	function notThis(e: SubmitEvent) {
		e.preventDefault();
		onnotthis(reason.trim());
		explaining = false;
		reason = '';
	}
</script>

<section class="sticker whatnow" aria-live="polite" data-task={pick.task_id}>
	<button class="close" aria-label="Close" onclick={onclose}>✕</button>
	<h2>{step?.description ?? task?.title ?? 'That task'}</h2>
	{#if step}<p class="parent">{task?.title}</p>{/if}
	<p>{pick.reason}</p>
	{#if explaining}
		<form onsubmit={notThis}>
			<input class="field" bind:value={reason} placeholder="Why not? (optional)" aria-label="Why not this one" />
			<button class="btn" disabled={asking}>Ask again</button>
		</form>
	{:else}
		<div class="actions">
			<button class="btn primary" onclick={begin}>Começar</button>
			<button class="btn" disabled={asking} onclick={() => (explaining = true)}>
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
