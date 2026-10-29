---
type: "query"
date: "2026-10-25T01:40:47.710016+00:00"
question: "como o backend se conecta com o front-end"
contributor: "graphify"
outcome: "useful"
source_nodes: ["App", "Service", "Hub", "httpapi: REST and SSE under /api/v1"]
---

# Q: como o backend se conecta com o front-end

## Answer

Consulta de arquitetura orientada pelo grafo, com conferência pontual no código. Vocabulário: httpapi, service, sse, refetch, sqlite, embedded, terminal, facts, undo, fallback, prompt.
O projeto tem um núcleo Go compartilhado por uma interface web SvelteKit e uma interface terminal Bubble Tea. A web chama endpoints REST JSON em /api/v1; api.ts centraliza fetch e tipos que espelham o JSON Go. O servidor publica eventos do Hub por SSE em /api/v1/stream. O design documenta refetch do estado autoritativo após eventos de alteração e SSE também transporta partes do chat.
O frontend é compilado em arquivos estáticos e embutido no executável Go com go:embed. Node é usado na construção do frontend, sem exigir servidor Node para servir o pacote compilado. O terminal chama os serviços Go diretamente; app.Open constrói os componentes para ambas as interfaces. Processos separados não compartilham o Hub em memória.
As decisões documentadas priorizam uso pessoal local, SQLite, operações de domínio compartilhadas com batches e undo, código para fatos e restrições, modelos para interpretação e planejamento, captura persistida antes da IA e fallback em falhas de planejamento.
Interpretação: o conjunto funcional se beneficia de responsabilidades delimitadas e de mecanismos reutilizáveis. Isso não comprova ausência de bugs nem que a stack sozinha explique a qualidade do produto. Para adaptar o fork, preservar o núcleo e introduzir preferências no planejamento e nas interfaces é um ponto de partida coerente.
Fontes: internal/app/app.go:56; web/src/lib/api.ts:291; internal/httpapi/server.go:61; internal/httpapi/stream.go:10; internal/web/web.go:12; docs/superpowers/specs/2026-10-05-whatnow-design.md; docs/superpowers/specs/2026-10-05-whatnow-web-ui-today-design.md; docs/superpowers/specs/2026-10-06-whatnow-tui-design.md.

## Outcome

- Signal: useful

## Source Nodes

- App
- Service
- Hub
- httpapi: REST and SSE under /api/v1