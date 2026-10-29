<script lang="ts">
 import { api } from '#lib/api.ts';
 import { chat } from '#lib/chat.svelte.ts';
 import { execution } from '#lib/execution.svelte.ts';
 let {taskId}: {taskId: number} = $props();
 let busy = $state(false); let problem = $state('');
 async function divide() {
  busy = true; chat.busy = true; problem = '';
  const before = execution.steps(taskId).length;
  try { const reply = await api.decomposeTask(taskId); chat.message(reply); await execution.refresh(); if (reply.error) problem = reply.text || 'Não foi possível concluir a decomposição.'; else { if (reply.tools?.some((t) => t.error)) problem = 'Parte da decomposição falhou. Confira os passos e a resposta no chat.'; else if (execution.steps(taskId).length <= before) problem = 'A IA não acrescentou passos. Veja a resposta no chat.'; } }
  catch (err) { problem = err instanceof Error ? err.message : String(err); }
  finally { busy = false; chat.busy = false; }
 }
</script>
<button type="button" class="decompose" disabled={busy || chat.busy} onclick={divide}>{busy ? 'Dividindo tarefa…' : 'Dividir em passos com IA'}</button>
{#if problem}<p role="alert">{problem} Você pode adicionar passos manualmente.</p>{/if}
<style>.decompose { border: 0; background: none; padding: 8px 0; min-height: 44px; text-align: left; font-size: 12px; text-decoration: underline; text-underline-offset: 3px; } p { font-size: 13px; background: var(--must); padding: 8px; } button:disabled { opacity: .5; }</style>
