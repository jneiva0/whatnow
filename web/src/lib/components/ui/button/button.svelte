<script lang="ts" module>
	import { type VariantProps, tv } from "tailwind-variants";
	import { cn, type WithElementRef } from "#lib/utils.js";
	import type { HTMLAnchorAttributes, HTMLButtonAttributes } from "svelte/elements";

	export const buttonVariants = tv({
		base: "group/button inline-flex shrink-0 items-center justify-center whitespace-nowrap rounded-md [border-width:var(--stroke)] border-transparent bg-clip-padding text-sm [font-weight:var(--heavy)] transition-[background-color,box-shadow,transform] outline-none focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-ring aria-invalid:border-destructive select-none disabled:pointer-events-none disabled:opacity-50 [&_svg]:pointer-events-none [&_svg]:size-[var(--icon-size,1.125rem)] [&_svg]:shrink-0 motion-reduce:transition-none",
		variants: {
			variant: {
				default: "border-border bg-primary text-primary-foreground shadow-control hover:bg-primary/90 active:translate-x-0.5 active:translate-y-0.5 active:shadow-none",
				outline: "border-border bg-card text-card-foreground shadow-control hover:bg-accent hover:text-accent-foreground active:translate-x-0.5 active:translate-y-0.5 active:shadow-none",
				secondary: "border-border bg-secondary text-secondary-foreground shadow-control hover:bg-accent hover:text-accent-foreground active:translate-x-0.5 active:translate-y-0.5 active:shadow-none",
				metadata: "border-border/25 bg-card/60 text-foreground [font-weight:500] hover:border-border hover:bg-accent aria-expanded:bg-accent",
				ghost: "text-foreground hover:bg-accent hover:text-accent-foreground aria-expanded:bg-accent",
				destructive: "border-border bg-destructive text-destructive-foreground shadow-control hover:bg-destructive/90",
				link: "text-primary underline-offset-4 hover:underline",
				inline: "w-full justify-start whitespace-normal wrap-anywhere text-left [font:inherit] text-foreground cursor-text hover:bg-muted hover:border-input",
				"inline-number": "w-full justify-end whitespace-nowrap text-right [font:inherit] text-foreground cursor-text hover:bg-muted hover:border-input tabular-nums",
				playback: "rounded-[var(--pill)] bg-[var(--play-background,transparent)] text-[var(--play-ink,var(--foreground))] hover:bg-accent",
				"playback-active": "rounded-[var(--pill)] bg-[var(--play-background,var(--accent))] text-[var(--play-ink,var(--foreground))] hover:bg-accent",
			},
			size: {
				default: "min-h-11 min-w-11 gap-2 px-3 py-1.5",
				xs: "h-8 gap-1 px-2 text-xs",
				sm: "h-9 gap-1 px-3",
				lg: "min-h-12 gap-2 px-4",
				icon: "size-11",
				"icon-xs": "size-8",
				"icon-sm": "size-9",
				"icon-lg": "size-12",
				playback: "size-[var(--play-size,2.75rem)]",
				inline: "h-auto min-h-11 min-w-0 px-1.5 py-1.5",
			},
		},
		defaultVariants: {
			variant: "default",
			size: "default",
		},
	});

	export type ButtonVariant = VariantProps<typeof buttonVariants>["variant"];
	export type ButtonSize = VariantProps<typeof buttonVariants>["size"];

	export type ButtonProps = WithElementRef<HTMLButtonAttributes> &
		WithElementRef<HTMLAnchorAttributes> & {
			variant?: ButtonVariant;
			size?: ButtonSize;
		};
</script>

<script lang="ts">
	let {
		class: className,
		variant = "default",
		size = "default",
		ref = $bindable(null),
		href = undefined,
		type = "button",
		disabled,
		children,
		...restProps
	}: ButtonProps = $props();
</script>

{#if href}
	<a
		bind:this={ref}
		data-slot="button"
		class={cn(buttonVariants({ variant, size }), className)}
		href={disabled ? undefined : href}
		aria-disabled={disabled}
		role={disabled ? "link" : undefined}
		tabindex={disabled ? -1 : undefined}
		{...restProps}
	>
		{@render children?.()}
	</a>
{:else}
	<button
		bind:this={ref}
		data-slot="button"
		class={cn(buttonVariants({ variant, size }), className)}
		{type}
		{disabled}
		{...restProps}
	>
		{@render children?.()}
	</button>
{/if}
