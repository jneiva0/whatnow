<script lang="ts">
	import { cn, type WithElementRef } from "#lib/utils.js";
	import type { HTMLInputAttributes, HTMLInputTypeAttribute } from "svelte/elements";

	type InputType = Exclude<HTMLInputTypeAttribute, "file">;

	type Props = WithElementRef<
		Omit<HTMLInputAttributes, "type"> &
			({ type: "file"; files?: FileList } | { type?: InputType; files?: undefined })
	> & { variant?: "default" | "inline"; textAlign?: "start" | "end" };

	let {
		ref = $bindable(null),
		value = $bindable(),
		type,
		files = $bindable(),
		class: className,
		variant = "default",
		textAlign = "start",
		"data-slot": dataSlot = "input",
		...restProps
	}: Props = $props();
</script>

{#if type === "file"}
	<input
		bind:this={ref}
		data-slot={dataSlot}
		class={cn(
			"min-h-11 w-full min-w-0 rounded-md [border-width:var(--stroke)] border-input bg-card px-3 py-1.5 text-base text-card-foreground outline-none focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-ring aria-invalid:border-destructive file:border-0 file:bg-transparent file:font-semibold file:text-foreground placeholder:text-muted-foreground disabled:cursor-not-allowed disabled:opacity-50",
			variant === "inline" && "[font:inherit] px-1.5",
			textAlign === "end" && "text-right tabular-nums",
			className
		)}
		type="file"
		bind:files
		bind:value
		{...restProps}
	/>
{:else}
	<input
		bind:this={ref}
		data-slot={dataSlot}
		class={cn(
			"min-h-11 w-full min-w-0 rounded-md [border-width:var(--stroke)] border-input bg-card px-3 py-1.5 text-base text-card-foreground outline-none focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-ring aria-invalid:border-destructive file:border-0 file:bg-transparent file:font-semibold file:text-foreground placeholder:text-muted-foreground disabled:cursor-not-allowed disabled:opacity-50",
			variant === "inline" && "[font:inherit] px-1.5",
			textAlign === "end" && "text-right tabular-nums",
			className
		)}
		{type}
		bind:value
		{...restProps}
	/>
{/if}
