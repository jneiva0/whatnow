<script lang="ts">
	import { api, type Task, type TaskFields } from '#lib/api.ts';
	import * as Select from '#lib/components/ui/select/index.js';
	import * as Popover from '#lib/components/ui/popover/index.js';
	import * as Field from '#lib/components/ui/field/index.js';
	import { Button, buttonVariants } from '#lib/components/ui/button/index.js';
	import { Textarea } from '#lib/components/ui/textarea/index.js';
	import ClockIcon from '@lucide/svelte/icons/clock';
	import ZapIcon from '@lucide/svelte/icons/zap';
	import BatteryLowIcon from '@lucide/svelte/icons/battery-low';
	import BatteryMediumIcon from '@lucide/svelte/icons/battery-medium';
	import BatteryFullIcon from '@lucide/svelte/icons/battery-full';
	import MessageCircleIcon from '@lucide/svelte/icons/message-circle';
	import CalendarIcon from '@lucide/svelte/icons/calendar';
	import CheckIcon from '@lucide/svelte/icons/check';
	import XIcon from '@lucide/svelte/icons/x';
	import InlineEdit from './InlineEdit.svelte';
	import { execution } from '#lib/execution.svelte.ts';
	import { taskEstimate } from '#lib/execution.ts';
	import { today } from '#lib/today.svelte.ts';
	import { toasts } from '#lib/toasts.svelte.ts';
	import { cn } from '#lib/utils.js';

	let { task }: { task: Task } = $props();
	const uid = $props.id();
	let busy = $state(false);
	let feelingOpen = $state(false);
	let feeling = $state('');
	let feelingError = $state('');
	let energyError = $state('');
	const estimate = $derived(taskEstimate(task, execution.view?.steps ?? []));
	const energyOptions = [
		{ value: 'none', label: 'Não definida', icon: ZapIcon },
		{ value: 'low', label: 'Baixa', icon: BatteryLowIcon },
		{ value: 'med', label: 'Média', icon: BatteryMediumIcon },
		{ value: 'high', label: 'Alta', icon: BatteryFullIcon }
	];
	const energy = $derived(energyOptions.find((option) => option.value === task.energy) ?? energyOptions[0]);

	async function save(patch: Partial<Pick<TaskFields, 'estimate_minutes' | 'energy' | 'feeling'>>) {
		if (busy) throw new Error('Aguarde o salvamento em andamento.');
		busy = true;
		try {
			// This API replaces all fields. Refetch first to preserve edits from other devices.
			const latest = (await api.tasks()).find((entry) => entry.id === task.id);
			if (!latest) throw new Error('Tarefa não encontrada.');
			const fields = {
				title: latest.title, notes: latest.notes, size: latest.size,
				energy: latest.energy, feeling: latest.feeling, due_on: latest.due_on,
				not_before: latest.not_before, refs: latest.refs ?? [],
				estimate_minutes: latest.estimate_minutes ?? null, ...patch
			};
			const result = await api.updateTask(task.id, fields);
			toasts.show('Tarefa atualizada', result.batch_id);
			// A refresh failure must not misreport a confirmed write as unsaved.
			await Promise.all([today.refresh(), execution.refresh()]).catch(() => {
				toasts.error('Alteração salva, mas a tela não atualizou. Recarregue a página.');
			});
		} finally {
			busy = false;
		}
	}

	async function changeEnergy(value: string) {
		const next = value === 'none' ? null : value;
		if (next === task.energy || busy) return;
		energyError = '';
		try {
			// The bound getter keeps the server value selected until confirmation.
			await save({ energy: next as Task['energy'] });
		} catch (err) {
			energyError = `Energia não salva: ${err instanceof Error ? err.message : String(err)}. Selecione novamente para tentar.`;
		}
	}

	async function saveFeeling() {
		if (busy) return;
		if (feeling === task.feeling) { feelingOpen = false; return; }
		feelingError = '';
		try {
			await save({ feeling });
			feelingOpen = false;
		} catch (err) {
			feelingError = `Não salvo: ${err instanceof Error ? err.message : String(err)}. Seu texto continua aqui.`;
		}
	}
</script>

