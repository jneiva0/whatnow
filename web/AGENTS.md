You are able to use the Svelte MCP server, where you have access to comprehensive Svelte 5 and SvelteKit documentation. Here's how to use the available tools effectively:

Use Bun for frontend dependencies and scripts. Install with `bun install --frozen-lockfile`,
run scripts with `bun run <script>`, and invoke package CLIs with `bunx`.
Run the Vitest suite with `bun run test`; `bun test` invokes a different test runner.
Keep `bun.lock` as the frontend lockfile and respect the Bun version in `package.json`.

## UI system

Read `.agents/skills/shadcn-svelte/SKILL.md` in this frontend directory (also exposed
at the repository root). Use shadcn-svelte primitives and its CLI according to
`components.json`; consult https://shadcn-svelte-extras.com/ when the official
collection lacks a needed interaction. Preserve the Whatnow themes through the
semantic token bridge in `src/routes/layout.css`. Follow `../docs/ui-migration.md`.
Do not recreate shared buttons, menus, dialogs or form controls in individual
product components when an installed primitive can provide the interaction.

## Available Svelte MCP Tools:

### 1. list-sections

Use this FIRST to discover all available documentation sections. Returns a structured list with titles, use_cases, and paths.
When asked about Svelte or SvelteKit topics, ALWAYS use this tool at the start of the chat to find relevant sections.

### 2. get-documentation

Retrieves full documentation content for specific sections. Accepts single or multiple sections.
After calling the list-sections tool, you MUST analyze the returned documentation sections (especially the use_cases field) and then use the get-documentation tool to fetch ALL documentation sections that are relevant for the user's task.

### 3. svelte-autofixer

Analyzes Svelte code and returns issues and suggestions.
You MUST use this tool whenever writing Svelte code before sending it to the user. Keep calling it until no issues or suggestions are returned.

### 4. playground-link

Generates a Svelte Playground link with the provided code.
After completing the code, ask the user if they want a playground link. Only call this tool after user confirmation and NEVER if code was written to files in their project.
