# Feedback de uso — execução no Whatnow

Fonte: correções explícitas de João durante QA em 25/10/2026, delegadas pela conversa `01a17327-9bf8-79b2-a173-dd587d532f5b`. Registro contínuo do QA. A primeira rodada apenas capturou os achados; João autorizou implementação, pediu pausa para esboço e depois autorizou retomar com as correções visuais abaixo. Estados distinguem implementação, verificação técnica e aceite visual/funcional pelo uso. Preview e dados devem ser preservados.

## Aba Ações — refinamentos registrados em 27/10/2026

Fonte: [screenshot de João](feedback/2026-10-27/actions-view.png) e comentário
direto após estacionar a migração visual.
São ajustes para a próxima fatia de apresentação; não substituem a frente de
granularidade/continuidade da IA nem significam rejeição da direção atual.

- O passo deve ter mais ênfase que a tarefa na aba Ações: hoje título, notas,
  justificativa e metadados do pai ocupam a maior parte do cartão.
- Reunir metadados e três pontos em uma faixa compacta quando houver espaço,
  adaptando no celular sem forçar uma linha que transborde.
- Energia pode mostrar somente ícone/indicador do valor atual e abrir o seletor
  ao clicar. Manter nome acessível e dica para que o significado seja descoberto.
- Estimativa e sentimento também precisam de menos espaço; o texto completo e
  sua edição devem continuar alcançáveis. Reduzir o peso do contexto pai sem
  perder a relação com a tarefa ou criar outra interação para passos.

Retomada autorizada por João em27/10: implementado destaque do passo no modo
compacto compartilhado, contexto pai recolhido, metadados compactos e menu na
mesma faixa. Energia usa indicador com nome acessível; sentimento abre o mesmo
editor e preserva texto completo no nome/dica. Foco aplica os mesmos metadados
e hierarquia tarefa/ação. Ver CHANGELOG.md para validações. Aceite visual aberto.

## Migração de componentes — 27/10/2026

- Menu dos três pontos de tarefas e passos migrado para shadcn-svelte/Bits UI.
  A pedido de João, o cartão também abre as mesmas ações por menu de contexto;
  `TaskActionItems` mantém uma única composição para ambas as entradas.
  Check Svelte sem erros/avisos, 167 testes existentes e build passaram.
  Gesto touch e apresentação ficam para o QA de uso de João.
- **Correção:** as screenshots de João mostraram o controle do chat ausente.
  A classe legada `collapse` colidiu com a utility Tailwind `visibility: collapse`.
  Substituída por `Button` compartilhado com ícone X, nome acessível e área de
  toque de 44px; fechamento e preferência salva mantidos. João confirmou em
  seguida: “close button is ok now”. Isso aceita esse controle, não toda a UI.
- Metadados compartilhados Dia/Foco: estimativa mantém edição inline; energia
  usa Select com ícone/valor; sentimento é um chip editável por Popover com
  campo de texto livre. Erro de escrita mantém o rascunho; valor de energia
  acompanha o servidor. Composição sem o formulário avulso abaixo do cartão.
  Check Svelte sem erros/avisos, 167 testes existentes, build e detector
  determinístico Impeccable passaram. Aceite visual/interação desses controles
  permanece com João; nenhum teste em aparelho físico foi realizado.
- **Condução:** manter commits por fatia funcional/refatoração para revisão e
  rollback. Feedback pontual não substitui automaticamente a sequência de entrega.

## Rodada coesa U01–U10 — implementação autorizada

João respondeu **“beleza, vamos tentar, prossiga”** e autorizou implementar o conjunto proposto. Abaixo está o estado da nova rodada; o diagnóstico anterior permanece como origem dos requisitos, não como descrição atual de tudo que ainda falta. Sem nova entrevista, redesign ou torneio visual. Não há aceite automático de UX: implementar, compilar e testar lógica é diferente de aprovação pelo uso.

