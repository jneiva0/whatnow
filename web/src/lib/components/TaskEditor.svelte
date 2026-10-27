<script lang="ts">
	import ExecutionSteps from './ExecutionSteps.svelte';
	import { execution } from '#lib/execution.svelte.ts';
	import { taskEstimate } from '#lib/execution.ts';
	import { api, type Task } from '#lib/api.ts';
	import { editor } from '#lib/editor.svelte.ts';
	import { blankForm, fieldsFrom, formFor } from '#lib/taskForm.ts';
	import { toasts } from '#lib/toasts.svelte.ts';
	import { today } from '#lib/today.svelte.ts';

	interface Props {
		/** Null makes a new task. */
		task: Task | null;
	}
	let { task }: Props = $props();

	// Seeded once: the sheet is keyed by task, so a refetch mid-edit doesn't wipe what's typed.
	// svelte-ignore state_referenced_locally
	let form = $state(task ? formFor(task) : blankForm());
	// svelte-ignore state_referenced_locally
	let estimate = $state<number | undefined>(task?.estimate_minutes ?? undefined);
	const total = $derived(task ? taskEstimate({...task, estimate_minutes: estimate ?? null}, execution.view?.steps ?? []) : null);
	let problem = $state('');
	let saving = $state(false);
	let adding = $state('');
	/** A new task's blockers, added once it exists. */
	let pending = $state<number[]>([]);

	const titleOf = (id: number) => today.tasks.find((t) => t.id === id)?.title ?? `#${id}`;
	const blockedBy = $derived(
		task ? today.blockers.filter((b) => b.task_id === task.id).map((b) => b.blocked_by_id) : pending
	);
	const candidates = $derived(
		today.tasks.filter((t) => t.status === 'open' && t.id !== task?.id && !blockedBy.includes(t.id))
	);

	async function save(e: SubmitEvent) {
		e.preventDefault();
		const result = fieldsFrom(form, task);
		if ('error' in result) {
			problem = result.error;
			return;
		}
		saving = true;
		const fields = {...result.fields, estimate_minutes: estimate ?? null};
		const write = task
			? toasts.track(`Saved: ${fields.title}`, () => api.updateTask(task.id, fields))
			: toasts.track(`Added: ${fields.title}`, async () => {
					const created = await api.createTask(fields);
					for (const id of pending) await api.addBlocker(created.data.id, id);
					return created;
				});
		if (await write) editor.close();
		else problem = 'Não foi possível salvar a tarefa. Suas alterações continuam aqui.';
		saving = false;
	}

	function addBlocker() {
		const id = Number(adding);
		adding = '';
		if (!id) return;
		if (task) toasts.track(`Now waits on: ${titleOf(id)}`, () => api.addBlocker(task.id, id));
		else pending.push(id);
	}

	function removeBlocker(id: number) {
		if (task) toasts.track(`No longer waits on: ${titleOf(id)}`, () => api.removeBlocker(task.id, id));
		else pending = pending.filter((p) => p !== id);
	}
</script>

<header>
 <h2>{task ? 'Edit task' : 'New task'}</h2>
 <button type="button" class="close" aria-label="Close" onclick={() => editor.close()}><svg width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" aria-hidden="true"><path d="m6 6 12 12M18 6 6 18"/></svg></button>
