<script lang="ts">
 import { api, type Task, type Energy } from '#lib/api.ts';
 import InlineEdit from './InlineEdit.svelte';
 import { execution } from '#lib/execution.svelte.ts';
 import { taskEstimate } from '#lib/execution.ts';
 import { today } from '#lib/today.svelte.ts';
 import { toasts } from '#lib/toasts.svelte.ts';
 let {task}: {task: Task} = $props();
 let busy = $state(false);
 let feelingOpen = $state(false);
 let feeling = $state('');
 const estimate = $derived(taskEstimate(task, execution.view?.steps ?? []));
 async function save(field: 'estimate_minutes' | 'energy' | 'feeling', value: string) {
  busy = true;
  try {
   // Read the latest authoritative task just before this full-field PUT.
   const latest = (await api.tasks()).find((t) => t.id === task.id);
   if (!latest) throw new Error('Tarefa não encontrada');
   const fields = {title: latest.title, notes: latest.notes, size: latest.size, energy: latest.energy, feeling: latest.feeling, due_on: latest.due_on, not_before: latest.not_before, refs: latest.refs ?? [], estimate_minutes: latest.estimate_minutes ?? null};
   if (field === 'estimate_minutes') fields.estimate_minutes = value ? Number(value) : null;
   else if (field === 'energy') fields.energy = (value || null) as Energy | null;
   else fields.feeling = value;
   const result = await api.updateTask(task.id, fields);
   toasts.show('Tarefa atualizada', result.batch_id);
   await today.refresh(); await execution.refresh(); feelingOpen = false;
  } catch (err) { if (field === 'estimate_minutes') throw err; toasts.error(err); }
  finally { busy = false; }
 }
</script>
<div class="metadata">
 <div class="icon-field estimate-field"><svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" aria-hidden="true"><circle cx="12" cy="12" r="9"/><path d="M12 7v5l3 2"/></svg><InlineEdit numeric value={task.estimate_minutes == null ? '' : String(task.estimate_minutes)} display={estimate.minutes == null ? '— min' : `${estimate.minutes}${estimate.partial ? '+' : ''} min`} label="Editar estimativa da tarefa em minutos" disabled={busy} onsave={(value) => save('estimate_minutes', value)} /></div>
 <label class="icon-field"><svg width="15" height="15" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" aria-hidden="true"><path d="m13 2-8 12h6l-1 8 9-13h-6Z"/></svg><select aria-label="Energia da tarefa" title="Energia da tarefa" value={task.energy ?? ''} disabled={busy} onchange={(e) => {const value = e.currentTarget.value; e.currentTarget.value = task.energy ?? ''; void save('energy', value);}}><option value="">Energia</option><option value="low">Baixa</option><option value="med">Média</option><option value="high">Alta</option></select></label>
 <button class="feeling-value" disabled={busy} aria-expanded={feelingOpen} onclick={() => {feelingOpen = !feelingOpen; feeling = task.feeling;}}>{task.feeling || 'Sentimento'}</button>
 {#if task.due_on}<span class="due">Prazo: {task.due_on}</span>{/if}
</div>
{#if feelingOpen}<form class="feeling-edit" onsubmit={(e) => {e.preventDefault(); void save('feeling', feeling);}}><input class="field" aria-label="Sentimento sobre a tarefa" bind:value={feeling} /><button class="btn" disabled={busy}>Salvar</button><button type="button" class="cancel" onclick={() => feelingOpen = false}>Cancelar</button></form>{/if}
<style>.estimate-field { width: 110px; } .icon-field { border: 1px solid color-mix(in srgb, var(--ink) 20%, transparent); border-radius: var(--radius-sm); padding-inline: 6px; background: color-mix(in srgb, var(--card) 55%, transparent); } .icon-field { display: inline-flex; align-items: center; gap: 3px; color: var(--ink-soft); }.metadata { display: flex; align-items: center; gap: 2px 8px; flex-wrap: wrap; font-size: 12px; margin-top: 3px; } select, .feeling-value { min-height: 44px; max-width: 100%; background: transparent; border: 0; padding: 4px 3px; color: var(--ink-soft); } .feeling-value { text-align: left; text-decoration: underline; text-underline-offset: 3px; overflow-wrap: anywhere; } .due { color: var(--ink-soft); } .feeling-edit { display: flex; flex-wrap: wrap; gap: 8px; padding-block: 8px; } .feeling-edit input { flex: 1; min-width: 140px; } .cancel { border: 0; background: none; min-height: 44px; } select:focus-visible, button:focus-visible { outline: 2px solid var(--ink); outline-offset: 2px; }</style>