| Grupo | Estado atual e evidência | Aceite ainda aberto |
| --- | --- | --- |
| U01/U09 — decomposição contextual | Implementado `TaskDecompose` junto da lista, com **Dividir em passos com IA** para a tarefa correta. Novo `/tasks/{id}/decompose` usa o mesmo chat/SSE/histórico e só oferece/aceita `list_steps` e `add_step` desse ID; não pode alterar outra tarefa, apagar progresso, concluir ou trocar foco. **Adicionar passo** é manual; **Dividir este passo** substitui apenas o passo pendente com consequência explícita. Ferramentas com erro, provider indisponível e resposta sem novos passos não fingem conclusão; manual permanece. Testes agente/API com fixtures passaram. | Qualidade da decomposição contextual com provedor real não validada nesta rodada; a boa decomposição anterior do brain dump foi relato de João, não prova desta nova entrada. |
| U02/U03/U04 — mesma interação de passo | Implementado um único `ExecutionSteps.svelte` no Dia/Foco/editor; `StepsEditor.svelte` removido. Texto edita inline; estimativa à direita edita inline; checkbox conclui/reabre; play/pause por ícone; adição no próprio local; remover/dividir em menu secundário. Grip usa Pointer Events para mouse/toque e setas para alternativa de teclado, com anúncio de posição. Ordem salva somente após servidor confirmar. | Drag por toque/mouse e edição nos layouts reais não foram executados em navegador/aparelho nesta rodada. Não confundir teste da ordem e persistência com teste físico do gesto. |
| U07/U08 — execução, seleção e próximo | Implementado `focusTarget`: execução realmente rodando tem precedência sobre pedido de URL. Clicar texto edita e não troca alvo do painel. Ao concluir o alvo atual, apresenta o próximo pendente parado; concluir outra linha não muda o atual. Todos concluídos deixam a conclusão do pai explícita, sem iniciar próximo ou manter timer obsoleto. Testes de representação cobrem alvo ativo versus pedido de outra tarefa, próximo parado e conclusão alheia. | Uso visual/funcional de alternar/editar/concluir no celular permanece com João; os testes são lógicos. |
| U05/U06 — dado também é controle | Ícones SVG + valor para tamanho/energia no mesmo controle editável; sentimento continua livre. Tempo compartilhado compacto `05:00 / 15:00`; sem estimativa `05:00 / —`. Foco mantém timer principal/progresso e sticky; total pai soma registros uma vez, passo soma seus próprios segmentos. Alterar estimativa não modifica tempo real. | Legibilidade/densidade e descoberta dos controles dependem do QA de uso; não declarar visual aprovado por build. |
| U10/F05 — condução e preservação | Implementação tratada como uma rodada de componente/interação e estado, preservando os relatos e referências. Sem browser/review estético nem teste em dados reais. Backup consistente realizado antes da atualização;7777/configuração de IA preservados, sem reset de seed. | Preview atualizada em25/10 às15:44; hash/HTTP e preservação dos dados verificados abaixo. |

No Dia a lista é recolhível; mesmo sem passos ela oferece criação/decomposição no cartão. No editor é o mesmo componente de passos, não um formulário diferente para a entidade. Selecionar texto para editar não significa selecionar a ação que roda. O foco ativo pode manter prioridade sobre outra tarefa pedida na URL; para iniciar outra ação, usar seu play no Dia, preservando a regra de pausar a anterior.

### Verificação da rodada compartilhada

- Check Svelte: **0 erros/0 avisos**; **164 testes unitários frontend passaram**. Cobrem timer/restante/excedente, soma de tempo, execução ativa com seleção divergente, próximo apresentado parado, conclusão alheia sem trocar foco e cálculo de ordem por identidade/destino válido. MCP aprovou componentes alterados. Testes lógicos não comprovam gesto touch nem edição visual no aparelho.
- Testes Go do agente/API passaram, incluindo escopo de tarefa na decomposição, rejeição de ferramenta/ID alheio, preservação de passo concluído e manual após falha da IA. `go vet ./...` passou. Modelo real **não foi chamado para validar a nova entrada** nesta rodada.
- Smoke HTTP do executável final, em **banco temporário descartável**: criar/editar/reordenar passos, concluir sem iniciar timer, foco preservado ao editar descrição/estimativa, concluir pelo checkbox encerrando o registro sem concluir pai e falha de IA comunicada. Não usou nem modificou tarefas reais da preview.
- Preview7790 atualizada em **25/10/2026 às15:44 America/Sao_Paulo**. HTML local/Tailscale HTTP200 idêntico ao build final SHA256 `0f71fd56c7ef74ae76200275abbf47468a2b9b351641e4941ad40396e329864f`. Comparação somente de leitura com backup consistente preservou integralmente tarefas, passos, registros de foco, correções, tokens/seleção, settings e mensagens (snapshot:5 tarefas/13 passos/10 registros/6 mensagens). IA configurada preservada, integridade SQLite e referências válidas.7777 permanece no processo original.
- Nenhum browser/screenshot/review estético executado nesta rodada. **Implementado/verificado tecnicamente não significa aprovado por uso.** Arrastar por toque/mouse, legibilidade e descoberta das interações continuam no aceite de João. As referências e relatos abaixo ficam como histórico da origem, não itens apagados ou garantia de perfeição.

