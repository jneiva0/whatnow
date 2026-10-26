<script lang="ts">
 import { api, type Task, type Size, type Energy } from '#lib/api.ts';
 import { today } from '#lib/today.svelte.ts';
 import { toasts } from '#lib/toasts.svelte.ts';
 let {task}: {task: Task} = $props();
 let busy = $state(false);
 let feelingOpen = $state(false);
 let feeling = $state('');
 async function save(field: 'size' | 'energy' | 'feeling', value: string) {
  busy = true;
  try {
   // Read the latest authoritative task just before this full-field PUT.
   const latest = (await api.tasks()).find((t) => t.id === task.id);
   if (!latest) throw new Error('Tarefa não encontrada');
   const fields = {title: latest.title, notes: latest.notes, size: latest.size, energy: latest.energy, feeling: latest.feeling, due_on: latest.due_on, not_before: latest.not_before, refs: latest.refs ?? [], estimate_minutes: latest.estimate_minutes ?? null};
   if (field === 'size') fields.size = (value || null) as Size | null;
   else if (field === 'energy') fields.energy = (value || null) as Energy | null;
   else fields.feeling = value;
   const result = await api.updateTask(task.id, fields);
   toasts.show('Tarefa atualizada', result.batch_id);
   await today.refresh(); feelingOpen = false;
  } catch (err) { toasts.error(err); }
  finally { busy = false; }
 }
</script>
<div class="metadata">
 <label class="icon-field"><svg width="15" height="15" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" aria-hidden="true"><rect x="4" y="4" width="16" height="16" rx="3"/><path d="M8 8v8M12 8v8M16 8v8"/></svg><select aria-label="Tamanho da tarefa" title="Tamanho da tarefa" value={task.size ?? ''} disabled={busy} onchange={(e) => {const value = e.currentTarget.value; e.currentTarget.value = task.size ?? ''; void save('size', value);}}><option value="">Tamanho</option><option value="S">S</option><option value="M">M</option><option value="L">L</option></select></label>
 <label class="icon-field"><svg width="15" height="15" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" aria-hidden="true"><path d="m13 2-8 12h6l-1 8 9-13h-6Z"/></svg><select aria-label="Energia da tarefa" title="Energia da tarefa" value={task.energy ?? ''} disabled={busy} onchange={(e) => {const value = e.currentTarget.value; e.currentTarget.value = task.energy ?? ''; void save('energy', value);}}><option value="">Energia</option><option value="low">Baixa</option><option value="med">Média</option><option value="high">Alta</option></select></label>
 <button class="feeling-value" disabled={busy} aria-expanded={feelingOpen} onclick={() => {feelingOpen = !feelingOpen; feeling = task.feeling;}}>{task.feeling || 'Sentimento'}</button>
 {#if task.due_on}<span class="due">Prazo: {task.due_on}</span>{/if}
</div>
{#if feelingOpen}<form class="feeling-edit" onsubmit={(e) => {e.preventDefault(); void save('feeling', feeling);}}><input class="field" aria-label="Sentimento sobre a tarefa" bind:value={feeling} /><button class="btn" disabled={busy}>Salvar</button><button type="button" class="cancel" onclick={() => feelingOpen = false}>Cancelar</button></form>{/if}
<style>.icon-field { display: inline-flex; align-items: center; gap: 3px; color: var(--ink-soft); }.metadata { display: flex; align-items: center; gap: 2px 8px; flex-wrap: wrap; font-size: 12px; margin-top: 3px; } select, .feeling-value { min-height: 44px; max-width: 100%; background: transparent; border: 0; padding: 4px 3px; color: var(--ink-soft); } .feeling-value { text-align: left; text-decoration: underline; text-underline-offset: 3px; overflow-wrap: anywhere; } .due { color: var(--ink-soft); } .feeling-edit { display: flex; flex-wrap: wrap; gap: 8px; padding-block: 8px; } .feeling-edit input { flex: 1; min-width: 140px; } .cancel { border: 0; background: none; min-height: 44px; } select:focus-visible, button:focus-visible { outline: 2px solid var(--ink); outline-offset: 2px; }</style>
