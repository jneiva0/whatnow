<script lang="ts">
	import { onDestroy, untrack } from 'svelte';
	import { suggestion } from '#lib/suggestion.svelte.ts';
	import { goto } from '$app/navigation';
	import { page } from '$app/state';
	import { api, type Tier } from '#lib/api.ts';
	import { live } from '#lib/live.svelte.ts';
	import { actionCards, doneOn, tierLabels, tiers, unplanned } from '#lib/plan.ts';
	import { execution } from '#lib/execution.svelte.ts';
	import { survived } from '#lib/time.ts';
	import { toasts } from '#lib/toasts.svelte.ts';
	import { editor } from '#lib/editor.svelte.ts';
	import { today } from '#lib/today.svelte.ts';
	import { welcomeDue } from '#lib/welcome.ts';
	import Capture from '#lib/components/Capture.svelte';
	import DayStrip from '#lib/components/DayStrip.svelte';
	import DaysDone from '#lib/components/DaysDone.svelte';
	import DoneToday from '#lib/components/DoneToday.svelte';
	import Header from '#lib/components/Header.svelte';
	import TaskCard from '#lib/components/TaskCard.svelte';
	import WhatNow from '#lib/components/WhatNow.svelte';
	import Working from '#lib/components/Working.svelte';
	import { work } from '#lib/work.svelte.ts';

	let interacting = $state(false);
 let editingField = $state(false);
 function isEditing(target: EventTarget | null) {
  return target instanceof Element && !!target.closest('input, textarea, [contenteditable="true"]');
 }
	const asking = $derived(suggestion.asking);
	let replanning = $state(false);


	const t = $derived(today.today);
	const planned = $derived(t ? tiers(t, today.tasks) : null);
	const tray = $derived(t ? unplanned(t, today.tasks) : []);
	const done = $derived(t ? doneOn(t.date, today.tasks) : []);
	const pick = $derived(suggestion.pick);
 const pickTask = $derived(pick ? today.tasks.find((x) => x.id === pick.task_id) : undefined);
 const running = $derived(execution.view?.sessions.find((s) => s.state === 'running'));
 const runningTask = $derived(today.tasks.find((task) => task.id === running?.task_id));
 const runningStep = $derived(execution.view?.steps.find((step) => step.id === running?.step_id));
	const actionView = $derived(page.url.searchParams.get('view') === 'actions');
	const actions = $derived(t ? actionCards(t, today.tasks, today.blockers, execution.view?.sessions.find((s) => s.state === 'running')?.task_id) : []);
	const tierNames: Tier[] = ['must', 'should', 'could'];
	// A replan can come from anywhere: this button, a pace change, the agent, or an edit that matters.
	const replanBusy = $derived(replanning || work.busy('replan'));

 $effect(() => {
  const snapshot = t && execution.view ? {day: t, tasks: today.tasks, blockers: today.blockers, execution: execution.view, events: today.events, revision: today.revision} : null;
  const active = actionView;
  const editing = interacting || editingField || editor.taskId !== null || editor.creating;
  const busy = execution.busy || work.busy('cleanup') || replanBusy;
  if (snapshot) untrack(() => suggestion.sync(snapshot, active, editing, busy));
 });
 onDestroy(() => suggestion.deactivate());

	// A day with no plan starts with the morning check-in, once first-run setup has had its turn.
	$effect(() => {
		if (t && !t.plan && live.view && !welcomeDue(live.view.settings, '/')) goto('/morning', { replaceState: true });
	});

	// A nudge's What now? arrives as ?whatnow, from any page.
	$effect(() => {
		if (t?.plan && page.url.searchParams.has('whatnow')) {
			goto(actionView ? '/?view=actions' : '/', { replaceState: true });
			if (!t.day?.wrapped_at) whatNow();
		}
	});

	// The new plan arrives over the stream as plan.updated, which refreshes the page.
	async function replan() {
		replanning = true;
		try {
			await api.replan();
		} catch (err) {
			toasts.error(err);
		} finally {
			replanning = false;
		}
	}

 function whatNow(exclude = 0, context = '') {
  if (running) { void goto(`/focus?task=${running.task_id}&${running.step_id == null ? 'direct=1' : `step=${running.step_id}`}`); return; }
  void suggestion.ask(exclude, context);
 }

</script>

<svelte:document onfocusin={(event) => editingField = isEditing(event.target)} onfocusout={(event) => editingField = isEditing(event.relatedTarget)} />

