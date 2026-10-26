<script lang="ts">
 import { onMount } from 'svelte';
 import { chat } from '#lib/chat.svelte.ts';
 import { goto } from '$app/navigation';
 import { page } from '$app/state';
 import { api, type FocusSession } from '#lib/api.ts';
 import { execution } from '#lib/execution.svelte.ts';
 import { available, timerText, actionActual, taskActual, taskEstimate, actionClock, focusTarget } from '#lib/execution.ts';
 import { today } from '#lib/today.svelte.ts';
 import { live } from '#lib/live.svelte.ts';
 import { editor } from '#lib/editor.svelte.ts';
 import { toasts } from '#lib/toasts.svelte.ts';
 import ExecutionSteps from '#lib/components/ExecutionSteps.svelte';
 import TaskMetadata from '#lib/components/TaskMetadata.svelte';
 import PlayToggle from '#lib/components/PlayToggle.svelte';
 import { linkify } from '#lib/linkify.ts';
 onMount(() => { chat.open = false; });
 let correctingID = $state<number | null>(null);
 let correctedMinutes = $state<number | undefined>();
 let correctionReason = $state('');
 const selected = $derived(execution.selected);
 const requestedTask = $derived(Number(page.url.searchParams.get('task')) || null);
 const requestedStep = $derived(Number(page.url.searchParams.get('step')) || null);
 const target = $derived(focusTarget(execution.view, requestedTask ? {taskID: requestedTask, stepID: requestedStep ?? (page.url.searchParams.has('direct') ? null : execution.steps(requestedTask).find((s) => !s.done)?.id ?? null)} : null));
 const taskID = $derived(target?.taskID ?? null);
 const task = $derived(today.tasks.find((t) => t.id === taskID));
 const steps = $derived(taskID ? execution.steps(taskID) : []);
 const stepID = $derived(target?.stepID ?? null);
 const step = $derived(execution.view?.steps.find((s) => s.task_id === taskID && s.id === stepID));
 const session = $derived(execution.view?.sessions.find((s) => s.state === 'running' && s.task_id === taskID && s.step_id === stepID) ?? (selected?.task_id === taskID && selected.step_id === stepID ? selected : undefined));
 const allSessions = $derived(execution.view?.sessions ?? []);
 const actual = $derived(taskID ? stepID ? actionActual(allSessions, taskID, stepID, execution.now) : taskActual(allSessions, taskID, execution.now) : 0);
 const aggregate = $derived(taskID ? taskActual(allSessions, taskID, execution.now) : 0);
 const planned = $derived(task ? taskEstimate(task, execution.view?.steps ?? []) : {minutes: null, partial: false});
 const targetEstimate = $derived(step ? step.estimate_minutes : planned.minutes);
 const clock = $derived(actionClock(targetEstimate, actual));
 const progress = $derived(targetEstimate ? Math.min(100, actual / (targetEstimate * 60000) * 100) : null);
 const completed = $derived(steps.filter((s) => s.done).length);
 const allStepsDone = $derived(steps.length > 0 && completed === steps.length);
 const done = $derived(task?.status !== 'open' || !!step?.done || !!step?.deleted_at);
 const eligible = $derived(task ? available(task, today.tasks, today.blockers, today.today?.date ?? '') : false);
 const busy = $derived(execution.busy);
 const history = $derived(allSessions.filter((s) => s.task_id === taskID));
 const correctionSession = $derived(allSessions.find((s) => s.id === correctingID));
 const status = $derived(target?.ready && !done ? 'Próximo passo · parado' : done ? 'Ação concluída' : session?.state === 'running' ? 'Em foco' : session?.state === 'paused' ? 'Pausado' : 'Pronto para começar');
 function choose(id: number | null) { goto(`/focus?task=${taskID}${id ? `&step=${id}` : '&direct=1'}`); }
 async function toggle() {
  if (!task) return;
  if (session?.state === 'running') await execution.action(session.id, 'pause');
  else if (session?.state === 'paused') await execution.action(session.id, 'resume');
  else await execution.start(task.id, stepID);
 }
 function sessionTitle(s: FocusSession) { return execution.view?.steps.find((x) => x.id === s.step_id)?.description ?? task?.title ?? 'Tarefa'; }
 function openCorrection(s: FocusSession) { correctingID = s.id; correctedMinutes = Math.round(execution.elapsed(s) / 600) / 100; correctionReason = ''; }
 async function correctTime(e: SubmitEvent) {
  e.preventDefault();
  if (correctingID && correctedMinutes !== undefined && correctionReason.trim()) {
   if (await execution.correct(correctingID, Math.round(correctedMinutes * 60000), correctionReason.trim())) correctingID = null;
  }
 }