</header>
<div class={['editor-body', {existing: !!task}]}>
<form id="task-fields" onsubmit={save}>
	<label>
		<span class="label">Title</span>
		<input class="field" bind:value={form.title} />
	</label>
	{#if problem}<p role="alert" class="problem">{problem}</p>{/if}
	<label>
		<span class="label">Notes</span>
		<textarea class="field" rows="3" bind:value={form.notes}></textarea>
	</label>
	<div class="row">
        <label>
            <span class="label">Estimativa (min)</span>
            <input class="field" type="number" min="1" max="1440" step="1" bind:value={estimate} placeholder="Não definida" />
        </label>
		<label>
			<span class="label">Energy</span>
			<select class="field" bind:value={form.energy}>
				<option value="">none</option>
				<option value="low">low</option>
				<option value="med">med</option>
				<option value="high">high</option>
			</select>
		</label>
	</div>

	<label>
		<span class="label">Feeling</span>
		<input class="field" bind:value={form.feeling} placeholder="How does this one feel?" />
	</label>
	<details class="advanced" open={!!form.due_on || !!form.not_before || blockedBy.length > 0}>
	<summary>Datas e dependências</summary>
	<div class="row">
		<label>
			<span class="label">Due</span>
			<input class="field" type="date" bind:value={form.due_on} />
		</label>
		<label>
			<span class="label">Not before</span>
			<input class="field" type="date" bind:value={form.not_before} />
		</label>
	</div>

	<section>
		<span class="label">Blocked by</span>
		<ul>
			{#each blockedBy as id (id)}
				<li>
					{titleOf(id)}
					<button type="button" class="close" aria-label="Remove blocker {titleOf(id)}" onclick={() => removeBlocker(id)}>✕</button>
				</li>
			{:else}
				<li class="muted">Nothing. It's free to start.</li>
			{/each}
		</ul>
		{#if candidates.length}
			<select class="field" bind:value={adding} onchange={addBlocker} aria-label="Add a blocker">
				<option value="">Add a blocker…</option>
				{#each candidates as t (t.id)}
					<option value={String(t.id)}>{t.title}</option>
				{/each}
			</select>
		{/if}
	</section>

	</details>
	{#if task?.raw_text}
		<details class="raw-details">
			<summary>Captura original</summary>
			<p class="raw">{task.raw_text}</p>
		</details>
	{/if}

</form>
{#if task}<section class="steps-panel" aria-label="Passos da tarefa"><div class="steps-title"><h3>Passos</h3>{#if total?.minutes != null}<span>{total.minutes}{total.partial ? '+' : ''} min</span>{/if}</div><p class="save-hint">Alterações nos passos são salvas na hora.</p><ExecutionSteps taskId={task.id} /></section>{/if}
</div>
<footer class="actions">
 <button form="task-fields" class="btn primary" disabled={saving}>{saving ? 'Saving…' : 'Save'}</button>
 <button type="button" class="btn" onclick={() => editor.close()}>Cancel</button>
</footer>

<style>
 .editor-body { display: grid; gap: 24px; padding-block: 20px; }
 .steps-panel { min-width: 0; border-top: 1px solid var(--paper-deep); padding-top: 20px; }
 .steps-title { display: flex; align-items: baseline; justify-content: space-between; gap: 12px; }
 h3 { margin: 0; font-size: 18px; } .steps-title span { font-variant-numeric: tabular-nums; }
 .save-hint { font-size: 12px; color: var(--ink-soft); margin: 4px 0 12px; }
 summary { cursor: pointer; padding-block: 12px; min-height: 44px; }
 .advanced, .raw-details { border-top: 1px solid var(--paper-deep); }
 header { position: sticky; top: 0; z-index: 2; background: var(--paper); padding-block: 12px; border-bottom: 1px solid var(--paper-deep); }
 .actions { position: sticky; bottom: 0; background: var(--paper); padding-block: 16px; border-top: 1px solid var(--paper-deep); z-index: 2; }
 @media(min-width: 900px) { .editor-body.existing { grid-template-columns: minmax(0, 1fr) minmax(0, 1.25fr); gap: 32px; } .steps-panel { border-top: 0; border-left: 1px solid var(--paper-deep); padding: 0 0 0 24px; } }
 @media(max-width: 480px) { .row { flex-wrap: wrap; } .row label { min-width: 120px; } }
	form {
		display: flex;
		flex-direction: column;
		gap: 12px;
	}
	header {
		display: flex;
		justify-content: space-between;
		align-items: center;
	}
	h2 {
		margin: 0;
		font-weight: var(--heavy);
	}
	label {
		display: flex;
		flex-direction: column;
		gap: 4px;
		flex: 1;
		min-width: 0;
	}
	.row {
		display: flex;
		gap: 12px;
	}
	ul {
		margin: 6px 0;
		padding-left: 18px;
	}
	.close {
        min-width: 44px; min-height: 44px; display: inline-grid; place-items: center;
		background: none;
		border: none;
	}
	.raw {
		margin: 4px 0 0;
		white-space: pre-wrap;
		font-size: 14px;
		color: var(--ink-soft);
	}
	.problem {
		margin: 0;
		color: var(--danger);
		font-weight: var(--heavy);
	}
	.actions {
		display: flex;
		gap: 10px;
	}
</style>
