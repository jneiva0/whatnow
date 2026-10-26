<script lang="ts">
 import { goto } from '$app/navigation';
 import PlayToggle from './PlayToggle.svelte';
 import TaskDecompose from './TaskDecompose.svelte';
 import { api, type Step, type Mutation } from '#lib/api.ts';
 import { execution } from '#lib/execution.svelte.ts';
 import { actionActual, timerText, reorderStepIDs } from '#lib/execution.ts';
 import { toasts } from '#lib/toasts.svelte.ts';
 let {taskId, selectedId = null, eligible = true}: {taskId: number; selectedId?: number | null; eligible?: boolean} = $props();
 let saving = $state(false); let problem = $state(''); let description = $state('');
 let editing = $state<number | null>(null); let text = $state('');
 let estimating = $state<number | null>(null); let minutes = $state<number | undefined>();
 let splitting = $state<number | null>(null); let splitText = $state('');
 let dragged = $state<number | null>(null); let over = $state<number | null>(null); let announcement = $state('');
 const steps = $derived(execution.steps(taskId));
 async function write(label: string, op: () => Promise<Mutation>) {
  if (saving) return false; saving = true; problem = '';
  try { const result = await op(); await execution.refresh(); toasts.show(label, result.batch_id); return true; }
  catch (err) { problem = `${err instanceof Error ? err.message : String(err)}. A alteração não foi salva.`; return false; }
  finally { saving = false; }
 }
 function fields(step: Step) { return {description: step.description, estimate_minutes: step.estimate_minutes, done: step.done}; }
 async function toggle(step: Step) { await write(step.done ? 'Passo reaberto' : 'Passo concluído', () => api.updateStep(taskId, step.id, {...fields(step), done: !step.done})); }
 async function saveText(step: Step) { if (text.trim() && await write('Passo salvo', () => api.updateStep(taskId, step.id, {...fields(step), description: text.trim()}))) editing = null; }
 async function saveEstimate(step: Step) { if (await write('Estimativa salva', () => api.updateStep(taskId, step.id, {...fields(step), estimate_minutes: minutes ?? null}))) estimating = null; }
 async function add(e: SubmitEvent) { e.preventDefault(); if (description.trim() && await write('Passo adicionado', () => api.addStep(taskId, {description: description.trim(), estimate_minutes: null}))) description = ''; }
 async function split(step: Step) { const parts = splitText.split('\n').map((s) => s.trim()).filter(Boolean); if (parts.length < 2) { problem = 'Escreva pelo menos dois passos, um por linha.'; return; } if (await write('Passo dividido', () => api.splitStep(taskId, step.id, parts.map((description) => ({description, estimate_minutes: null}))))) splitting = null; }
 async function start(step: Step) { if (await execution.start(taskId, step.id)) goto(`/focus?task=${taskId}&step=${step.id}`); }
 async function reorder(id: number, destination: number) {
  const ids = reorderStepIDs(steps.map((s) => s.id), id, destination); if (!ids) return;
  const to = ids.indexOf(id);
  if (await write('Ordem dos passos salva', () => api.reorderSteps(taskId, ids))) announcement = `Passo movido para posição ${to + 1} de ${ids.length}.`;
 }
 function pointerDown(e: PointerEvent, id: number) {
  if (saving || e.button !== 0) return; const grip = e.currentTarget as HTMLButtonElement;
  grip.setPointerCapture(e.pointerId); dragged = id; over = id;
 }
 function pointerMove(e: PointerEvent) {
  if (dragged === null) return;
  const row = document.elementFromPoint(e.clientX, e.clientY)?.closest<HTMLElement>('[data-step-row]');
  if (row?.dataset.task === String(taskId)) over = Number(row.dataset.stepRow);
 }
 function pointerUp() { if (dragged !== null && over !== null) void reorder(dragged, over); dragged = null; over = null; }
 function keyboard(e: KeyboardEvent, id: number) { const index = steps.findIndex((s) => s.id === id); const delta = e.key === 'ArrowUp' ? -1 : e.key === 'ArrowDown' ? 1 : 0; if (delta) {e.preventDefault(); const dest = steps[index + delta]; if (dest) void reorder(id, dest.id);} }
