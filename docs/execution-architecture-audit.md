# Auditoria de arquitetura — planejamento e execução

27/10/2026 · código inspecionado no commit `147777b`, branch
`codex/execution-direction`. Auditoria solicitada por João após o chat responder
verbalmente a uma sequência que ele esperava ver organizada em passos.

## Conclusão

A estrutura de persistência e execução é aproveitável. Go, SQLite, serviço único,
atividade/desfazer, API/SSE e timer autoritativo continuam formando uma base
coerente. Não foi encontrada evidência que justifique reescrever a stack ou
reinventar tarefas e passos.

**A integração comportamental está incompleta:** os passos foram incorporados ao
banco, às ferramentas e à UI com mais profundidade do que aos fatos derivados,
ao contrato de organização e à continuidade entre superfícies. Acrescentar
instruções aos prompts e campos às ferramentas habilita um comportamento, mas
não define um fluxo verificável de organizar → salvar → apresentar → executar.

O princípio original, em `docs/superpowers/specs/2026-10-05-whatnow-design.md`, é
útil para corrigir isso: código estabelece fatos e valida operações; IA faz os
julgamentos. Parte da nova lógica deixou cálculos e a própria escolha entre
organizar e apenas aconselhar implícitos no modelo.

Esta auditoria não mede eficácia clínica nem demonstra que a versão original
era estatisticamente melhor. Identifica diferenças de contrato e falhas concretas
na extensão atual. Não houve correção de código nesta rodada.

## Fluxo atual e limites de cada módulo

| Área | O que existe de fato | Alinhamento |
| --- | --- | --- |
| `domain` / `store` | Task, passos planos, estimativa explícita, foco persistente, correções independentes e histórico de atividades. | Suficiente para executar o fluxo atual; insuficiente para prometer calibração histórica automática. |
| `service` | Validação, CRUD/reordenação/divisão, tempo, elegibilidade, idempotência e undo compartilhados pela API e IA. | Base forte a preservar. Valida forma/estado, não a adequação semântica de uma decomposição. |
| `cleanup` | Entrada salva antes da IA; saída estruturada de tarefas/eventos/passos; associação a intenção existente, preservação da captura, lote reversível. | Sabe acrescentar passos, mas não reconciliar/reordenar uma sequência existente. |
| `facts` / `prompt` | Disponibilidade, bloqueios, datas, intervalos livres e estimativa total; contexto textual com notas/captura. Execução acrescentada como JSON bruto pelo planner. | Progresso, tempo agregado por ação e próximas ações ainda não são fatos de primeira classe. |
| `planner` | Plano de prioridades versionado e escolha estruturada de uma tarefa/passo elegível, retry/fallback. | “Escolher entre ações existentes” está definido; “preparar uma ação ainda vaga” fica fora desse contrato. |
| `agent` | Histórico do chat e loop livre de ferramentas; uma resposta pode consultar, modificar ou apenas conversar. | Não há resultado tipado que distinga organização aplicada, conversa, sugestão e impossibilidade/ambiguidade. |
| `httpapi` / SSE / `web` | Mutações passam pelo serviço, execução atualiza clientes; Dia/Foco usam componentes compartilhados. | Sincronização integra os dados; isso não torna chat, E agora e Ações uma mesma decisão de execução. |
| Encerramento / histórico | Resumos de tarefas feitas, não feitas e triagem; execução preservada separadamente. | O resumo do dia ainda não incorpora passos concluídos ou tempo dedicado. Não há ciclo de aprendizado implementado. |

## Cinco achados principais

### 1. Organizar uma sequência ainda é uma decisão implícita do chat

`internal/agent/agent.go:173` monta o contexto e entra no loop. Em `:212`, qualquer
resposta sem ferramentas encerra normalmente o turno. `Message` guarda texto,
chips e batch; ausência de escrita não é distinguida de um pedido de organização
atendido. As ferramentas têm parâmetros estruturados, mas o resultado da intenção
do turno não tem contrato correspondente.

No caso observado no app, a primeira resposta teve apenas `list_steps`, batch 0.
Após a cobrança explícita, o aplicativo executou `list_steps`, três `add_step` e
`reorder_steps`, batch 120. Isso foi leitura do histórico local, não escrita da
auditoria. As screenshots fornecidas por João mostram a mesma diferença. As
ferramentas conseguem executar; a entrada implícita não é confiável.

O teste real fictício abaixo reforça o limite sem reproduzir exatamente o mesmo
resultado: houve escrita, mas uma operação semanticamente inadequada. Portanto
“forçar pelo menos uma escrita” também não seria uma solução correta.

**Consequência:** não basta outro prompt dizendo para agir. Precisamos distinguir
intenção conversacional de organização e verificar a consequência no estado,
mantendo a possibilidade legítima de só conversar.

### 2. Captura, chat e botão de decompor têm capacidades diferentes