## Diagnóstico anterior à rodada coesa — devolutiva de João, 25/10/2026

**Escopo desta atualização: leitura do código e documentação apenas.** João relatou melhora, mas perguntou se todo o feedback estava implementado e apontou fricção ao tentar criar dois passos no próprio teste sugerido pelo assistente. Isso faz parte do fluxo central, não foi tangente. Não implementar, alterar preview/dados, rodar testes ou QA de navegador nesta rodada. O resumo anterior de entrega e testes verdes não comprovam atendimento de todos os critérios de UX.

**Entregue com evidência anterior, preservando seus limites:** duração independente removida de schema/API/cliente, restante/excedente por estimativa e tempo acumulado, registros/correções preservados na migração; checkboxes sem iniciar timer, play/pause por ícone, lista de execução compartilhada entre Dia/Foco, metadados editáveis junto ao título, acesso compacto ao foco e Ver dia sem pausa. Sincronização computador/celular e decomposição útil do brain dump foram confirmadas por João. O editor completo, a edição inline de passos, a reordenação por drag, a decomposição contextual de tarefa e os estados seleção/execução não estão unificados ou completos.

**Conclusão honesta: não, nem todo o feedback anterior foi atendido.** F02–F04 e W01–W03/W05–W06 têm atendimento parcial; a tabela histórica abaixo foi corrigida para não confundir implementação de trechos com resolução do fluxo. F01 permanece tecnicamente implementado/verificado; o problema U08 é qual alvo o painel representa, não prova de perda do tempo no servidor. W04 (lápis) está implementado; aceite visual integral não foi estabelecido.

## Achados preservados — estado antes da rodada coesa

Fonte de todos os relatos U01–U10: devolutiva explícita de João nesta conversa, delegada por `01a17327-9bf8-79b2-a173-dd587d532f5b` em 25/10/2026. “Falta” abaixo significa ausência no caminho inspecionado, não auditoria do produto inteiro. Evidência é leitura estática desta versão; sem reprodução por navegador nesta rodada.

### Criar, decompor e editar a mesma entidade

| ID | Origem, estado e evidência | Resultado esperado / limite |
| --- | --- | --- |
| U01 | **Parcial — decomposição da tarefa disponível pelo chat, sem affordance contextual.** João não encontrou como pedir para quebrar a TAREFA. `internal/agent/instructions.go:18` orienta `add_step` sob demanda; `internal/agent/steps.go:18` oferece essa ferramenta. `split_step` exige um passo existente (`:22`). Em Foco, “Ainda está grande: adicionar passos” abre o formulário manual quando não há passos; não solicita decomposição da tarefa à IA. | Tornar a decomposição de tarefa central e encontrável, distinta de adicionar manualmente um passo e de dividir um passo existente. Não confundir capacidade do chat com entrada visível no contexto da tarefa. |
| U02 | **Parcial — criar inline existe, editar continua distante.** João percorreu tarefa → menu → editar → scroll → lista textual; também encontrou Adicionar passo inline e outra edição com uma lista duplicada. `ExecutionSteps.svelte:34` adiciona só descrição; `TaskEditor.svelte:155` monta `StepsEditor`; Foco `:92–93` monta outro editor abaixo da lista já exibida. | Uma superfície coerente para criar/editar os mesmos passos, sem exigir essa navegação nem abrir outro formulário/lista para a mesma operação. Há um atalho manual novo, mas o problema anterior F02–F04 não está resolvido. |
| U03 | **Parcial — compartilhado apenas entre Dia e Foco.** Ambos usam `ExecutionSteps`, com checkbox e play/pause. O editor usa `StepsEditor`, que oferece Editar/Subir/Descer/Remover e outros formulários. Texto em `ExecutionSteps.svelte:29` seleciona/navega; não edita inline. A estimativa fica em `small`, não é um campo clicável à direita. | Padrão compartilhado Dia/Foco/editor: checkbox, texto clicável para edição inline, estimativa à direita clicável, play/pause por ícone. Eliminar listas/formulários duplicados e ações textuais redundantes. Compartilhar dados ou duas renderizações não satisfaz sozinho esse critério. |
| U04 | **Falta — drag para reordenar, pedido explícito nesta devolutiva.** API `reorder_steps`/REST existe; UI inspecionada usa `StepsEditor.svelte:40–43,71–73` com Subir/Descer. Não há drag nessa superfície. | Arrastar para reordenar e manter alternativa acessível, sem botões textuais dominantes. O pedido de drag é explícito agora; não atribuí-lo retroativamente como recurso que já havia sido entregue. |
| U05 | **Parcial — metadados editam, representação visual falta.** `TaskMetadata.svelte` compartilha dropdowns nativos de tamanho/energia e sentimento livre entre Dia/Foco. Não há ícones de energia/tamanho; opções mostram S/M/L e Baixa/Média/Alta. | Tratar tamanho/energia visualmente com ícones/valor, sem converter sentimento livre em enum. Origem é o feedback visual reiterado de João; W01 não estava completo. |

