# Changelog

## 2026-10-27

### Continuidade após concluir um passo

- Corrigida a escolha no Foco após concluir a ação selecionada: abrir outro passo
  da mesma tarefa ou escolher execução direta agora respeita essa escolha.
  Antes, a seleção concluída forçava o primeiro pendente. Um timer em execução
  continua tendo prioridade; avanço ao próximo passo continua parado.
- Teste de regressão reproduziu a troca indevida antes da correção. Passaram
  174 testes frontend, Svelte check e testes Go de troca/pausa/retomada,
  idempotência, concorrência, conclusão e preservação de tempo/undo.
- João confirmou em uso: pedir uma ação no chat mostrou “Começar esta ação”,
  e clicar abriu o Foco e iniciou como esperado. Esse caminho está aceito.

### Mobile e ambiente de desenvolvimento

- Corrigido overflow horizontal no Dia: colunas da grade podem encolher abaixo
  da largura intrínseca do conteúdo. Removidas regras mobile duplicadas.
- Chat limita sua largura e quebra textos longos sem alargar a página.
- João confirmou a correção no navegador do celular, pelo acesso Tailscale.
- Vite permite o hostname do Mac na tailnet. HTTPS `:7001` encaminha para
  `localhost:5173`; `just dev` no tmux serve o frontend e a API atual em `7777`.
  O frontend usa HMR; alterações Go exigem reiniciar `just dev`.
- Confirmados HTTP 200 pelo HTTPS e campo `progress` na API de execução.
  O aceite do caminho chat → Começar esta ação → foco foi confirmado depois,
  conforme registro acima.
- Validação desta fatia: 173 testes frontend passaram, Svelte check sem erros
  ou avisos e build estático concluído. Autofixer sem issues; sugestões sobre
  efeitos preexistentes de navegação e scroll foram revisadas e preservadas.

O histórico anterior de implementação e validação está em
[docs/execution-direction.md](docs/execution-direction.md).
