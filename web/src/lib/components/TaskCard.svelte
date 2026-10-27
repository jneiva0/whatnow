<script lang="ts">
	import { goto } from '$app/navigation';
	import { execution } from '#lib/execution.svelte.ts';
	import TaskMetadata from './TaskMetadata.svelte';
	import PlayToggle from './PlayToggle.svelte';
	import ExecutionSteps from './ExecutionSteps.svelte';
	import { available, actionActual, taskActual, taskEstimate, timerText } from '#lib/execution.ts';
	import { api, type Tier } from '#lib/api.ts';
	import { dismiss } from '#lib/dismiss.ts';
	import { editor } from '#lib/editor.svelte.ts';
	import { linked } from '#lib/linked.svelte.ts';
	import { linkify } from '#lib/linkify.ts';
	import { links, tierLabels, type Card } from '#lib/plan.ts';
	import { addDays } from '#lib/time.ts';
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
	}
	let { card, tier, date, listed = false }: Props = $props();

	let menuOpen = $state(false);
	let expanded = $state(false);
	const task = $derived(card.task);
	const steps = $derived(execution.steps(task.id));
	const running = $derived(execution.view?.sessions.find((s) => s.task_id === task.id && s.state === 'running'));
 const nextStep = $derived(steps.find((s) => s.id === running?.step_id) ?? steps.find((s) => !s.done));
 const nextActual = $derived(nextStep ? actionActual(execution.view?.sessions ?? [], task.id, nextStep.id, execution.now) : 0);
	const completed = $derived(steps.filter((s) => s.done).length);
	const actual = $derived(taskActual(execution.view?.sessions ?? [], task.id, execution.now));
	const estimate = $derived(taskEstimate(task, execution.view?.steps ?? []));
	const canStart = $derived(available(task, today.tasks, today.blockers, date));
	const chain = $derived(links(task.id, today.blockers, today.tasks));
	const movable = $derived(!listed && !!today.today?.plan && !today.today.day?.wrapped_at);
	const moves: { to: Tier | null; label: string }[] = [
		{ to: 'must', label: tierLabels.must },
		{ to: 'should', label: tierLabels.should },
		{ to: 'could', label: tierLabels.could },
		{ to: null, label: 'Unplanned' }
	];
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
	function push() {
		menuOpen = false;
		toasts.track(`Pushed to tomorrow: ${task.title}`, () => api.push(task.id, addDays(date, 1)));
	}
	function drop() {
		menuOpen = false;
		toasts.track(`Dropped: ${task.title}`, () => api.setStatus(task.id, 'dropped'));
	}
	function move(to: Tier | null, label: string) {
		menuOpen = false;
		toasts.track(`Moved to ${label}: ${task.title}`, () => api.move(task.id, to));
	}
	function edit() {
		menuOpen = false;
		editor.open(task.id);
	}
	function restore() {
		toasts.track(`Restored: ${task.title}`, () => api.setStatus(task.id, 'open'));
	}
	const short = (date: string) =>
		new Date(`${date}T12:00:00`).toLocaleDateString(undefined, { month: 'short', day: 'numeric' });
</script>

<article
	data-task={task.id}
	class={['sticker card', tier ?? (listed ? 'listed' : 'unplanned'), {closed: task.status !== 'open', linked: linked.ids.includes(task.id)}]}
