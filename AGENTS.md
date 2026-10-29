# Whatnow development

## Layer boundaries

- Keep execution rules, eligibility, time, persistence and undo in the Go service.
  The web UI, terminal client and future context integrations must be able to use
  the same operations and data without depending on Svelte or browser navigation.
- Frontend state owns synchronization and confirmed server snapshots. Components
  own interaction and presentation; shared web navigation belongs in a small web
  helper, not in domain/service code. Do not duplicate execution rules per screen.
- Prefer cohesive changes with relevant contract tests and a project changelog.
  Consolidate demonstrated duplication; do not add speculative abstractions.

## Codebase navigation

Use current source files and checked-in project docs as the primary evidence.
Graphify was an experiment, not part of the required development workflow.
Do not query or rebuild its graph unless João specifically asks to use it.

## Frontend components

- Before frontend UI work, read `.agents/skills/shadcn-svelte/SKILL.md`. This is a relative symlink to the canonical skill installed in `web/.agents/skills/shadcn-svelte/`; keep one copy.
- Use **shadcn-svelte**, with the configuration in `web/components.json`, for shared UI primitives. Read that skill and the official component documentation before adding or composing components. Do not use the React shadcn CLI or React component APIs in this Svelte app.
- Run the CLI from `web/` with Bun: `bunx --bun shadcn-svelte add <component>`. Check `web/src/lib/components/ui/` first; review the files and aliases produced by the CLI.
- Consult [shadcn-svelte-extras](https://shadcn-svelte-extras.com/) when the official components do not cover a concrete interaction. Inspect its installation instructions, dependencies and compatibility before adding an item; it is a reference, not an additional mandatory dependency.
- Migrate incrementally: semantic theme tokens first, shared controls next, then the surfaces that use them. Preserve Whatnow's visual identity and theme choices. See `docs/ui-migration.md` for the compatibility bridge.
- Product components own task/step/focus behavior; UI primitives own common presentation and interaction. Keep server state, SSE, execution history and undo intact during visual work.

`web/AGENTS.md` contains the frontend tooling and Svelte validation instructions.