{#snippet suggestionPanel(inline: boolean)}
 <WhatNow {pick} task={pickTask} {asking} problem={suggestion.problem} stale={suggestion.stale} nothing={suggestion.nothing} {inline}
  onask={(context) => whatNow(0, context)} onnotthis={(context) => whatNow(pick?.task_id ?? 0, context)}
  onclose={() => suggestion.dismiss()} oninteraction={(value) => interacting = value} />
{/snippet}

{#if t?.plan}
	<Header today={t} {asking} onwhatnow={() => whatNow()} />
	<DayStrip events={today.events} now={live.now} endOfDay={live.endOfDay} startedAt={t.day?.started_at ?? null} />
	<div class="add">
		<Capture />
		<button class="btn" onclick={() => editor.openNew()}>+ Task</button>
	</div>

 {#if !actionView && !running && (pick || suggestion.asking || suggestion.nothing || suggestion.problem)}
  {@render suggestionPanel(false)}
 {/if}

	{#if t.day?.wrapped_at}
		<DaysDone day={t.day} done={done.length} meetings={survived(today.events, live.now)} />
	{:else if planned}
		<nav class="day-view" aria-label="Visualização do dia">
			<a href="/" aria-current={!actionView ? 'page' : undefined}>Prioridades</a>
			<a href="/?view=actions" aria-current={actionView ? 'page' : undefined}>Ações</a>
		</nav>
		{#if replanBusy}<Working label="Replanning…" />{/if}
		{#if t.plan.fallback}
			<p class="chip rough">rough plan, model was out</p>
		{/if}
		{#if actionView}
			<section aria-label="Próximas ações">
    {#if running}
     <div class="sticker continuation">
      <h2>Continue de onde está</h2>
      <p>{runningStep?.description ?? runningTask?.title ?? 'Foco em andamento'}</p>
      <a class="btn" href={`/focus?task=${running.task_id}&${running.step_id == null ? 'direct=1' : `step=${running.step_id}`}&from=actions`}>Voltar ao foco</a>
     </div>
    {:else}
     {@render suggestionPanel(true)}
    {/if}
    <p class="action-hint">Outras ações disponíveis · você também pode escolher por aqui.</p>
				{#each actions.filter((item) => running || !pick || item.card.task.id !== pick.task_id) as item (item.card.task.id)}
					<TaskCard card={item.card} tier={item.tier} date={t.date} actionView />
				{:else}
					<p class="muted empty">Sem outras ações disponíveis agora.</p>
				{/each}
			</section>
		{:else}
		{#each tierNames as tier (tier)}
			<section>
				<h2 class="label">{tierLabels[tier]}</h2>
				{#each planned[tier] as card (card.task.id)}
					<TaskCard {card} {tier} date={t.date} />
				{:else}
					<p class="muted empty">Nothing here.</p>
				{/each}
			</section>
		{/each}
		{/if}
	{/if}

	<div class={['lower', {actions: actionView && !t.day?.wrapped_at}]}>
		{#if !actionView || t.day?.wrapped_at}
		<section>
			<div class="tray-head">
				<h2 class="label">Unplanned</h2>
				{#if tray.length && !t.day?.wrapped_at}
					<button class="btn" disabled={replanBusy} onclick={replan}>
						{replanBusy ? 'Replanning…' : 'Replan'}
					</button>
				{/if}
			</div>
			{#each tray as card (card.task.id)}
				<TaskCard {card} date={t.date} />
			{:else}
				<p class="muted empty">Tray's empty.</p>
			{/each}
		</section>
		{/if}
		<section>
			<DoneToday {done} meetings={survived(today.events, live.now)} />
		</section>
	</div>
{:else}
	<p class="muted">Loading…</p>
{/if}

<style>
	.day-view { display: flex; gap: 4px; align-self: flex-start; padding: 4px; border: var(--line); border-radius: var(--radius-sm); background: var(--card); }
	.day-view a { display: inline-flex; align-items: center; min-height: 44px; padding: 6px 16px; color: var(--ink); text-decoration: none; border-radius: calc(var(--radius-sm) / 2); }
	.day-view a[aria-current='page'] { background: var(--chip); }
	.day-view a:focus-visible { outline: 2px solid var(--primary); outline-offset: 2px; }
	.continuation { padding: 12px 16px; margin-bottom: 12px; background: var(--chip); }
	.continuation h2 { font-size: 18px; margin: 0; }
	.continuation p { margin: 4px 0 8px; overflow-wrap: anywhere; }
	.action-hint { margin: 0 0 12px; color: var(--ink-soft); font-size: 13px; }
	.lower.actions { grid-template-columns: minmax(0, 1fr); }
	h2 {
		margin: 6px 0 8px;
	}
	.empty {
		margin: 0 0 8px;
		font-size: 14px;
	}
	.rough {
		align-self: flex-start;
		margin: 0;
	}
	.add {
		display: flex;
		gap: 10px;
		align-items: center;
	}
	.add > :global(:first-child) {
		flex: 1;
		min-width: 0;
	}
	.tray-head {
		display: flex;
		align-items: center;
		justify-content: space-between;
		margin-bottom: 1rem;
	}
	.lower {
		display: grid;
		grid-template-columns: repeat(2, minmax(0, 1fr));
		gap: 16px;
		align-items: start;
	}

	@media (max-width: 800px) {
		.lower {
			grid-template-columns: minmax(0, 1fr);
		}
	}
</style>
