<script lang="ts">
 import { onDestroy } from 'svelte';
 import { startAndOpenFocus } from '#lib/focus-navigation.ts';
 import { execution } from '#lib/execution.svelte.ts';
 import { chat } from '#lib/chat.svelte.ts';
 import { editor } from '#lib/editor.svelte.ts';
 import { timerText, actionActual } from '#lib/execution.ts';
 import type { Pick, Task } from '#lib/api.ts';
 import { Button } from '#lib/components/ui/button/index.js';
 import { Input } from '#lib/components/ui/input/index.js';
 import * as Field from '#lib/components/ui/field/index.js';
 import XIcon from '@lucide/svelte/icons/x';

 let {pick, task, asking, problem = '', stale = false, nothing = false, inline = false, onask, onnotthis, onclose, oninteraction}: {
  pick: Pick | null; task: Task | undefined; asking: boolean; problem?: string; stale?: boolean; nothing?: boolean; inline?: boolean;
  onask: (context: string) => void; onnotthis: (context: string) => void; onclose: () => void; oninteraction: (value: boolean) => void;
 } = $props();
 const uid = $props.id();
 let contextOpen = $state(false);
 let context = $state('');
 let starting = $state(false);
 let startProblem = $state('');
 const step = $derived(execution.view?.steps.find((s) => s.id === pick?.step_id));
 const minutes = $derived(pick?.step_id != null ? step?.estimate_minutes : task ? execution.estimate(task.id).minutes : null);
 const partial = $derived(pick?.step_id == null && task ? execution.estimate(task.id).partial : false);
 const actual = $derived(pick ? actionActual(execution.view?.sessions ?? [], pick.task_id, pick.step_id ?? null, execution.now) : 0);

 function lock(value: boolean) { contextOpen = value; oninteraction(value || starting); }
 onDestroy(() => oninteraction(false));
 async function begin() {
  if (!pick || starting || execution.busy || asking) return;
  starting = true; startProblem = ''; oninteraction(true);
  try {
   const error = await startAndOpenFocus({task_id: pick.task_id, step_id: pick.step_id ?? null}, inline);
   if (error) startProblem = error;
   else onclose();
  } finally { starting = false; oninteraction(contextOpen); }
 }
 function refine(event: SubmitEvent) {
  event.preventDefault();
  const note = context.trim();
  onask(note);
  lock(false); context = '';
 }
</script>

<section class="sticker whatnow" aria-label="Sugestão para começar" aria-busy={asking}>
 <div class="suggestion-heading">
  <h2>{pick ? (step?.description ?? task?.title ?? 'Ação sugerida') : nothing ? 'Nenhuma ação disponível agora' : asking ? 'Escolhendo uma ação…' : 'Por onde começar?'}</h2>
  {#if !inline}<Button variant="ghost" size="icon" disabled={starting} aria-label="Fechar sugestão" onclick={onclose}><XIcon aria-hidden="true" /></Button>{/if}
 </div>
 {#if pick}
  {#if step}<p class="parent">{task?.title}</p>{/if}
  <p class="estimate">{timerText(actual)} / {minutes == null ? '—' : `${timerText(minutes * 60000)}${partial ? '+' : ''}`} <span>dedicado / estimativa</span></p>
  <p class="reason">{pick.reason}</p>
  {#if pick.fallback}<p class="fallback">Sugestão básica por prazo; a IA não pôde avaliar o contexto.</p>{/if}
 {/if}
 {#if asking}<p class="status" role="status">Buscando uma sugestão com os dados atuais…</p>{/if}
 {#if problem}<p class="problem" role="alert">{problem}</p>{/if}
 {#if startProblem}<p class="problem" role="alert">{startProblem}</p>{/if}
 <div class="actions">
  {#if pick}
   <Button disabled={starting || execution.busy || asking} onclick={begin}>{starting ? 'Iniciando…' : 'Começar'}</Button>
   <Button variant="outline" disabled={asking || starting} onclick={() => onnotthis('')}>Outra sugestão</Button>
   <Button variant="ghost" disabled={starting} onclick={() => editor.open(pick!.task_id)}>Ver tarefa</Button>
  {:else if !asking}
   <Button disabled={execution.busy} onclick={() => onask('')}>{problem ? 'Tentar novamente' : 'Buscar sugestão'}</Button>
  {/if}
  {#if stale && pick && !asking}<Button variant="ghost" disabled={starting} onclick={() => onask('')}>Atualizar sugestão</Button>{/if}
  <Button variant="ghost" disabled={starting || asking || chat.busy || !!chat.draft} onclick={() => { void chat.helpOrganize(pick?.task_id); }}>Me ajude a organizar</Button>
  <Button variant="ghost" disabled={starting || asking} aria-expanded={contextOpen} onclick={() => lock(!contextOpen)}>Ajustar contexto</Button>
 </div>
 {#if contextOpen}
  <form onsubmit={refine}>
   <Field.Group>
    <Field.Field>
     <Field.Label for={`${uid}-context`}>O que considerar agora?</Field.Label>
     <Input id={`${uid}-context`} bind:value={context} placeholder="Estou travado, tenho dez minutos…" readonly={asking} />
     <Field.Description>Vale para esta consulta, sem alterar os dados da tarefa.</Field.Description>
    </Field.Field>
    <div><Button type="submit" disabled={asking || starting}>Usar este contexto</Button></div>
   </Field.Group>
  </form>
 {/if}
</section>

<style>
 .whatnow { background: var(--chip); padding: 12px 16px; margin-bottom: 12px; }
 .suggestion-heading { display: flex; align-items: center; justify-content: space-between; gap: 8px; }
 h2 { margin: 0; font-size: 20px; font-weight: var(--heavy); overflow-wrap: anywhere; }
 p { margin: 6px 0; overflow-wrap: anywhere; }
 .parent, .estimate, .fallback, .status { font-size: 13px; }
 .parent, .estimate span { color: var(--ink-soft); }
 .estimate { font-variant-numeric: tabular-nums; }
 .estimate span { margin-left: 6px; }
 .actions { display: flex; flex-wrap: wrap; gap: 4px 8px; margin-top: 8px; }
 .problem { color: var(--destructive); }
 form { margin-top: 12px; }
</style>
