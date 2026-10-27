# UI migration to shadcn-svelte

## Direction

Use the installed shadcn-svelte CLI and skill as the common component foundation.
Keep Whatnow's Sticker identity: cream canvas, dark ink, tier fills, purple primary
actions, chunky outlines and offset shadows. Existing alternative themes remain
available. A separate theme package is not required.

The source skill is `web/.agents/skills/shadcn-svelte/SKILL.md`. The repository-root
`.agents/skills/shadcn-svelte` symlink points there, and root/frontend `AGENTS.md`
explicitly require reading it. This makes the entry point available from either
working directory; a running agent may need a new session to refresh its automatic
skill catalog. The checked-in instructions work without relying on that refresh.

## Theme compatibility

`web/src/routes/layout.css` is the single stylesheet entry point configured in
`web/components.json`. It loads Tailwind, the existing token/theme sources and
the semantic bridge. The old sources still own their current values; the bridge
maps them to the shadcn-svelte contract rather than duplicating the palette.

| Whatnow source | Semantic token / utility |
| --- | --- |
| `--paper`, `--ink` | `--background`, `--foreground` / `bg-background`, `text-foreground` |
| `--card`, `--ink` | `bg-card`, `text-card-foreground` |
| `--primary`, `--on-primary` | `bg-primary`, `text-primary-foreground` |
| `--paper-deep`, `--ink-soft` | `bg-muted`, `text-muted-foreground` |
| `--chip`, `--ink` | `bg-accent`, `text-accent-foreground` |
| `--border-color`, `--stroke` | `border-border` and the shared outline width |
| `--shadow-sm`, `--radius-sm` | `shadow-control`, `rounded-md` |

Aliases are declared on both `:root` and `[data-theme]` so nested theme previews
resolve their own colors. Neon/Terminal variants use the existing theme attribute;
there is no second theme state to synchronize. Font follows the existing `--font`
token. The CLI-installed font remains available but does not replace it implicitly.

Legacy base rules live in Tailwind's `base` layer; `.btn`, `.field` and `.sticker`
live in `components`. Existing special theme treatments keep their specificity;
Plain's legacy anchor rules exclude the new `data-slot` primitives. New utilities
can style primitives without fighting the old global button/input reset.

## First migrated controls

- `ui/button` and `ui/input` were added by the official CLI, then customized in
  their source to share Whatnow tokens, outlines, radii and accessible focus.
- `PlayToggle` composes `Button` and configured Lucide icons. Normal lines and the
  large focus control share the same component, with the existing size variables.
- `InlineEdit` composes `Button` and `Input`: Enter/blur save, Escape cancels,
  pending state remains visible and failed writes retain the draft. No persistence
  logic moved into the UI primitives.

Keep the registry's bindable element references: `bind:this` in these primitives
is the public `ref` API, not incidental DOM manipulation to rewrite mechanically.

## Task and step menus

`DropdownMenu` replaces the hand-built task/step menus. `ContextMenu` opens the
same task actions from the card with a right click or a long press on its surface.
`TaskActionItems` is the single composition for both entry points: labels,
eligibility, move submenu and commands are shared. Inline editing fields keep
their native context menu; long presses on controls keep their own interaction.
The three-dot trigger remains available for touch and keyboard use.

The primitives own focus, keyboard navigation, outside dismissal and collision
handling through Bits UI. Menu surfaces use Whatnow's existing tokens, outline
and shadow, with 44px action rows. Inside the native task-editor dialog,
`BitsConfig.defaultPortalTo` scopes overlays to the dialog so they stay in its
top layer instead of appearing behind it.

The dropdown source was added through the
[neobrutalism-svelte registry](https://neobrutalism-svelte.flenze.com/docs/installation),
then adapted to the existing token bridge. Context menu and the other primitives
come from the configured official shadcn-svelte registry. There is no parallel
theme or additional application-state store.

## Shared task metadata

`TaskMetadata` is used in both Day and Focus. Minute estimates still use
`InlineEdit`; energy now composes `Select`, and feeling uses a compact `Popover`
with `Field`/`Textarea` and icon actions. Long feelings are truncated in the chip
and remain fully editable in the popover. Enter/blur saves the estimate, choosing
an energy saves it, and the check icon or Ctrl/Cmd+Enter saves free text.

The `metadata` variant on Button/Select.Trigger holds the shared treatment rather
than repeating per-surface control CSS. Select/Popover/Textarea use the same
existing semantic theme bridge. Bits UI manages the overlay and its keyboard/focus
behavior; there is no separate manual menu state or document-level listener.

Writes retain the existing full-field API, refetching the latest task first so an
edit does not replace unrelated fields with an old local snapshot. Energy's
selection follows the authoritative task value. A failed write keeps the feeling
draft and shows an inline error. A failed refresh after a confirmed write is
reported as a refresh failure, not an unsaved edit. Undo batch handling is intact;
no timer, step, priority or capture operations changed.

## Remaining work — parked on 2026-10-27

João considers the UI sufficient for now and asked to return to the overall
product direction. Keep these follow-ups for a near-term UI pass; do not start
another migration slice automatically. The current controls remain usable while
the product discussion resumes in [execution-direction.md](execution-direction.md).

| Order | Cohesive slice | Keep / intended result |
| --- | --- | --- |
| 1 | Task editor fields and dialog composition | Compose the existing editor with shared Input, Textarea, Field, Select and Button. Evaluate moving its native dialog to Bits UI only when it simplifies focus, overlays and maintenance. Preserve drafts across refreshes, explicit task save, immediate step saves, blockers, original capture and responsive layout. Do not create another step editor. |
| 2 | Day cards and shared execution rows | Polish density, hierarchy, long content, inline editing, reorder affordances and compact actual/estimate presentation as a set. Retain one ExecutionSteps interaction across Day, Focus and editor, and keep the main focus timer prominent/reachable on mobile. Known feedback is in feedback.md; João does not need to repeat it. |
| 3 | Remaining recurring controls and consistency | Consolidate remaining manual controls in capture/suggestion/chat/navigation when the shared primitives remove code or interaction friction. Check labels, icons, focus/error/pending states and existing themes together. Avoid converting every native element merely for library coverage. |

João's next screenshot clarified slice 2: Actions should emphasize the concrete
step over its parent task. Compress parent metadata into the row with the task
menu when space allows; energy can be an icon-only value control with an
accessible name. Keep full feeling/context reachable and allow mobile wrapping.
This is recorded in feedback.md under “Aba Ações”; it is refinement along the
way, not authorization to restart a broad visual migration now.

When resuming, use one primitive in every surface with the same intent. Add
missing components through `bunx --bun shadcn-svelte add <component>` from `web/`
and review their source. Consult https://shadcn-svelte-extras.com/ only for a
concrete missing interaction.

Separately, the architecture review identified estimate aggregation mirrored in
Go and TypeScript. Evaluate consuming the API's authoritative estimate consistently
when returning to that boundary; it is not a reason to interrupt product work or
change the agreed estimate meaning. This is a refactoring candidate, not a new
feature or a newly confirmed bug.

Keep task/step/timer operations, capture, chat, SSE and undo independent of this
visual migration. João continues visual acceptance through actual use; build and
static checks do not establish visual approval or a physical-device test.