>
	{#if task.status === 'open'}
		<input type="checkbox" aria-label="Mark done" onchange={complete} />
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
            <div class="execution">
                <div class="step-progress">{#if nextStep}<span class="next-action">{nextStep.description}</span><span class="next-time">{timerText(nextActual)} / {nextStep.estimate_minutes == null ? '—' : timerText(nextStep.estimate_minutes * 60000)}{running ? ' · Em foco' : ''}</span>{/if}<span class={['total-time', {secondary: !!nextStep}]}>{nextStep ? 'Total ' : ''}{timerText(actual)} / {estimate.minutes == null ? '—' : timerText(estimate.minutes * 60000)}{estimate.partial ? '+' : ''}</span></div>
                {#if !steps.length || nextStep}<PlayToggle running={execution.selected?.task_id === task.id && execution.selected.state === 'running'} disabled={!canStart || execution.busy} label={execution.selected?.task_id === task.id && execution.selected.state === 'running' ? `Pausar ${task.title}` : `Começar ${task.title}`} onclick={async () => { if (execution.selected?.task_id === task.id && execution.selected.state === 'running') await execution.action(execution.selected.id, 'pause'); else if (await execution.start(task.id, nextStep?.id ?? null)) goto(`/focus?task=${task.id}${nextStep ? `&step=${nextStep.id}` : ''}`); }} />{/if}
                <a class="focus-icon" href={`/focus?task=${task.id}${nextStep ? `&step=${nextStep.id}` : ''}`} aria-label={`Abrir foco de ${task.title}`} title="Abrir foco"><svg width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" aria-hidden="true"><path d="M8 3H3v5M16 3h5v5M21 16v5h-5M8 21H3v-5"/><circle cx="12" cy="12" r="3"/></svg></a>
            </div>
            <details class="card-steps"><summary>{steps.length ? `${completed}/${steps.length} passos` : 'Adicionar passos'}{steps.length && !nextStep ? ' · Todos concluídos' : ''}</summary><ExecutionSteps taskId={task.id} eligible={canStart} /></details>
        {/if}
	</div>
	{#if task.status === 'dropped'}
		<button class="btn" onclick={restore}>Restore</button>
	{:else if task.status === 'open'}
		<details class="more" bind:open={menuOpen} {@attach dismiss}>
			<summary aria-label="More actions">⋯</summary>
			<div class="sticker menu">
				<button onclick={edit}>Edit</button>
				<button onclick={push}>Push to tomorrow</button>
				<button onclick={drop}>Drop</button>
				{#if movable}
					<details class="sub">
						<summary>Move ›</summary>
						{#each moves as m (m.label)}
							<button disabled={m.to === (tier ?? null)} onclick={() => move(m.to, m.label)}>{m.label}</button>
						{/each}
					</details>
				{/if}
			</div>
		</details>
	{/if}
</article>

<style>
 .focus-icon { display: inline-grid; place-items: center; min-width: 44px; min-height: 44px; color: var(--ink); }
	.execution { display: flex; gap: 8px; align-items: center; flex-wrap: wrap; margin-top: 6px; }
	.step-progress { flex: 1; min-width: 0; display: flex; flex-direction: column; gap: 4px; font-variant-numeric: tabular-nums; }
 .next-action { font-size: 16px; font-weight: 650; overflow-wrap: anywhere; } .next-time, .total-time { font-size: 18px; } .total-time.secondary { font-size: 12px; color: var(--ink-soft); }
 .execution { border-top: 1px solid color-mix(in srgb, var(--ink) 20%, transparent); padding-top: 12px; }
	.card-steps { font-size: 13px; }
	.card-steps summary { cursor: pointer; min-height: 36px; padding-block: 7px; overflow-wrap: anywhere; }
	.start { min-height: 44px; }
	.start:disabled { opacity: .5; cursor: default; }
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
	.badges {
		display: flex;
		gap: 6px;
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
	.feeling {
		font-style: italic;
		font-size: 13px;
		color: var(--ink-soft);
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
	.more {
		position: relative;
	}
	.more summary {
		list-style: none;
		cursor: pointer;
		font-weight: var(--heavy);
		padding: 0 4px;
	}
	.menu {
		position: absolute;
		right: 0;
		top: 24px;
		z-index: 10;
		display: flex;
		flex-direction: column;
		padding: 6px;
		min-width: 160px;
	}
	.menu button {
		background: none;
		border: none;
		text-align: left;
		padding: 6px 8px;
		border-radius: var(--radius-sm);
	}
	.menu button:hover:not(:disabled) {
		background: var(--chip);
	}
	.menu button:disabled {
		opacity: 0.4;
	}
	.sub summary {
		list-style: none;
		cursor: pointer;
		padding: 6px 8px;
		border-radius: var(--radius-sm);
	}
	.sub summary:hover {
		background: var(--chip);
	}
	.sub button {
		display: block;
		width: 100%;
		padding-left: 20px;
	}
</style>