</script>
<p class="sr-only" aria-live="polite">{announcement}</p>
{#if problem}<p class="problem" role="alert">{problem}</p>{/if}
<ol class="step-list">
 {#each steps as step (step.id)}
  {@const actual = actionActual(execution.view?.sessions ?? [], taskId, step.id, execution.now)}
  {@const running = execution.view?.sessions.find((s) => s.task_id === taskId && s.step_id === step.id && s.state === 'running')}
  <li data-step-row={step.id} data-task={taskId} class={{selected: selectedId === step.id || running, done: step.done, dragging: dragged === step.id, over: dragged !== null && over === step.id}}>
   <div class="step-row">
    <button type="button" class="grip" aria-label={`Reordenar ${step.description}. Use setas para cima ou baixo.`} disabled={saving} onpointerdown={(e) => pointerDown(e, step.id)} onpointermove={pointerMove} onpointerup={pointerUp} onpointercancel={() => {dragged = null; over = null;}} onkeydown={(e) => keyboard(e, step.id)}><svg width="14" height="18" viewBox="0 0 14 18" fill="currentColor" aria-hidden="true"><circle cx="4" cy="4" r="1.5"/><circle cx="10" cy="4" r="1.5"/><circle cx="4" cy="9" r="1.5"/><circle cx="10" cy="9" r="1.5"/><circle cx="4" cy="14" r="1.5"/><circle cx="10" cy="14" r="1.5"/></svg></button>
    <button type="button" class="check" aria-label={`${step.done ? 'Reabrir' : 'Concluir'} ${step.description}`} aria-pressed={step.done} disabled={saving || execution.busy} onclick={() => toggle(step)}><svg viewBox="0 0 24 24" width="20" height="20" aria-hidden="true"><rect x="3" y="3" width="18" height="18" rx="4" fill="none" stroke="currentColor" stroke-width="2" />{#if step.done}<path d="m7 12 3 3 7-7" fill="none" stroke="currentColor" stroke-width="2" />{/if}</svg></button>
    {#if editing === step.id}<form class="text-edit" onsubmit={(e) => {e.preventDefault(); void saveText(step);}}><input class="field" aria-label="Descrição do passo" bind:value={text} onkeydown={(e) => {if (e.key === 'Escape') editing = null;}} /><button disabled={saving}>Salvar</button></form>{:else}<button type="button" class="description" title="Editar passo" onclick={() => {editing = step.id; text = step.description;}}>{step.description}</button>{/if}
    <span class="actual">{timerText(actual)}</span><span class="divider">/</span>
    {#if estimating === step.id}<form class="estimate-edit" onsubmit={(e) => {e.preventDefault(); void saveEstimate(step);}}><input class="field" type="number" min="1" max="1440" step="1" aria-label="Estimativa do passo em minutos" bind:value={minutes} /><button disabled={saving}>Salvar</button></form>{:else}<button type="button" class="estimate" aria-label={`Editar estimativa de ${step.description}`} title="Editar estimativa em minutos" onclick={() => {estimating = step.id; minutes = step.estimate_minutes ?? undefined;}}>{step.estimate_minutes == null ? '—' : timerText(step.estimate_minutes * 60000)}</button>{/if}
    {#if !step.done}<PlayToggle running={!!running} disabled={!eligible || saving || execution.busy} label={`${running ? 'Pausar' : 'Começar'} ${step.description}`} onclick={() => {if (running) void execution.action(running.id, 'pause'); else void start(step);}} />{/if}
    <details class="more"><summary aria-label={`Mais ações para ${step.description}`}>⋯</summary><div class="step-menu">{#if !step.done}<button type="button" onclick={() => {splitting = step.id; splitText = '';}}>Dividir este passo</button>{/if}<button type="button" disabled={saving} onclick={() => write('Passo removido', () => api.deleteStep(taskId, step.id))}>Remover passo</button></div></details>
   </div>
   {#if splitting === step.id}<form class="split" onsubmit={(e) => {e.preventDefault(); void split(step);}}><label>Substituir este passo por ações menores<textarea class="field" rows="3" bind:value={splitText}></textarea></label><p>Um passo por linha. Substitui só esta ação pendente; preserva tempo e passos concluídos.</p><button disabled={saving}>Dividir passo</button><button type="button" onclick={() => splitting = null}>Cancelar</button></form>{/if}
  </li>
 {/each}
</ol>
<form class="add-step" onsubmit={add}><input class="field" aria-label="Novo passo" placeholder="+ Adicionar passo" bind:value={description} /><button disabled={saving || !description.trim()} aria-label="Adicionar passo">Adicionar</button></form>
<TaskDecompose {taskId} />
<style>
 .step-list { list-style: none; padding: 0; margin: 8px 0 0; } li { border-top: 1px solid color-mix(in srgb, var(--ink) 20%, transparent); border-radius: var(--radius-sm); } .step-row { display: flex; align-items: center; gap: 2px; min-width: 0; } .selected { background: color-mix(in srgb, var(--chip) 35%, transparent); } .dragging { opacity: .6; } .over { outline: 2px dashed var(--primary); }
 button { background: none; border: 0; min-height: 44px; } .grip { touch-action: none; cursor: grab; padding: 4px; width: 24px; flex-shrink: 0; } .check { display: grid; place-items: center; min-width: 36px; padding: 6px; } .description { flex: 1; min-width: 40px; text-align: left; padding: 6px 2px; overflow-wrap: anywhere; } .done .description { text-decoration: line-through; }
 .actual,.estimate { font-size: 14px; font-variant-numeric: tabular-nums; white-space: nowrap; } .actual { margin-left: auto; } .estimate { padding: 4px; text-decoration: underline; text-underline-offset: 3px; } .divider { font-size: 12px; color: var(--ink-soft); }
 .text-edit { flex: 1; min-width: 100px; display: flex; flex-wrap: wrap; } .estimate-edit { max-width: 100px; } .estimate-edit input { padding: 4px; } .text-edit button,.estimate-edit button { font-size: 12px; }
 .more { position: relative; } summary { list-style: none; cursor: pointer; min-width: 28px; min-height: 44px; display: grid; place-items: center; } .step-menu { position: absolute; right: 0; top: 40px; z-index: 15; background: var(--card); border: var(--line); border-radius: var(--radius-sm); padding: 4px 8px; min-width: 160px; display: flex; flex-direction: column; } .step-menu button { text-align: left; font-size: 12px; }
 .add-step { display: flex; gap: 8px; margin-top: 8px; } .add-step input { flex: 1; min-width: 60px; padding: 6px; font-size: 13px; } .add-step button { font-size: 12px; }
 .split { padding: 8px; } .split label { display: flex; flex-direction: column; font-size: 13px; gap: 4px; } .split p { font-size: 12px; } .split button { margin-right: 12px; text-decoration: underline; }
 .problem { font-size: 13px; background: var(--must); padding: 8px; } .sr-only { position: absolute; width: 1px; height: 1px; overflow: hidden; clip-path: inset(50%); } button:disabled { opacity: .5; cursor: default; } button:focus-visible,summary:focus-visible { outline: 2px solid var(--ink); outline-offset: 2px; }
 @media(max-width: 480px) { .step-row { flex-wrap: wrap; } .description { min-width: calc(100% - 68px); } .actual { margin-left: 60px; } }
</style>