### Tempo, continuidade e identidade do foco

| ID | Origem, estado e evidência | Resultado esperado / limite |
| --- | --- | --- |
| U06 | **Parcial — tempo existe, apresentação do Dia segue pequena/verbosa.** João descreveu “3 minutos e 10 segundos dedicado” e “sem estimativa”. Nesta fonte, `TaskCard.svelte:136,168` usa `timerText` (MM:SS), sufixos “dedicado / N min previstos” ou “sem estimativa”, em 13px; `ExecutionSteps.svelte:29,42` usa a mesma prosa em 12px. A fala do usuário é mantida como percepção, sem afirmar que o código formata literalmente segundos por extenso. | Mostrar compacto e legível, por exemplo **05:00 / 15:00**, sem sobrecarregar labels. Sem estimativa João sugeriu infinito OU traço/ausência; são alternativas, não aprovação de infinito. **Recomendação do assistente: traço**, porque infinito não equivale a desconhecido. |
| U07 | **Parcial — próximo é encontrado/oferecido, mas não selecionado automaticamente.** Ao concluir, `ExecutionSteps.svelte:14–17` só atualiza e refaz consulta. Foco `:27–30` mantém `step` explícito na URL e `:84` exige clicar “Ir para”. Sem `step` explícito, o fallback ao primeiro pendente pode mudar, mas isso não garante o fluxo inteiro. | Após concluir o atual, selecionar/apresentar o próximo automaticamente **parado**, sem iniciar timer nem concluir o pai. Não apresentar a oferta atual como atendimento integral desta regra. |
| U08 | **Bug de representação confirmado por leitura; execução no servidor não é o mesmo que seleção da tela.** João relata: segundo passo rodando → clicar primeiro concluído → painel volta ao antigo e execução real fica só no indicador pequeno. Foco `:23–30` deriva o alvo da URL antes do foco persistido; `choose` (`:47`) muda a URL sem pausar. `clock/status` (`:32–46`) seguem o alvo visual; `ExecutionSteps.svelte:26–29` conserva apenas o pequeno “Em foco” e play de pausa na linha que roda. `.selected` destaca a linha selecionada (`:27,39`), não necessariamente a rodando. | Foco em execução permanece inequívoco e com timer correspondente. Separar estados “rodando”, “selecionado para ver/editar” e “próximo pronto” em vez de deixar a seleção substituir silenciosamente o painel de execução. Preservar navegação sem pausa e a sincronização aprovada. O critério precisa ser tratado no conjunto do fluxo, não só mudar a cor de um número. |
| U09 | **Parcial — decomposição manual existe, consequência do rótulo é obscura.** “Ainda está grande” assustou João, que não sabia o que aconteceria. No editor abre texto para substituir um passo pendente por pares; em Foco sem passos abre adição manual. Mesmo vocabulário representa consequências diferentes. | Decomposição contextual com linguagem que deixe claro tarefa versus passo e ação que ocorrerá. A preocupação não cria uma regra universal de confirmação nem autoriza popups em todas as mudanças. |
| U10 | **Pendência de condução, não uma nova funcionalidade.** Feedback apareceu ao tentar criar só dois passos para o teste sugerido; João pede visão geral e questiona sequência de remendos. | Inventário honesto e próxima rodada coerente de componente/interação compartilhada e estados do foco. Não exigir repetição dos achados, chamar o relato de tangente ou iniciar correções cosméticas nesta rodada. |