- `internal/cleanup/schema.go:25` só propõe passos novos. `cleanup.go:383`
  acrescenta esses passos e evita descrições iguais por normalização textual,
  inclusive quando um passo anterior está concluído. Não expressa inserir antes
  de um passo, reutilizar sua identidade ou adaptar a ordem existente.
- `internal/agent/tools.go:105` limita o botão de decomposição a `list_steps` e
  `add_step`. Ele é seguro no escopo, mas não consegue dividir um passo existente
  ou reorganizar os demais, embora o chat completo consiga.
- `internal/agent/steps.go:17` expõe operações menores no chat completo. A IA
  precisa compô-las corretamente: acrescentar, editar, reordenar ou substituir.
  `service.SplitStep` substitui a ação original por tombstone; não é sinônimo de
  inserir uma preparação antes dela.

No ensaio real, `split_step` substituiu “Meditar por três minutos” pela preparação
de água/banheiro/alongamento/flexões. O serviço preservou histórico e estado
corretamente, mas não tinha como saber que isso não decompõe meditar. O passo
antigo não ficou duplicado como pendente: está excluído logicamente. É erro de
decisão da IA, não evidência de corrupção na rotina de reordenação.

**Consequência:** os três caminhos precisam compartilhar significados e um caminho
de aplicação de organização, com modos/capacidades explícitos. Isso não exige
que todo clique possa apagar ou alterar qualquer coisa. Repetir uma ação concluída
deliberadamente também precisa ser distinguido de duplicar acidentalmente uma
pendência; não se resolve reabrindo progresso antigo nem por deduplicação textual.

### 3. Parte dos fatos novos foi delegada ao modelo, e os contratos divergem

`internal/facts/facts.go` recebe passos, mas `TaskFacts` oferece essencialmente
estimativa total, bloqueios e anotações. `domain.TaskEstimate` é correto para o
total acordado: maior entre estimativa explícita e soma conhecida, sem confundir
concluído com excluído. Não calcula progresso/tempo por ação.

`planner/planner.go:135` serializa `Service.Execution` inteiro e o anexa ao prompt.
Esse estado inclui todos os passos (também excluídos), sessões e correções sem
recorte histórico (`service/focus.go:136`). O modelo precisa agregar registros,
aplicar regras de relógio/estado e separar relevante de histórico. O prompt pede
avaliar o trabalho que falta, mas o código entrega o total mais registros brutos.

Outras superfícies mostram recortes diferentes: `agent.briefOf` em
`tools.go:176` retorna só a estimativa explícita; a captura descreve total e estado
dos passos, mas não seus minutos individualmente (`cleanup.go:152`); o frontend
recalcula total em `web/src/lib/execution.ts:44` apesar do mapa `estimates` na API.
As fórmulas Go/TS estão equivalentes hoje; duplicação é risco de divergência,
não uma divergência numérica demonstrada nesta rodada.

Snapshot local: aproximadamente 7,5 KB de JSON de execução, 25 linhas de passos e
7 sessões. Isso não prova problema de desempenho atual. A consulta sem limite
histórico é um risco de crescimento, não motivo para infraestrutura nova.

**Consequência:** criar uma projeção pura de fatos de execução para os consumidores:
estimativa total e sua incompletude, pendências/concluídos, estimativas conhecidas
das pendências, tempo registrado por alvo, foco atual e referências candidatas.
Não converter tempo registrado em porcentagem de tarefa realizada nem presumir
que uma estimativa excedida significa conclusão.

### 4. Não existe uma próxima ação compartilhada entre todas as superfícies

`planner.Pick` contém task_id/step_id e pode escolher um passo que não é o primeiro
pendente. O fallback usa o primeiro. Isso é uma sugestão válida, não uma fila.
No chat, `tools.go:540` retorna o Pick ao modelo para sua resposta em texto; não
o transforma em um resultado de interface endereçável pela mensagem.

`web/src/lib/plan.ts:24` monta Ações com ordem de prioridades/plano e foco ativo
primeiro. `TaskCard.svelte:38` e `ExecutionSteps.svelte` escolhem o passo em execução
ou o primeiro pendente. Não recebem o Pick do chat. O botão E agora tem seu próprio
estado de sugestão e só inicia foco por uma ação explícita do usuário.

O commit `147777b` melhorou a passagem do pedido atual quando o agente chama
`what_now`. Ele não muda a decisão de chamar ferramentas, não integra uma sequência
e não dá ao botão comum o histórico do chat. Era uma correção de fronteira parcial.

**Consequência:** separar claramente plano do dia, sugestão e foco ativo, mas usar
uma referência de ação consistente nas saídas/entradas. Alinhar escolher uma ação
existente versus precisar prepará-la antes; sugerir não pode fingir decomposição.
Não criar automaticamente uma segunda fila persistente ou replanejar Must/Should/
Could a cada passo concluído. A estabilidade do plano continua um requisito útil.

### 5. O histórico atual registra execução, mas não fecha o aprendizado

