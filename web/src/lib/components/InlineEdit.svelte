<script lang="ts">
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
<span class={['inline-edit', {numeric}]} aria-busy={pending}>
 {#if editing}
  <input {@attach (node) => { node.focus(); node.select(); }} bind:value={draft} aria-label={label} aria-invalid={!!error} inputmode={numeric ? 'numeric' : 'text'} readonly={pending} onkeydown={keydown} onblur={() => {void save();}} />
 {:else}
  <button type="button" {disabled} aria-label={label} title={label} onclick={open}>{display ?? value}</button>
 {/if}
 {#if pending}<span class="pending" role="status" aria-label="Salvando">…</span>{/if}
 {#if error}<span class="error" role="alert">{error}</span>{/if}
</span>
<style>
 .inline-edit { display: block; position: relative; min-width: 0; width: 100%; }
 button, input { box-sizing: border-box; min-height: 44px; width: 100%; font: inherit; color: var(--ink); border: 1px solid transparent; border-radius: var(--radius-sm); padding: 6px; text-align: left; background: transparent; }
 button { overflow-wrap: anywhere; cursor: text; }
 button:hover { background: color-mix(in srgb, var(--card) 60%, transparent); border-color: var(--ink-soft); }
 input { min-width: 0; background: var(--card); border-color: var(--ink); }
 .numeric { font-variant-numeric: tabular-nums; }
 .numeric button, .numeric input { text-align: right; }
 button:focus-visible, input:focus-visible { outline: 2px solid var(--primary); outline-offset: 2px; }
 button:disabled { opacity: .5; cursor: default; }
 .pending { position: absolute; right: 4px; top: -8px; font-size: 12px; }
 .error { display: block; min-width: 0; overflow-wrap: anywhere; font-size: 12px; color: var(--danger); }
</style>