## Rodada proposta e depois autorizada

1. **Unificar a interação de passos** no Dia/Foco/editor: criar, editar texto e estimativa na linha, concluir/reabrir, play/pause e reordenar, com alternativa acessível ao drag. Remover a segunda lista de edição. Incluir decomposição da tarefa e do passo com contexto/consequência claros.
2. **Resolver a continuidade como um conjunto:** concluir atual apresenta próximo parado; ver/editar outro passo não troca silenciosamente o timer de quem está rodando; seleção e execução ficam distintas com destaque inequívoco.
3. **Aplicar a apresentação compartilhada** de tempos e metadados: tempos compactos comparáveis; ausências sem orçamento inventado; ícones para tamanho/energia e sentimento livre. Reaproveitar a organização já documentada de Sunsama/Passos/Shift, preservando Whatnow, banco e sincronização.

Na rodada de registro isto era uma proposta, sem autorização de executar. João autorizou posteriormente esta unidade de trabalho; ver o estado atual acima. Antes da próxima implementação, levar em conta que o novo relato atualiza o contrato anterior de “oferecer próximo” para “selecionar próximo sem iniciar”. Não é necessário redesenhar stack, infraestrutura, calendário ou produto. O QA visual/funcional permanece com João; a implementação e os testes novos constam acima, sem transformar testes antigos/build em aceite da experiência.

## Correção de origem e prioridade

A duração independente de sessão foi introduzida pelo assistente; João esclareceu que não a pediu e a rejeitou. Sua presença no plano e na implementação anteriores não a transforma em requisito dele. Este feedback prevalece sobre descrições conflitantes em `PRODUCT.md`, `.impeccable/focus.brief.md`, `docs/execution.md` e na nota Basic Memory «Whatnow — passos leves e foco com tempo corrigível». `PRODUCT.md`, `.impeccable/focus.brief.md` e `docs/execution.md` foram reconciliados nesta implementação. A nota histórica Basic Memory descrevia o modelo anterior e não deve ser usada para restaurar a duração rejeitada.

## Backlog confirmado

| ID | Estado | Achado de João e resultado necessário |
| --- | --- | --- |
| F01 | Implementado · lógica, migração e preview verificadas | **Remover a duração independente e o default de 15 minutos do comportamento, API e UI**, não só esconder o campo. A ação tem estimativa e tempo efetivo acumulado: estimativa de 2min menos 1m50 dedicados = 10s restantes; depois de zero, mostrar excedente explícito sem concluir automaticamente. Pausa, troca e retomada preservam tempo por ação; correção explícita permanece. Segmentos internos podem continuar servindo à persistência, sem virarem outro orçamento de tempo para o usuário. |
| F02 | Parcial · contexto entregue; edição de passos e distinção seleção/foco pendentes (U01–U04, U08) | **Manter contexto e edição da tarefa no foco.** O pai está apagado/secundário demais, seus dados não ficam acessíveis e a edição não é evidente. Deve ser claro que a ação pertence à tarefa, com informações e edição ao alcance, sem procurar em outra tela. João não afirmou que destacar o título do passo é errado. |
| F03 | Parcial · lista contínua entregue; interação única e avanço sem play faltam (U03–U04, U07) | **Mostrar tarefa e seus passos como um conjunto contínuo.** «0 de 2 passos concluídos» isolado não substitui a lista. A referência Sunsama mostra tarefa como cabeçalho/contexto, passos em lista, realizado/planejado comparáveis e detalhes abaixo. É referência de organização e continuidade; não autoriza copiar sidebar, calendário, Pomodoro, estética ou todo o produto. |
| F04 | Parcial · menos controles grandes, mas duplicação e fricção persistem (U03–U06, U09) | **Reduzir escala, dispersão e conceitos que atrapalham executar.** João considera a UI atual grande, dispersa e obscura. Recuperar os aprendizados já existentes em Passos Svelte e Shift antes de propor mudança futura, preservando a identidade Whatnow. Não reabrir exploração ampla nem inventar um novo sistema de execução. |
| F05 | Em acompanhamento contínuo | **Conduzir QA com orientação curta e guardar achados aqui.** João faz o QA visual/funcional pelo uso real; não pedir que repita o feedback nem criar nova matriz extensa. Novos relatos entram com fonte/data e estado, sem serem declarados resolvidos antes da mudança e verificação correspondente. |

