<script lang="ts">
	import { goto } from '$app/navigation';
	import { execution } from '#lib/execution.svelte.ts';
	import TaskMetadata from './TaskMetadata.svelte';
	import PlayToggle from './PlayToggle.svelte';
	import ExecutionSteps from './ExecutionSteps.svelte';
	import { available, actionActual, taskActual, timerText } from '#lib/execution.ts';
	import { api, type Tier } from '#lib/api.ts';
	import * as DropdownMenu from '#lib/components/ui/dropdown-menu/index.js';
	import * as ContextMenu from '#lib/components/ui/context-menu/index.js';
	import TaskActionItems from './TaskActionItems.svelte';
	import { buttonVariants } from '#lib/components/ui/button/index.js';
	import EllipsisIcon from '@lucide/svelte/icons/ellipsis';
	import { editor } from '#lib/editor.svelte.ts';
	import { linked } from '#lib/linked.svelte.ts';
	import { linkify } from '#lib/linkify.ts';
	import { links, type Card } from '#lib/plan.ts';
	import { toasts } from '#lib/toasts.svelte.ts';
	import { today } from '#lib/today.svelte.ts';

	interface Props {
		card: Card;
		/** Absent for an unplanned card. */
		tier?: Tier;
		/** Today's date, for push to tomorrow. */
		date: string;
		/** All tasks shows dates, and opens the editor on a click instead of expanding the card. */
		listed?: boolean;
		/** Present one executable action using the shared step controls. */
		actionView?: boolean;
	}
	let { card, tier, date, listed = false, actionView = false }: Props = $props();

	let expanded = $state(false);
	const task = $derived(card.task);
	const steps = $derived(execution.steps(task.id));
	const running = $derived(execution.view?.sessions.find((s) => s.task_id === task.id && s.state === 'running'));
 const nextStep = $derived(running ? steps.find((s) => s.id === running.step_id) : steps.find((s) => !s.done));
 const nextActual = $derived(nextStep ? actionActual(execution.view?.sessions ?? [], task.id, nextStep.id, execution.now) : 0);
	const completed = $derived(steps.filter((s) => s.done).length);
	const actual = $derived(taskActual(execution.view?.sessions ?? [], task.id, execution.now));
	const estimate = $derived(execution.estimate(task.id));
	const canStart = $derived(available(task, today.tasks, today.blockers, date));
	const chain = $derived(links(task.id, today.blockers, today.tasks));
	const showingLinks = $derived(linked.source === task.id);
	const titles = (tasks: { title: string }[]) => tasks.map((t) => t.title).join(', ');
	const toggleLinks = () => linked.toggle(task.id, [...chain.waitingOn, ...chain.unblocks].map((t) => t.id));

	// The card leaves once the refetch lands. If the write failed, untick it again.
	async function complete(e: Event & { currentTarget: HTMLInputElement }) {
		const box = e.currentTarget;
		if (!(await toasts.track(`Done: ${task.title}`, () => api.setStatus(task.id, 'done')))) {
			box.checked = false;
		}
	}
	function restore() {
		toasts.track(`Restored: ${task.title}`, () => api.setStatus(task.id, 'open'));
	}
	// Keep the browser's editing menu available inside inline text/number fields.
	function editingField(target: EventTarget | null) {
		return target instanceof Element && !!target.closest('input, textarea, [contenteditable="true"]');
	}
</script>

