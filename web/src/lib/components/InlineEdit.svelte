<script lang="ts">
 import { Button } from '#lib/components/ui/button/index.js';
 import { Input } from '#lib/components/ui/input/index.js';
 import { cn } from '#lib/utils.js';
 let {value, display, label, numeric = false, disabled = false, onsave}: {
  value: string; display?: string; label: string; numeric?: boolean; disabled?: boolean;
  onsave: (value: string) => Promise<void>;
 } = $props();
 let editing = $state(false);
 let draft = $state('');
 let pending = $state(false);
 let error = $state('');
 function open() {
  draft = value; error = ''; editing = true;
 }
 async function save() {
  if (!editing || pending) return;
  const next = draft.trim();
  if (next === value) { editing = false; return; }
  if (numeric && next !== '' && (!/^\d+$/.test(next) || Number(next) < 1 || Number(next) > 1440)) {
   error = 'Use minutos inteiros entre 1 e 1440.'; return;
  }
  if (!numeric && !next) { error = 'Escreva uma ação.'; return; }
  pending = true; error = '';
  try { await onsave(next); editing = false; }
  catch (err) { error = `${err instanceof Error ? err.message : String(err)}. Não salvo; pressione Enter para tentar novamente.`; }
  finally { pending = false; }
 }
 function keydown(e: KeyboardEvent) {
  if (e.isComposing) return;
  if (e.key === 'Escape' && !pending) { e.preventDefault(); editing = false; error = ''; }
  if (e.key === 'Enter') { e.preventDefault(); void save(); }
 }
</script>
<span class={cn('relative block w-full min-w-0', numeric && 'tabular-nums')} aria-busy={pending}>
 {#if editing}
  <Input variant="inline" textAlign={numeric ? 'end' : 'start'} {@attach (node: HTMLInputElement) => { node.focus(); node.select(); }} bind:value={draft} aria-label={label} aria-invalid={!!error} inputmode={numeric ? 'numeric' : 'text'} readonly={pending} onkeydown={keydown} onblur={() => {void save();}} />
 {:else}
  <Button variant={numeric ? 'inline-number' : 'inline'} size="inline" {disabled} aria-label={label} title={label} onclick={open}>{display ?? value}</Button>
 {/if}
 {#if pending}<span class="absolute -top-2 right-1 text-xs" role="status" aria-label="Salvando">…</span>{/if}
 {#if error}<span class="block min-w-0 text-xs wrap-anywhere text-destructive" role="alert">{error}</span>{/if}
</span>