<div class="mt-1 flex min-w-0 flex-wrap items-center gap-x-2 gap-y-1 text-sm" aria-busy={busy}>
	<div class="flex w-28 min-w-0 items-center gap-1 rounded-md border border-border/25 bg-card/60 pr-1 pl-2">
		<ClockIcon class="size-4 shrink-0" aria-hidden="true" />
		<InlineEdit
			numeric value={task.estimate_minutes == null ? '' : String(task.estimate_minutes)}
			display={estimate.minutes == null ? '— min' : `${estimate.minutes}${estimate.partial ? '+' : ''} min`}
			label="Editar estimativa da tarefa em minutos" disabled={busy}
			onsave={(value) => save({ estimate_minutes: value ? Number(value) : null })}
		/>
	</div>
	<Select.Root type="single" items={energyOptions} disabled={busy}
		bind:value={() => task.energy ?? 'none', (value) => { void changeEnergy(value); }}>
		<Select.Trigger variant="metadata" aria-label={`Energia da tarefa: ${energy.label}`} title="Editar energia" aria-invalid={!!energyError}>
			<energy.icon aria-hidden="true" />
			<span>{task.energy ? energy.label : 'Energia'}</span>
		</Select.Trigger>
		<Select.Content align="start">
			<Select.Group>
				<Select.GroupHeading>Energia da tarefa</Select.GroupHeading>
				{#each energyOptions as option (option.value)}
					<Select.Item value={option.value} label={option.label}>
						<option.icon aria-hidden="true" />{option.label}
					</Select.Item>
				{/each}
			</Select.Group>
		</Select.Content>
	</Select.Root>
	<Popover.Root bind:open={feelingOpen} onOpenChange={(open) => {
		if (open) { feeling = task.feeling; feelingError = ''; }
	}}>
		<Popover.Trigger class={cn(buttonVariants({ variant: 'metadata' }), 'max-w-full min-w-0')}
			disabled={busy} title={task.feeling || 'Adicionar sentimento'} aria-label="Editar sentimento sobre a tarefa">
			<MessageCircleIcon aria-hidden="true" />
			<span class="max-w-64 truncate">{task.feeling || 'Sentimento'}</span>
		</Popover.Trigger>
		<Popover.Content align="start" aria-labelledby={`${uid}-feeling-label`}>
			<form onsubmit={(event) => { event.preventDefault(); void saveFeeling(); }}>
				<Field.Group class="gap-3">
					<Field.Field data-invalid={!!feelingError} data-disabled={busy}>
						<Field.Label id={`${uid}-feeling-label`} for={`${uid}-feeling`}>Como essa tarefa parece?</Field.Label>
						<Textarea id={`${uid}-feeling`} bind:value={feeling} readonly={busy}
							aria-invalid={!!feelingError} aria-describedby={feelingError ? `${uid}-feeling-error` : undefined}
							onkeydown={(event) => {
								if (event.key === 'Enter' && (event.ctrlKey || event.metaKey) && !event.isComposing) {
									event.preventDefault(); void saveFeeling();
								}
							}} />
						{#if feelingError}<Field.Error id={`${uid}-feeling-error`}>{feelingError}</Field.Error>{/if}
					</Field.Field>
					<div class="flex justify-end gap-2">
						<Popover.Close class={buttonVariants({ variant: 'ghost', size: 'icon' })} disabled={busy} aria-label="Cancelar edição do sentimento" title="Cancelar">
							<XIcon aria-hidden="true" />
						</Popover.Close>
						<Button type="submit" size="icon" disabled={busy} aria-label="Salvar sentimento" title="Salvar sentimento"><CheckIcon aria-hidden="true" /></Button>
					</div>
				</Field.Group>
			</form>
		</Popover.Content>
	</Popover.Root>
	{#if task.due_on}<span class="inline-flex items-center gap-1 text-muted-foreground"><CalendarIcon class="size-4 shrink-0" aria-hidden="true" />{task.due_on}</span>{/if}
	{#if busy}<span class="text-xs text-muted-foreground" role="status">Salvando…</span>{/if}
</div>
{#if energyError}<p class="mt-1 text-sm text-destructive" role="alert">{energyError}</p>{/if}
