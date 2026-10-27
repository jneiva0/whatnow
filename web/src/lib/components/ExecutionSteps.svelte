<script lang="ts">
 import { startAndOpenFocus } from '#lib/focus-navigation.ts';
 import InlineEdit from './InlineEdit.svelte';
 import PlayToggle from './PlayToggle.svelte';
 import TaskDecompose from './TaskDecompose.svelte';
 import * as DropdownMenu from '#lib/components/ui/dropdown-menu/index.js';
 import { buttonVariants } from '#lib/components/ui/button/index.js';
 import EllipsisIcon from '@lucide/svelte/icons/ellipsis';
 import ScissorsIcon from '@lucide/svelte/icons/scissors';
 import TrashIcon from '@lucide/svelte/icons/trash-2';
 import { api, type Step, type Mutation } from '#lib/api.ts';
 import { execution } from '#lib/execution.svelte.ts';
 import { actionActual, timerText, reorderStepIDs } from '#lib/execution.ts';
 import { toasts } from '#lib/toasts.svelte.ts';
 let {taskId, selectedId = null, eligible = true, preview = false}: {taskId: number; selectedId?: number | null; eligible?: boolean; preview?: boolean} = $props();
 let showAll = $state(false);
 let saving = $state(false); let problem = $state(''); let description = $state('');
 let splitting = $state<number | null>(null); let splitText = $state('');
 let dragged = $state<number | null>(null); let over = $state<number | null>(null); let announcement = $state('');
 const steps = $derived(execution.steps(taskId));
 const current = $derived(steps.find((step) => execution.view?.sessions.some((s) => s.task_id === taskId && s.step_id === step.id && s.state === 'running')) ?? steps.find((step) => !step.done));
 const compact = $derived(preview && !showAll);
 const visibleSteps = $derived(compact ? current ? [current] : [] : steps);
 async function write(label: string, op: () => Promise<Mutation>) {
  if (saving) return false; saving = true; problem = '';
  try { const result = await op(); await execution.refresh(); toasts.show(label, result.batch_id); return true; }
  catch (err) { problem = `${err instanceof Error ? err.message : String(err)}. A alteração não foi salva.`; return false; }
  finally { saving = false; }
 }
 function fields(step: Step) { return {description: step.description, estimate_minutes: step.estimate_minutes, done: step.done}; }
 async function toggle(step: Step) { await write(step.done ? 'Passo reaberto' : 'Passo concluído', () => api.updateStep(taskId, step.id, {...fields(step), done: !step.done})); }
 async function editField(step: Step, change: {description?: string; estimate_minutes?: number | null}) {
  // Re-read the row after editing, so a completion arriving from another device is preserved.
  const latest = (await api.execution()).steps.find((s) => s.id === step.id && s.task_id === taskId && !s.deleted_at);
  if (!latest) throw new Error('Passo não encontrado');
  const result = await api.updateStep(taskId, step.id, {...fields(latest), ...change});
  await execution.refresh(); toasts.show('Passo atualizado', result.batch_id);
 }
 async function add(e: SubmitEvent) { e.preventDefault(); if (description.trim() && await write('Passo adicionado', () => api.addStep(taskId, {description: description.trim(), estimate_minutes: null}))) description = ''; }
 async function split(step: Step) { const parts = splitText.split('\n').map((s) => s.trim()).filter(Boolean); if (parts.length < 2) { problem = 'Escreva pelo menos dois passos, um por linha.'; return; } if (await write('Passo dividido', () => api.splitStep(taskId, step.id, parts.map((description) => ({description, estimate_minutes: null}))))) splitting = null; }
 async function toggleFocus(step: Step) {
  if (execution.busy) return;
  problem = '';
  const running = execution.view?.sessions.find((s) => s.task_id === taskId && s.step_id === step.id && s.state === 'running');
  if (running) {
   if (!await execution.action(running.id, 'pause')) problem = execution.problem;
  } else {
   problem = await startAndOpenFocus({task_id: taskId, step_id: step.id}, preview) ?? '';
  }
 }

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
<ol class={['step-list', {compact}]}>
 {#each visibleSteps as step (step.id)}
  {@const actual = actionActual(execution.view?.sessions ?? [], taskId, step.id, execution.now)}
  {@const running = execution.view?.sessions.find((s) => s.task_id === taskId && s.step_id === step.id && s.state === 'running')}
  <li data-step-row={step.id} data-task={taskId} class={{selected: selectedId === step.id || running, done: step.done, dragging: dragged === step.id, over: dragged !== null && over === step.id}}>
   <div class="step-row">
    <button type="button" class="grip" aria-label={`Reordenar ${step.description}. Use setas para cima ou baixo.`} disabled={saving} onpointerdown={(e) => pointerDown(e, step.id)} onpointermove={pointerMove} onpointerup={pointerUp} onpointercancel={() => {dragged = null; over = null;}} onkeydown={(e) => keyboard(e, step.id)}><svg width="14" height="18" viewBox="0 0 14 18" fill="currentColor" aria-hidden="true"><circle cx="4" cy="4" r="1.5"/><circle cx="10" cy="4" r="1.5"/><circle cx="4" cy="9" r="1.5"/><circle cx="10" cy="9" r="1.5"/><circle cx="4" cy="14" r="1.5"/><circle cx="10" cy="14" r="1.5"/></svg></button>
    <button type="button" class="check" aria-label={`${step.done ? 'Reabrir' : 'Concluir'} ${step.description}`} aria-pressed={step.done} disabled={saving || execution.busy} onclick={() => toggle(step)}><svg viewBox="0 0 24 24" width="20" height="20" aria-hidden="true"><rect x="3" y="3" width="18" height="18" rx="4" fill="none" stroke="currentColor" stroke-width="2" />{#if step.done}<path d="m7 12 3 3 7-7" fill="none" stroke="currentColor" stroke-width="2" />{/if}</svg></button>
    <div class="description"><InlineEdit value={step.description} label={`Editar ${step.description}`} disabled={saving} onsave={(description) => editField(step, {description})} /></div>
    <div class="time-pair"><span class="actual">{timerText(actual)}</span><span class="divider">/</span><div class="estimate"><InlineEdit numeric value={step.estimate_minutes == null ? '' : String(step.estimate_minutes)} display={step.estimate_minutes == null ? '—' : timerText(step.estimate_minutes * 60000)} label={`Estimativa de ${step.description} em minutos`} disabled={saving} onsave={(value) => editField(step, {estimate_minutes: value ? Number(value) : null})} /></div></div>
    {#if !step.done}<PlayToggle running={!!running} disabled={(!eligible && !running) || saving || execution.busy} label={`${running ? 'Pausar' : 'Começar'} ${step.description}`} onclick={() => {void toggleFocus(step);}} />{/if}
    <div class="more">
     <DropdownMenu.Root>
      <DropdownMenu.Trigger class={buttonVariants({variant: 'ghost', size: 'icon'})} aria-label={`Mais ações para ${step.description}`} disabled={saving}>
       <EllipsisIcon aria-hidden="true" />
      </DropdownMenu.Trigger>
      <DropdownMenu.Content align="end">
       <DropdownMenu.Group>
        {#if !step.done}<DropdownMenu.Item onSelect={() => {splitting = step.id; splitText = '';}}><ScissorsIcon />Dividir este passo</DropdownMenu.Item>{/if}
        <DropdownMenu.Item variant="destructive" onSelect={() => {void write('Passo removido', () => api.deleteStep(taskId, step.id));}}><TrashIcon />Remover passo</DropdownMenu.Item>
       </DropdownMenu.Group>
      </DropdownMenu.Content>
     </DropdownMenu.Root>
    </div>
   </div>
   {#if splitting === step.id}<form class="split" onsubmit={(e) => {e.preventDefault(); void split(step);}}><label>Substituir este passo por ações menores<textarea class="field" rows="3" bind:value={splitText}></textarea></label><p>Um passo por linha. Substitui só esta ação pendente; preserva tempo e passos concluídos.</p><button disabled={saving}>Dividir passo</button><button type="button" onclick={() => splitting = null}>Cancelar</button></form>{/if}
  </li>
 {/each}
</ol>
{#if preview}<button type="button" class="show-steps" aria-expanded={showAll} onclick={() => showAll = !showAll}>{showAll ? 'Mostrar só a próxima ação' : `Ver todos os passos · ${steps.filter((s) => s.done).length}/${steps.length}`}</button>{/if}
{#if !compact}
<form class="add-step" onsubmit={add}><input class="field" aria-label="Novo passo" placeholder="+ Adicionar passo" bind:value={description} /><button disabled={saving || !description.trim()} aria-label="Adicionar passo">Adicionar</button></form>
<TaskDecompose {taskId} />
{/if}
<style>
 .compact .step-row { grid-template-columns: 36px minmax(0, 1fr) auto 44px 44px; }
 .compact .grip { display: none; }
 .compact .description { font-size: 19px; font-weight: var(--heavy); line-height: 1.4; }
 .compact .time-pair { font-size: 13px; }
 @media (max-width: 600px) {
  .compact .step-row { grid-template-columns: 36px minmax(0, 1fr) 44px 44px; }
  .compact .description { grid-column: 2 / -1; }
  .compact .time-pair { grid-column: 2; grid-row: 2; }
 }

 .compact .grip { visibility: hidden; pointer-events: none; } .show-steps { font-size: 12px; color: var(--ink-soft); text-decoration: underline; text-underline-offset: 3px; padding-inline: 0; }
 .step-list { list-style: none; padding: 0; margin: 8px 0 0; } li { border-top: 1px solid color-mix(in srgb, var(--ink) 20%, transparent); border-radius: var(--radius-sm); } .step-row { display: grid; grid-template-columns: 24px 36px minmax(0, 1fr) auto 44px 44px; align-items: center; gap: 4px; min-width: 0; padding-block: 4px; } .selected { background: color-mix(in srgb, var(--chip) 35%, transparent); } .dragging { opacity: .6; } .over { outline: 2px dashed var(--primary); }
 button { background: none; border: 0; min-height: 44px; } .grip { touch-action: none; cursor: grab; padding: 4px; width: 24px; flex-shrink: 0; } .check { display: grid; place-items: center; min-width: 36px; padding: 6px; } .description { flex: 1; min-width: 40px; text-align: left; padding: 6px 2px; overflow-wrap: anywhere; } .done .description { text-decoration: line-through; }
 .time-pair { display: flex; align-items: center; gap: 2px; font-size: 14px; font-variant-numeric: tabular-nums; }
 .actual { white-space: nowrap; } .estimate { width: 66px; } .divider { color: var(--ink-soft); }

 .add-step { display: flex; gap: 8px; margin-top: 8px; } .add-step input { flex: 1; min-width: 60px; padding: 6px; font-size: 13px; } .add-step button { font-size: 12px; }
 .split { padding: 8px; } .split label { display: flex; flex-direction: column; font-size: 13px; gap: 4px; } .split p { font-size: 12px; } .split button { margin-right: 12px; text-decoration: underline; }
 .problem { font-size: 13px; background: var(--must); padding: 8px; } .sr-only { position: absolute; width: 1px; height: 1px; overflow: hidden; clip-path: inset(50%); } button:disabled { opacity: .5; cursor: default; } button:focus-visible { outline: 2px solid var(--ink); outline-offset: 2px; }
 .more { grid-column: 6; } .description { min-width: 0; }
 @media(max-width: 600px) { .step-row { grid-template-columns: 24px 36px minmax(0, 1fr) 44px 44px; } .description { grid-column: 3 / 6; } .time-pair { grid-column: 3; justify-self: end; } .more { grid-column: 5; } }
</style>
