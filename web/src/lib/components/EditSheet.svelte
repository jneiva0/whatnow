<script lang="ts">
	import { editor } from '#lib/editor.svelte.ts';
	import { today } from '#lib/today.svelte.ts';
	import TaskEditor from './TaskEditor.svelte';
	import Toasts from './Toasts.svelte';

	const task = $derived(today.tasks.find((t) => t.id === editor.taskId));
</script>



{#if task || editor.creating}
	<dialog class={['sticker sheet', {existing: !!task}]} aria-label={task ? 'Edit task' : 'New task'} oncancel={() => editor.close()} {@attach (node) => { node.showModal(); return () => node.close(); }}>
		{#key task?.id ?? 'new'}
			<TaskEditor task={task ?? null} />
		{/key}
		<Toasts />
	</dialog>
{/if}

<style>
 .sheet { margin: auto; width: min(560px, calc(100vw - 32px)); max-width: none; max-height: 90dvh; overflow-y: auto; background: var(--paper); padding: 0 24px; color: var(--ink); }
 .sheet.existing { width: min(1040px, calc(100vw - 48px)); }
 .sheet::backdrop { background: color-mix(in srgb, var(--ink) 40%, transparent); }
 @media(max-width: 600px) { .sheet, .sheet.existing { width: 100%; max-height: 100dvh; height: 100dvh; border-radius: 0; border: 0; padding-inline: 16px; } }
</style>
