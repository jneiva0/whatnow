<script lang="ts">
	import { execution } from '#lib/execution.svelte.ts';
	import { today } from '#lib/today.svelte.ts';
	import { actionActual, actionClock, timerText } from '#lib/execution.ts';
	import { toasts } from '#lib/toasts.svelte.ts';
	import { Progress } from '#lib/components/ui/progress/index.js';
	import PlayToggle from './PlayToggle.svelte';

	const running = $derived(execution.view?.sessions.find((session) => session.state === 'running'));
	const task = $derived(today.tasks.find((task) => task.id === running?.task_id));
	const step = $derived(execution.view?.steps.find((step) => step.id === running?.step_id));
	const title = $derived(step?.description ?? task?.title ?? 'Foco em andamento');
	const actual = $derived(running ? actionActual(execution.view?.sessions ?? [], running.task_id, running.step_id, execution.now) : 0);
	const estimate = $derived(running?.step_id != null ? step?.estimate_minutes : running ? execution.estimate(running.task_id).minutes : null);
	const clock = $derived(actionClock(estimate, actual));
	const progress = $derived(estimate ? Math.min(100, actual / (estimate * 60000) * 100) : null);
	const href = $derived(running ? `/focus?task=${running.task_id}&${running.step_id === null ? 'direct=1' : `step=${running.step_id}`}` : '/focus');

	async function pause() {
		if (!running || execution.busy) return;
		if (!await execution.action(running.id, 'pause')) toasts.error(execution.problem);
	}
</script>

{#if running}
	<aside class="active-focus" aria-label="Foco em andamento">
		<a {href} class="focus-link" aria-label={`Voltar ao foco: ${title}`}>
			<strong>Em foco</strong>
			<span class="action-title">{title}</span>
			<span class="timing"><span class="clock" aria-label={`${clock.label}: ${clock.text}`}>{clock.text}</span><span class="time-divider" aria-hidden="true">|</span><span class="time-pair" aria-label="Tempo dedicado / estimativa">{timerText(actual)} / {estimate == null ? '—' : timerText(estimate * 60000)}</span></span>
		</a>
		<PlayToggle running disabled={execution.busy} label={`Pausar ${title}`} onclick={pause} />
		{#if progress !== null}<Progress value={progress} max={100} class="absolute inset-x-0 bottom-0 h-1 rounded-none" aria-label="Tempo dedicado em relação à estimativa" />{/if}
	</aside>
{/if}

<style>
	.active-focus { position: sticky; top: 0; z-index: 15; display: flex; align-items: center; gap: 8px; padding: max(4px, env(safe-area-inset-top)) 16px 6px; background: var(--paper);  color: var(--ink); }
	.focus-link { display: flex; flex: 1; min-width: 0; align-items: center; gap: 12px; min-height: 44px; color: inherit; text-decoration: none; }
	.focus-link:hover .action-title { text-decoration: underline; }
	.focus-link:focus-visible { outline: 2px solid var(--primary); outline-offset: 2px; }
	strong { flex-shrink: 0; font-size: 13px; color: var(--primary); }
	.action-title { flex: 1; min-width: 0; overflow: hidden; text-overflow: ellipsis; white-space: nowrap; }
	.timing { flex-shrink: 0; display: flex; align-items: baseline; gap: 8px; white-space: nowrap; font-variant-numeric: tabular-nums; }
	.time-pair { font-size: 12px; color: var(--ink-soft); }
	.time-divider { color: var(--ink-soft); font-weight: 400; }
	.clock { font-size: 20px; flex-shrink: 0; font-variant-numeric: tabular-nums; font-weight: var(--heavy); }
	@media (max-width: 480px) { .active-focus { padding-inline: 12px; } .focus-link { gap: 8px; } strong { display: none; } }
</style>