## Refinamentos após o esboço — fonte explícita de João

A imagem `/Users/joao/.codex/generated_images/01a17327-9bf8-79b2-a173-dd587d532f5b/exec-f3b11389-db8b-4106-9756-f5c085a1e335.png` recebeu correções; não foi aprovada integralmente. João autorizou retomar sem outra aprovação de imagem. “Samsung” nas transcrições refere-se a Sunsama, confirmado por ele.

| ID | Correção confirmada | Estado nesta entrega |
| --- | --- | --- |
| W01 | Energia, tamanho e sentimento compactos abaixo do título; o valor também edita. Dropdown para enums, texto livre para sentimento; mesmo padrão Dia/Foco. | Parcial: edição compartilhada entregue; energia/tamanho ainda só em labels, tratamento visual/ícones falta (U05). |
| W02 | Um play/pause por ícone com nome acessível e alvo touch. Checkboxes concluem tarefa/passos, sem grandes botões duplicados. | Entregue nas listas Dia/Foco; editor não usa o mesmo padrão, unificação parcial (U03–U04). |
| W03 | Tempo é informação principal em Foco: número grande e indicador visual/progresso em desktop/mobile. O erro anterior era tempo errado e contexto separado, não simplesmente timer grande. | Parcial: tamanho/progresso implementados; painel pode representar seleção antiga em vez da execução atual (U08). |
| W04 | Edição por lápis discreto/reposicionado, sem texto “Editar” comprimindo título mobile. | Implementado; aceite de uso pendente. |
| W05 | Ação atual, timer e play/pause permanecem visíveis em região sticky ao rolar; lista e controles continuam acessíveis. | Parcial: sticky implementado; ação efetivamente rodando perde destaque ao selecionar outra linha (U08). |
| W06 | Dia sem rodapé gigante para foco: entrada compacta reconhecível; adicionar passo também compacto. | Entrada compacta/add inline entregues; tempo verboso e edição/criação duplicadas permanecem (U03, U06). |

Preservar a sincronização/controles que João confirmou funcionar no computador e celular e a IA que ele relatou decompor bem o brain dump real. Não reimplementar infraestrutura ou linha do tempo. A tarefa principal totaliza cada registro de trabalho uma vez (direto + passos), e o passo soma seus segmentos exatos. Metadados/estimativa não resetam foco. Migração 0011 remove apenas a duração independente; backup consistente precede a aplicação ao banco em uso.

## Evidência técnica e limites

- Testes Go existentes e focados passaram: migração de registro rodando com seleção/token/correção preservados, conclusão de passo sem iniciar foco, conclusão pelo checkbox encerrando o registro ativo sem perder tempo, soma ao longo de segmentos, correção e undo preservando tempo. API aceita start sem duração e não devolve o campo rejeitado. `go vet ./...` passou.
- Helpers frontend cobrem 2min−1m50=10s, excedente explícito, ausência de estimativa e soma por ação/pai sem duplicar registros. Check Svelte sem erros/avisos, 159 testes unitários frontend e build passaram nos refinamentos; nova rodada visual/browser não executada.
- Nenhuma nova rodada de screenshot, review estético ou browser foi realizada. A tentativa de um teste curto na rodada anterior foi recusada pela revisão automática por causa da orientação anterior de João; não houve contorno ou insistência. Testes de navegador antigos não validam a UI refinada atual.
- Não marcar UX como resolvida apenas por compilação; aprovação funcional/visual permanece com João. Preview atualizada em 25/10/2026 às 14:31 America/Sao_Paulo na porta 7790. HTTP local/Tailscale serviu exatamente o HTML do build final (`14dc1ff4889e497b6a619a64872f73a88bb5a220dd363a2393e650baf29de2f8`). Comparação somente de leitura com backup consistente confirmou todos os campos remanescentes das 4 tarefas, 11 passos, 4 registros de foco, tokens/seleção, correções e settings idênticos. Integridade SQLite ok, referências válidas, IA configurada preservada; duração removida do schema/API. Serviço7777 segue no mesmo processo. Estes números são um snapshot técnico; João pode continuar criando/alterando dados.

