<script lang="ts">
 let {running = false, disabled = false, label, onclick}: {running?: boolean; disabled?: boolean; label: string; onclick: () => void} = $props();
</script>
<button class={['play-toggle', {running}]} {disabled} aria-label={label} title={label} {onclick}><svg width="20" height="20" viewBox="0 0 24 24" aria-hidden="true">{#if running}<path d="M8 5v14M16 5v14" stroke="currentColor" stroke-width="3" stroke-linecap="round" />{:else}<path d="m9 5 11 7-11 7Z" fill="currentColor" />{/if}</svg></button>
<style>.play-toggle { min-width: 44px; height: 44px; display: inline-grid; place-items: center; border: 0; border-radius: var(--pill); background: transparent; color: var(--ink); flex-shrink: 0; } .play-toggle:hover { background: var(--chip); } .running { color: var(--primary); background: color-mix(in srgb, var(--primary) 10%, transparent); } button:disabled { opacity: .45; cursor: default; } button:focus-visible { outline: 2px solid var(--ink); outline-offset: 2px; }</style>