</script>
<svelte:head><title>Foco · whatnow</title></svelte:head>
<nav class="focus-nav" aria-label="Navegação do foco"><a href="/">Ver dia</a><span class="muted">O foco continua ao sair desta tela.</span></nav>
{#if execution.problem}<div class="sticker problem" role="alert"><p>{execution.problem}</p><button class="btn" onclick={() => execution.refresh()}>Tentar atualizar</button></div>{/if}
{#if !live.connected}<p class="connection" role="status">Reconectando. O tempo já confirmado está preservado.</p>{/if}
{#if !execution.view}
 <p>Buscando seu foco…</p>
{:else if !task}
 <section class="sticker empty"><h1>Nenhuma tarefa em foco</h1><a class="btn" href="/">Escolher no meu dia</a></section>
{:else}
 <section class="sticker focus-surface" aria-labelledby="task-title">
  <header class="task-header"><label class="task-check"><input type="checkbox" checked={task.status === 'done'} aria-label={task.status === 'done' ? 'Reabrir tarefa' : 'Concluir tarefa'} disabled={busy} onchange={async () => {if (await toasts.track('Tarefa atualizada', () => api.setStatus(task.id, task.status === 'done' ? 'open' : 'done'))) await Promise.all([today.refresh(), execution.refresh()]);}} /><h1 id="task-title">{task.title}</h1></label><button class="edit-icon" aria-label="Editar tarefa" title="Editar tarefa" onclick={() => editor.open(task.id)}><svg width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" aria-hidden="true"><path d="m15 5 4 4M4 20l4-1L20 7a2.8 2.8 0 0 0-4-4L4 15Z"/></svg></button></header>
  <TaskMetadata {task} />
  <div class="task-facts"><span>{timerText(aggregate)} / {planned.minutes == null ? '—' : timerText(planned.minutes * 60000)}{planned.partial ? ' (parcial)' : ''}</span></div>
  {#if !allStepsDone || step}
  <div class="current-action">
   <div class="action-heading"><h2>{step?.description ?? 'Executar tarefa'}</h2><span class="state">{status}</span></div>
   <div class="timer-controls"><div class="clock"><strong role="timer" aria-label={clock.label}>{clock.text}</strong><span>{clock.label}</span></div>{#if !done}<PlayToggle running={session?.state === 'running'} disabled={busy || !eligible} label={session?.state === 'running' ? 'Pausar foco' : 'Iniciar ou retomar foco'} onclick={() => {void toggle();}} />{/if}</div>
   {#if progress !== null}<progress max="100" value={progress} aria-label="Tempo dedicado em relação à estimativa"></progress>{/if}
   <p class="action-time">{timerText(actual)} / {targetEstimate == null ? '—' : timerText(targetEstimate * 60000)}{clock.label === 'Tempo excedido' ? ' · O registro continua; nada é concluído automaticamente.' : ''}</p>

   {#if !eligible && !done}<p>Esta tarefa está bloqueada ou adiada. Abra a edição para verificar.</p>{/if}

  </div>
  {/if}
  <div class="steps-heading"><h2>Passos</h2><span>{completed}/{steps.length} concluídos</span>{#if step}<button class="text-button" onclick={() => choose(null)}>Executar tarefa principal</button>{/if}</div>
   <ExecutionSteps taskId={task.id} selectedId={step?.id ?? null} {eligible} />
   {#if allStepsDone && task.status === 'open'}<div class="finish-task"><p>Todos os passos concluídos. Marque a tarefa quando estiver pronta.</p></div>{/if}
  {#if task.notes}<div class="task-notes"><h2>Notas</h2><p>{#each linkify(task.notes) as part, i (i)}{#if part.href}<a href={part.href} target="_blank" rel="noopener noreferrer">{part.text}</a>{:else}{part.text}{/if}{/each}</p></div>{/if}
  {#if task.refs?.length}<div class="task-refs">{#each task.refs as ref, i (i)}{#if ref.url}<a href={ref.url} target="_blank" rel="noopener noreferrer">{ref.label}</a>{:else}<span>{ref.label}</span>{/if}{/each}</div>{/if}
  {#if task.raw_text}<details class="raw"><summary>Captura original</summary><p>{task.raw_text}</p></details>{/if}
  {#if history.length}<details class="time-ledger"><summary>Corrigir tempo registrado</summary><p class="muted">A correção é explícita e preserva o histórico. Desfazer tarefas não desfaz tempo.</p>
   {#each history.slice(0, 5) as previous (previous.id)}<div class="ledger-row"><span>{sessionTitle(previous)}<small>{timerText(execution.elapsed(previous))} · {previous.state === 'running' ? 'Em foco' : previous.state === 'paused' ? 'Pausado' : 'Encerrado'}</small></span><button class="text-button" onclick={() => openCorrection(previous)}>Corrigir tempo</button></div>{/each}
  </details>{/if}
  {#if correctionSession}<form class="correction" onsubmit={correctTime}><p>Corrigir este registro de “{sessionTitle(correctionSession)}”.</p><label>Total correto em minutos<input class="field" type="number" min="0" max="525600" step="any" required bind:value={correctedMinutes} /></label><label>Motivo da correção<input class="field" required maxlength="2000" bind:value={correctionReason} /></label><div class="controls"><button class="btn" disabled={busy || correctedMinutes === undefined || correctedMinutes < 0 || !correctionReason.trim()}>Salvar correção</button><button type="button" class="text-button" onclick={() => correctingID = null}>Cancelar correção</button></div></form>{/if}
 </section>
{/if}
<style>
 .focus-nav { display: flex; align-items: center; justify-content: space-between; gap: 12px; font-size: 13px; padding: 4px 0; } .focus-nav a { color: var(--ink); }
 .focus-surface { padding: 20px 24px; max-width: 900px; width: 100%; overflow-wrap: anywhere; }
 .task-header { display: flex; align-items: start; gap: 16px; justify-content: space-between; } h1 { margin: 0; font-size: 26px; line-height: 1.2; font-weight: var(--heavy); }
 .edit-icon { min-width: 44px; min-height: 44px; display: grid; place-items: center; background: none; border: 0; flex-shrink: 0; } .task-check { display: flex; gap: 10px; align-items: flex-start; min-width: 0; } .task-check input { width: 20px; height: 20px; margin-top: 5px; accent-color: var(--primary); flex-shrink: 0; } progress { width: 100%; height: 6px; border: 0; accent-color: var(--primary); margin-top: 10px; }
 .task-facts { display: flex; gap: 8px 16px; flex-wrap: wrap; font-size: 13px; margin-top: 10px; font-variant-numeric: tabular-nums; }
 .current-action { position: sticky; top: 0; z-index: 5; background: var(--paper); border-block: 1px solid color-mix(in srgb, var(--ink) 25%, transparent); padding: 16px 0; margin-top: 18px; } .action-heading { display: flex; gap: 12px; align-items: baseline; flex-wrap: wrap; } h2 { font-size: 17px; margin: 0; font-weight: var(--heavy); } .state { font-size: 12px; color: var(--ink-soft); }
 .timer-controls { display: flex; align-items: center; justify-content: space-between; gap: 20px; flex-wrap: wrap; margin-top: 12px; } .clock { display: flex; flex-direction: column; } .clock strong { font-size: 72px; line-height: 1.1; letter-spacing: -.02em; font-variant-numeric: tabular-nums; } .clock span { font-size: 12px; margin-top: 3px; } .controls { display: flex; gap: 8px; flex-wrap: wrap; } .btn { min-height: 44px; } button:disabled { opacity: .5; cursor: default; }
 .action-time { font-size: 13px; margin: 10px 0 0; font-variant-numeric: tabular-nums; } .text-button { background: none; border: 0; padding: 8px 0; min-height: 44px; text-align: left; text-decoration: underline; text-underline-offset: 3px; font-size: 13px; }
 .steps-heading { display: flex; gap: 10px; align-items: baseline; flex-wrap: wrap; margin-top: 18px; } .steps-heading span { font-size: 12px; } .steps-heading .text-button { margin-left: auto; }
 .finish-task { padding-top: 12px; } .finish-task p { font-size: 13px; }
 .task-refs { display: flex; flex-wrap: wrap; gap: 10px; font-size: 13px; margin-top: 12px; }
 .task-notes { margin-top: 18px; } .task-notes p, .raw p { white-space: pre-wrap; font-size: 14px; max-width: 70ch; margin: 8px 0; } details { margin-top: 12px; font-size: 13px; } summary { cursor: pointer; min-height: 36px; padding-block: 8px; }
 .time-ledger { border-top: 1px solid color-mix(in srgb, var(--ink) 25%, transparent); padding-top: 8px; } .ledger-row { display: flex; gap: 12px; justify-content: space-between; align-items: center; } .ledger-row span { min-width: 0; } small { display: block; color: var(--ink-soft); font-variant-numeric: tabular-nums; }
 .correction { display: flex; flex-direction: column; gap: 10px; max-width: 65ch; } .correction label { display: flex; flex-direction: column; gap: 4px; } .problem { padding: 12px; background: var(--must); } .connection { font-size: 13px; margin: 0; } .empty { padding: 20px; } .empty h1 { margin-bottom: 16px; }
 button:focus-visible, a:focus-visible, summary:focus-visible { outline: 2px solid var(--ink); outline-offset: 3px; }
 @media (max-width: 800px) { .focus-surface { padding: 16px; } .task-header { gap: 6px; } .clock strong { font-size: 60px; } h1 { font-size: 23px; } .focus-nav { align-items: start; } .focus-nav span { max-width: 24ch; text-align: right; } .timer-controls { gap: 14px; } }
</style>
