# Whatnow development

## Frontend components

- Before frontend UI work, read `.agents/skills/shadcn-svelte/SKILL.md`. This is a relative symlink to the canonical skill installed in `web/.agents/skills/shadcn-svelte/`; keep one copy.
- Use **shadcn-svelte**, with the configuration in `web/components.json`, for shared UI primitives. Read that skill and the official component documentation before adding or composing components. Do not use the React shadcn CLI or React component APIs in this Svelte app.
- Run the CLI from `web/` with Bun: `bunx --bun shadcn-svelte add <component>`. Check `web/src/lib/components/ui/` first; review the files and aliases produced by the CLI.
- Consult [shadcn-svelte-extras](https://shadcn-svelte-extras.com/) when the official components do not cover a concrete interaction. Inspect its installation instructions, dependencies and compatibility before adding an item; it is a reference, not an additional mandatory dependency.
- Migrate incrementally: semantic theme tokens first, shared controls next, then the surfaces that use them. Preserve Whatnow's visual identity and theme choices. See `docs/ui-migration.md` for the compatibility bridge.
- Product components own task/step/focus behavior; UI primitives own common presentation and interaction. Keep server state, SSE, execution history and undo intact during visual work.

`web/AGENTS.md` contains the frontend tooling and Svelte validation instructions.