`domain.Step` não tem horário próprio de conclusão. FocusSession registra o alvo,
tempo acumulado e trecho em execução, mas não um snapshot da estimativa/descrição
nem cada intervalo de pausa/retomada. Atividades antes/depois permitem recuperar
parte do passado, sem constituírem hoje um contrato explícito de calibração.

`planner/close.go:174` produz o resumo do dia a partir de títulos de tarefas,
eventos e triagem. Uma tarefa ainda aberta com cinco passos feitos e meia hora de
foco continua aparecendo como inacabada sem esse progresso no resumo. O agente
lê resumos recentes, mas eles não contam essa parte da execução.

**Consequência:** primeiro incorporar progresso/tempo registrado à leitura do dia.
Calibração de estimativas vem depois, com identidade/versão da ação e estimativa
vigente na ocasião. Preservar correções e distinguir dado observado de hipótese.
Timer ligado não prova atividade, distração, energia ou esforço cognitivo.

## Evidência comportamental delimitada

[Resultado estruturado dos dois ensaios](evaluations/2026-10-27-assistance-boundary.json).
Modelo configurado: `gpt-6.1-sol`; um ensaio por caso, em SQLite temporário criado
do zero, sem copiar perfil, conversas, tarefas ou integrações do usuário. Somente
a configuração do provedor foi reutilizada. Ambos terminaram sem erro de API.

| Caso | Resultado observado | Interpretação |
| --- | --- | --- |
| Travado, com sequência de ações concretas antes de arrumar o espaço | Leu passos e chamou split_step; criou quatro ações, preservou concluídos e não iniciou timer. Substituiu a meditação pendente e deixou novas estimativas desconhecidas. | Capacidade técnica presente; transformação semântica de um passo não relacionado não atende à preservação pretendida. Não é reprodução exata da resposta só verbal do app. |
| Desabafo explícito, sem organizar/alterar tarefas | Conversou, zero ferramentas/escritas, mesmos seis passos e duas tarefas. | O caso de controle foi respeitado nessa amostra. Não transformar qualquer “estou travado” em mutação obrigatória. |

Uma execução por caso não é medida de confiabilidade. Código de saída verde,
quantidade de ferramentas ou existência de um batch não comprovam que o pedido
foi atendido corretamente. O resultado persistido e seu significado são o critério.

## O que preservar

- Tarefa como intenção/plano, passos leves em um nível e execução opcional direta.
- Serviço único, validação e histórico reversível para alterações de organização.
- Tempo real/correções independentes de undo, uma execução ativa, relógio do
  servidor e controles manuais sem depender da IA.
- Plano de prioridades estável; SSE e recuperação simples dos clientes.
- IA decide contexto/granularidade; não substituir isso por regex, classificação
  clínica, duração fixa ou contagem universal de passos.

## Ordem recomendada para a correção — proposta, não implementada

1. **Fatos de execução compartilhados.** Derivar no Go o estado útil de cada alvo
   e oferecer recortes adequados a captura, planejamento, chat e resumo do dia.
   Separar total original, estimativas pendentes conhecidas e tempo registrado;
   eliminar cálculo de domínio por interpretação de JSON no modelo.
2. **Contrato de organizar/preparar uma ação.** A IA identifica intenção/alvo e
   propõe acréscimos, reutilização e ordem; o serviço valida/aplica e devolve o
   estado efetivo. Reaproveitar as operações existentes. Tornar verificável a
   diferença entre criar passos, acrescentar antes deles e dividir um passo.
   Conversa sem organização continua sendo um resultado legítimo. Não criar
   um framework genérico de operações ou duplicar o domínio.
3. **Conectar o resultado à execução.** Chat, E agora e Ações devem referenciar
   tarefa/passo reais de maneira consistente. Começar continua explícito;
   conclusão/adiamento pede uma decisão nova quando fizer sentido, sem fila
   rígida nem troca automática do timer.

O critério da primeira entrega integrada deve cobrir, em conjunto: pedido implícito
de organização, tarefa existente com progresso, repetição deliberada de uma ação,
ação já concreta, desabafo sem mutação, falha parcial/undo e preservação do foco.
Usar testes existentes para integridade e poucas avaliações reais para a escolha
do modelo; não simular as ferramentas desejadas e chamar isso de teste de intenção.

## Verificação e alcance

Releitura da especificação original, módulos acima e seus testes. Os sete pacotes
Go `domain`, `facts`, `prompt`, `service`, `cleanup`, `planner` e `agent` passaram
com seus testes existentes (cache reutilizado quando aplicável). Nenhum teste novo
foi adicionado, nem QA visual executado.

O serviço real foi consultado somente para histórico recente e estado/configuração
necessários ao diagnóstico; nenhum POST/PUT/DELETE foi enviado. O ajuste observado
após a cobrança veio do próprio uso de João no chat. Não houve restart, migração,
edição de prompts, alteração de código do app ou modificação de tarefas reais.
Os bancos fictícios foram removidos pelo harness ao terminar.
