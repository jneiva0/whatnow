<script lang="ts">
 import { Button } from '#lib/components/ui/button/index.js';
 import PlayIcon from '@lucide/svelte/icons/play';
 import PauseIcon from '@lucide/svelte/icons/pause';
 let {running = false, disabled = false, label, onclick}: {running?: boolean; disabled?: boolean; label: string; onclick: () => void} = $props();
</script>
<Button variant={running ? 'playback-active' : 'playback'} size="playback" {disabled} aria-pressed={running} aria-label={label} title={label} {onclick}>
 {#if running}<PauseIcon aria-hidden="true" />{:else}<PlayIcon aria-hidden="true" />{/if}
</Button>
