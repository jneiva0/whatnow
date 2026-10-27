<script lang="ts">
	import * as DropdownMenu from '#lib/components/ui/dropdown-menu/index.js';
	import * as ContextMenu from '#lib/components/ui/context-menu/index.js';
	import PencilIcon from '@lucide/svelte/icons/pencil';
	import CalendarArrowUpIcon from '@lucide/svelte/icons/calendar-arrow-up';
	import ArchiveIcon from '@lucide/svelte/icons/archive';
	import { api, type Task, type Tier } from '#lib/api.ts';
	import { editor } from '#lib/editor.svelte.ts';
	import { tierLabels } from '#lib/plan.ts';
	import { addDays } from '#lib/time.ts';
	import { toasts } from '#lib/toasts.svelte.ts';
	import { today } from '#lib/today.svelte.ts';

	let { task, tier, date, listed = false, context = false }: {
		task: Task;
		tier?: Tier;
		date: string;
		listed?: boolean;
		context?: boolean;
	} = $props();

	// The two menu primitives share one set of labels, eligibility and commands.
	const Menu = $derived(context ? ContextMenu : DropdownMenu);
	const movable = $derived(!listed && !!today.today?.plan && !today.today.day?.wrapped_at);
	const moves: { to: Tier | null; label: string }[] = [
		{ to: 'must', label: tierLabels.must },
		{ to: 'should', label: tierLabels.should },
		{ to: 'could', label: tierLabels.could },
		{ to: null, label: 'Unplanned' }
	];
</script>

<Menu.Group>
	<Menu.Item onSelect={() => editor.open(task.id)}><PencilIcon />Editar tarefa</Menu.Item>
	<Menu.Item onSelect={() => { void toasts.track(`Pushed to tomorrow: ${task.title}`, () => api.push(task.id, addDays(date, 1))); }}>
		<CalendarArrowUpIcon />Adiar para amanhã
	</Menu.Item>
	{#if movable}
		<Menu.Sub>
			<Menu.SubTrigger>Mover para</Menu.SubTrigger>
			<Menu.SubContent>
				<Menu.Group>
					{#each moves as move (move.label)}
						<Menu.Item disabled={move.to === (tier ?? null)} onSelect={() => { void toasts.track(`Moved to ${move.label}: ${task.title}`, () => api.move(task.id, move.to)); }}>{move.label}</Menu.Item>
					{/each}
				</Menu.Group>
			</Menu.SubContent>
		</Menu.Sub>
	{/if}
</Menu.Group>
<Menu.Separator />
<Menu.Group>
	<Menu.Item variant="destructive" onSelect={() => { void toasts.track(`Dropped: ${task.title}`, () => api.setStatus(task.id, 'dropped')); }}>
		<ArchiveIcon />Descartar tarefa
	</Menu.Item>
</Menu.Group>