<ContextMenu.Root>
<ContextMenu.Trigger disabled={task.status !== 'open'}>
{#snippet child({ props })}
<article
	{...props}
	oncontextmenu={(event) => { if (!editingField(event.target) && typeof props.oncontextmenu === 'function') props.oncontextmenu(event); }}
	onpointerdown={(event) => {
		// A long press on an editor, link or drag handle keeps its own interaction.
		if (!(event.target instanceof Element && event.target.closest('button, a, input, textarea, summary, [contenteditable="true"]'))) { if (typeof props.onpointerdown === 'function') props.onpointerdown(event); }
	}}
	data-task={task.id}
	class={['sticker card', tier ?? (listed ? 'listed' : 'unplanned'), {closed: task.status !== 'open', linked: linked.ids.includes(task.id)}]}
>
	{#if task.status === 'open'}
		<input type="checkbox" aria-label={actionView ? `Concluir tarefa: ${task.title}` : 'Mark done'} onchange={complete} />
	{:else}
		<span class="mark" aria-label={task.status}>{task.status === 'done' ? '✓' : '✕'}</span>
	{/if}
	<div class="main">
		<button
			class="body"
			onclick={() => (listed ? editor.open(task.id) : (expanded = !expanded))}
			disabled={!listed && !card.rationale && !task.notes}
			aria-expanded={listed ? undefined : expanded}
		>
			<span class="title">{task.title}</span>
		</button>
		<TaskMetadata {task} />
			{#if task.notes && !expanded}<span class="notes">{task.notes}</span>{/if}
			{#if card.rationale}<span class={['why', {shown: expanded}]}>{card.rationale}</span>{/if}
		<!-- Expanded notes sit outside the button, so their links work and their text can be copied. -->
		{#if task.notes && expanded}
			<p class="notes shown">{#each linkify(task.notes) as s, i (i)}{#if s.href}<a href={s.href} target="_blank" rel="noopener noreferrer">{s.text}</a>{:else}{s.text}{/if}{/each}</p>
		{/if}
		{#if task.status === 'open' && (chain.waitingOn.length || chain.unblocks.length)}
			<div class="chain">
				{#if chain.waitingOn.length}
					<button
						class="badge link"
						aria-expanded={showingLinks}
						onclick={toggleLinks}
					>
						⛓ waits on {chain.waitingOn.length}
					</button>
				{/if}
				{#if chain.unblocks.length}
					<button
						class="badge link"
						aria-expanded={showingLinks}
						onclick={toggleLinks}
					>
						unblocks {chain.unblocks.length}
					</button>
				{/if}
			</div>
			{#if showingLinks}
				<p class="links">
					{#if chain.waitingOn.length}<span>Waiting on: {titles(chain.waitingOn)}</span>{/if}
					{#if chain.unblocks.length}<span>Unblocks: {titles(chain.unblocks)}</span>{/if}
				</p>
			{/if}
		{/if}

        {#if task.status === 'open'}
          {#if actionView && steps.length && !(running && running.step_id === null)}
            {#if !nextStep}<p class="steps-ready">Passos concluídos. Marque a tarefa quando estiver pronta.</p>{/if}
            <ExecutionSteps taskId={task.id} eligible={canStart} preview />
          {:else}
            <div class="execution">
                <div class="step-progress">{#if nextStep}<span class="next-action">{nextStep.description}</span><span class="next-time">{timerText(nextActual)} / {nextStep.estimate_minutes == null ? '—' : timerText(nextStep.estimate_minutes * 60000)}{running ? ' · Em foco' : ''}</span>{/if}<span class={['total-time', {secondary: !!nextStep}]}>{nextStep ? 'Total ' : ''}{timerText(actual)} / {estimate.minutes == null ? '—' : timerText(estimate.minutes * 60000)}{estimate.partial ? '+' : ''}</span></div>
                {#if !steps.length || nextStep || running}<PlayToggle running={!!running} disabled={(!canStart && !running) || execution.busy} label={running ? `Pausar ${task.title}` : `Começar ${task.title}`} onclick={async () => { if (running) await execution.action(running.id, 'pause'); else if (await execution.start(task.id, nextStep?.id ?? null)) goto(`/focus?task=${task.id}${nextStep ? `&step=${nextStep.id}` : '&direct=1'}${actionView ? '&from=actions' : ''}`); }} />{/if}
                <a class="focus-icon" href={`/focus?task=${task.id}${nextStep ? `&step=${nextStep.id}` : ''}`} aria-label={`Abrir foco de ${task.title}`} title="Abrir foco"><svg width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" aria-hidden="true"><path d="M8 3H3v5M16 3h5v5M21 16v5h-5M8 21H3v-5"/><circle cx="12" cy="12" r="3"/></svg></a>
            </div>
            <details class="card-steps"><summary>{steps.length ? `${completed}/${steps.length} passos` : 'Adicionar passos'}{steps.length && !nextStep ? ' · Todos concluídos' : ''}</summary><ExecutionSteps taskId={task.id} eligible={canStart} /></details>
          {/if}
        {/if}
	</div>
	{#if task.status === 'dropped'}
		<button class="btn" onclick={restore}>Restore</button>
	{:else if task.status === 'open'}
        <DropdownMenu.Root>
            <DropdownMenu.Trigger class={buttonVariants({variant: 'ghost', size: 'icon'})} aria-label={`Mais ações para ${task.title}`}>
                <EllipsisIcon aria-hidden="true" />
            </DropdownMenu.Trigger>
            <DropdownMenu.Content align="end">
                <TaskActionItems {task} {tier} {date} {listed} />
            </DropdownMenu.Content>
        </DropdownMenu.Root>
	{/if}
</article>
{/snippet}
</ContextMenu.Trigger>
<ContextMenu.Content>
	<TaskActionItems {task} {tier} {date} {listed} context />
</ContextMenu.Content>
</ContextMenu.Root>

<style>
 .steps-ready { margin: 8px 0 0; font-size: 13px; }
 .focus-icon { display: inline-grid; place-items: center; min-width: 44px; min-height: 44px; color: var(--ink); }
	.execution { display: flex; gap: 8px; align-items: center; flex-wrap: wrap; margin-top: 6px; }
	.step-progress { flex: 1; min-width: 0; display: flex; flex-direction: column; gap: 4px; font-variant-numeric: tabular-nums; }
 .next-action { font-size: 16px; font-weight: 650; overflow-wrap: anywhere; } .next-time, .total-time { font-size: 18px; } .total-time.secondary { font-size: 12px; color: var(--ink-soft); }
 .execution { border-top: 1px solid color-mix(in srgb, var(--ink) 20%, transparent); padding-top: 12px; }
	.card-steps { font-size: 13px; }
	.card-steps summary { cursor: pointer; min-height: 36px; padding-block: 7px; overflow-wrap: anywhere; }
	.card {
		display: flex;
		align-items: flex-start;
		gap: 10px;
		padding: 10px 12px;
		margin-bottom: 10px;
	}
	.must {
		background: var(--must);
	}
	.should {
		background: var(--should);
	}
	.could {
		background: var(--could);
	}
	.unplanned {
		border-style: dashed;
		box-shadow: none;
	}
	.closed {
		opacity: 0.7;
	}
	.mark {
		width: 20px;
		text-align: center;
		font-weight: var(--heavy);
	}
	input[type='checkbox'] {
		width: 20px;
		height: 20px;
		margin-top: 2px;
		accent-color: var(--primary);
	}
	.main {
		flex: 1;
		display: flex;
		flex-direction: column;
		gap: 4px;
		min-width: 0;
	}
	.body {
		display: flex;
		flex-direction: column;
		gap: 3px;
		background: none;
		border: none;
		padding: 0;
		text-align: left;
		cursor: default;
	}
	.body:not(:disabled) {
		cursor: pointer;
	}
	.title {
		font-weight: var(--heavy);
	}

	.badge {
		font-size: 11px;
		font-weight: var(--heavy);
		border: 1.5px solid var(--ink);
		border-radius: var(--pill);
		padding: 0 7px;
		background: var(--card);
	}
	.chain {
		display: flex;
		gap: 6px;
	}
	.link {
		cursor: pointer;
		color: var(--ink);
	}
	.link[aria-expanded='true'] {
		background: var(--chip);
	}
	.links {
		margin: 0;
		display: flex;
		flex-direction: column;
		font-size: 13px;
		color: var(--ink-soft);
	}
	.linked {
		outline: 2px dashed var(--primary);
		outline-offset: 3px;
	}

	.notes {
		font-size: 13px;
		white-space: pre-line;
		overflow: hidden;
		text-overflow: ellipsis;
		display: -webkit-box;
		-webkit-line-clamp: 1;
		line-clamp: 1;
		-webkit-box-orient: vertical;
	}
	.notes.shown {
		display: block;
		margin: 0;
		overflow-wrap: anywhere;
	}
	.why {
		display: none;
		font-size: 13px;
	}
	.why.shown,
	.card:hover .why {
		display: block;
	}
</style>