## Referências recuperadas — leitura direcionada

Foram localizados `/Users/joao/space/passos`, `/Users/joao/space/shift` e `/Users/joao/space/shift_1`. Os ponteiros de execução abaixo são do Passos Svelte e de `shift`; `shift_1` também existe, mas não foi tomado como a implementação de referência. Leitura de fonte não comprova que esses fluxos estejam rodando ou aceitos em seu estado atual.

- **Passos Svelte — restante, excedente e foco compacto:** [FocusPanel.svelte](/Users/joao/space/passos/src/lib/components/FocusPanel.svelte:35) calcula estimativa menos tempo decorrido e exibe `+` no excedente; [CompactFocus.svelte](/Users/joao/space/passos/src/lib/components/CompactFocus.svelte:13) usa a mesma relação. [Agora](/Users/joao/space/passos/src/routes/agora/+page.svelte:668) mantém edição da estimativa junto à ação e mostra restante/extra na lista. Reaproveitar o aprendizado de interação, sem transplantar stack ou presumir equivalência de persistência.
- **Shift — contexto e continuidade:** [focus.tsx](/Users/joao/space/shift/src/components/passos/focus.tsx:32) relaciona tarefa, filho selecionado, tempo planejado/real e excedente; a mesma superfície mantém a [lista de subtarefas](/Users/joao/space/shift/src/components/passos/focus.tsx:47) e controles por ação. É um ponto concreto para estudar continuidade, não uma afirmação de que já satisfaz toda a edição pedida em F02.
- **Shift — tempo acumulado e correção:** [passos-model.ts](/Users/joao/space/shift/src/lib/passos-model.ts:71) soma ajustes, segmentos encerrados e segmento ativo por alvo; [time-control.tsx](/Users/joao/space/shift/src/components/passos/time-control.tsx:93) separa editar realizado de editar estimativa. Regras específicas de agregação/entrada desse projeto não foram aprovadas automaticamente para Whatnow.
- **Basic Memory pessoal:** [Passos — direção do produto e guia de implementação](</Users/joao/space/memory/decisoes/Passos — direção do produto e guia de implementação.md:71>), permalink `codex/decisions/passos-direcao-do-produto-e-guia-de-implementacao`, seção «Execução e continuidade»: sequência curta de ações, estimativa para perceber tempo e foco como **visualização dos mesmos dados**, com edição coerente. A nota distingue relato, hipótese e aprovação; não importar seu catálogo de funcionalidades. A consulta a «Shift» trouxe resultados sem correspondência pertinente; a referência específica foi recuperada diretamente do código local.

## Evidência visual fornecida por João

Cópias preservadas neste projeto, sem edição ou nova captura:

- [Whatnow — foco considerado confuso](feedback/2026-10-25/whatnow-focus.png): estimativa de 2min coexistindo com sessão de 15min, grande ocupação pelo passo/timer e contexto pai secundário.
- [Sunsama — tarefa com passos](feedback/2026-10-25/sunsama-task.png): tarefa e lista com tempos no mesmo conjunto.
- [Sunsama — foco e detalhes](feedback/2026-10-25/sunsama-focus.png): tarefa, passos, realizado/planejado e detalhes em continuidade. Outros elementos visíveis não ampliam o pedido.

## Próximo recorte de uso

Os relatos U01–U10 estão preservados, sem pedir que João os repita. Com a nova preview, o recorte é criar dois passos no próprio cartão, editar texto/estimativa na linha, marcar um já feito e iniciar o seguinte. Editar outra linha não deve tirar o timer do foco ativo. O QA visual/funcional permanece com João; esta rodada não executa browser nem altera dados de trabalho para testá-los.
