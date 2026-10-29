# Changelog

## 2026-10-29

- Checkbox de tarefa integra a grade do cabeçalho e alinha ao título, em vez
  de ficar presa ao topo do cartão. No celular os metadados ocupam a linha
  seguinte; título e checkbox continuam alinhados.

- Estimativas numéricas permanecem em uma linha (`28+ min`). Energia usa número
  e um/dois/três raios Lucide da mesma cor para baixa/média/alta, inclusive no
  seletor. Cartão de Prioridades reúne ação, tempo e controles na mesma linha
  quando cabe; total da tarefa acompanha a contagem de passos à direita.

- Faixa persistente fora da tela Foco mostra a ação em andamento, tempo dedicado,
  pausa e link de retorno. Usa o estado sincronizado e as operações existentes;
  navegar de volta não inicia outro registro de tempo.
- A faixa destaca tempo restante/excedente, mostra dedicado / estimativa da
  ação e transforma o separador em progresso. Sem estimativa, mostra dedicado
  e total desconhecido, sem barra percentual fictícia.
- Refinamento pedido após uso: tempos na mesma linha (`restante | dedicado /
  estimativa`), restante maior, sem rótulo visual redundante; descrição acessível
  preservada e excedente continua com sinal `+`.
- Cartões compartilham metadados compactos em Prioridades e Ações. Título e
  controles aproveitam a mesma linha quando cabe; energia abre seleção pelo
  indicador, estimativa edita ao clicar e contexto expande explicitamente.
  Removida expansão por hover e reduzido espaçamento da área de execução.
- João confirmou instalação PWA funcionando; faixa azul em standalone permanece
  como pequeno bug de viewport sem causa confirmada.
- Validação: 181 testes unitários, Svelte check sem erros/avisos, build e
  autofixer dos componentes alterados passaram. Teste Chromium com banco
  temporário verificou navegação/retorno, faixa sticky, pausa sem nova sessão,
  edição de energia/estimativa e ausência de overflow em320/390/1440px.
  Capturas inspecionadas; aceite no aparelho físico permanece com João.

## 2026-10-28

- Adicionado manifesto PWA vinculado no HTML, com nome Whatnow, abertura
  standalone, escopo/entrada na raiz e ícones PNG192/512 derivados do favicon.
  Sem dependências ou cache de dados. Build passou; manifesto no HTML gerado
  e dimensões/caminhos dos ícones verificados. Instalação no Android fica para
  teste de João após `just install`.

## 2026-10-27

### Ações e Foco: ação primeiro, contexto compacto

- Ações destaca a descrição do passo, com contexto da tarefa em tipografia
  secundária. Notas e justificativa ficam sob “Ver contexto”, sem expansão por
  hover nessa perspectiva. Prioridades preserva a apresentação anterior.
- TaskMetadata compartilha modo compacto entre Ações/Foco: energia por indicador,
  sentimento por ícone com texto completo acessível no controle e popover,
  estimativa editável e data preservadas. Menu da tarefa ocupa a mesma faixa.
- A linha compacta usa os mesmos controles de ExecutionSteps; no celular,
  descrição e tempo/ações se distribuem em duas linhas. Foco reduz o peso do
  título pai e destaca a ação junto ao timer. Nenhuma regra de execução mudou.
- 181 testes, check Svelte, autofixer e build passaram; detector sem achados.
  Inspeção limitada em Chromium: Ações sem overflow em320/390/1440px, Foco
  em320/390px e captura desktop. Não houve alteração de tarefas/timers durante
  a inspeção. Aceite visual em aparelho continua com João.

### Recuperação do chat sem perder respostas recentes

- Uma leitura atrasada do histórico não remove mais mensagens confirmadas que
  chegaram enquanto ela estava em andamento, incluindo referências de ação.
- Novos deltas de uma resposta em andamento não são apagados pela recuperação.
  Leituras anteriores não substituem uma recuperação mais recente.
- Três testes reproduziram as falhas antes da correção. A lógica permanece no
  estado do chat; não altera API, banco ou regras de execução.
- Resposta duplicada por SSE/HTTP não encerra um rascunho posterior nem repete
  o indicador de não lida. Teste reproduziu a falha antes da correção.
  181 testes frontend, check Svelte e build passaram. Autofixer sem issues;
  Maps locais de recuperação são intencionalmente não reativos.

### Entrada de foco compartilhada e erros visíveis

- Cartões, passos, chat e E agora compartilham `focus-navigation.ts` para iniciar
  pela operação existente e navegar somente após confirmação. Falha de navegação
  informa que o foco iniciou, sem repetir ou desfazer a operação.
- Falhas de play/pausa nos cartões e passos agora aparecem onde ocorreu o clique;
  o último estado confirmado permanece disponível.
- Mantidas as fronteiras: Go possui regras/tempo/dados, estado frontend possui
  sincronização, componentes apresentam resultados. Documentado em AGENTS.md
  para preservar uso pelo terminal e futuras integrações de contexto.
- 177 testes frontend e Svelte check passaram; autofixer sem pendências nos
  quatro componentes. Testes novos verificam falha sem navegação, confirmação
  antes de navegar e falha de navegação sem duplicar o início.

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
