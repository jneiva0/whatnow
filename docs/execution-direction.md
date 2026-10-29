# Execução assistida — direção e diagnóstico

## Retomada de 29/10 — plano autorizado

### Ponto de retomada — manter o plano visível

João pediu explicitamente preservar o plano geral enquanto traz feedback pontual,
pois pode esquecer onde estávamos. Novos refinamentos entram na fatia atual sem
substituir silenciosamente o objetivo. Ao fechar uma rodada, informar o que foi
entregue e qual é a próxima fatia deste plano.

- **Objetivo:** ajudar a escolher uma ação viável, começar e continuar.
- **Entregue no código:** faixa persistente/cartões compactos; escolha contextual
  automática em Ações, com Não planejado e captura recente; ajuda por perguntas
  no chat para esclarecer tarefas.
- **Próximo recorte de uso:** abrir Ações, conferir a escolha e usar Começar;
  se faltar clareza, usar Me ajude a organizar e responder uma pergunta.
- **Depois:** avaliar continuidade pelo uso. Migração ampla, calibração histórica
  e voz permanecem posteriores; o pequeno bug PWA não assume a prioridade.

João confirmou instalação PWA funcionando, com pequeno bug de faixa azul no
viewport standalone ainda sem diagnóstico. O aceite de Ações revelou que a
ordem por prioridades não esclarece qual ação começar; o ritual capturado para
destravar ficou em Não planejado. Relato completo em `feedback.md`.

1. Entregar faixa persistente com ação em execução, tempo acumulado, pausa e
   retorno ao Foco fora dessa tela, consumindo o estado compartilhado existente.
2. Compactar cartões e metadados: aproveitar a mesma linha quando houver espaço,
   energia como indicador que abre seleção ao clicar, estimativa editável ao
   clicar e contexto sob expansão explícita. Preservar passos e controles.
3. Próxima fatia de produto: fazer Ações apresentar uma escolha contextual útil
   incluindo Não planejado, sem exigir pedido extra no chat. A forma de seleção
   ainda está aberta; não implementar ordenação universal por menor duração.

Os itens 1–3 estão implementados nesta rodada, com as verificações e limites
descritos abaixo. Migração ampla, calibração histórica e voz continuam posteriores.

### Próxima implementação — escolha contextual em Ações

**Estado após execução de 29/10:** plano abaixo implementado. Sugestão compartilhada
com E agora, atualização após captura aplicada, continuidade do foco e tratamento
de respostas antigas/falhas verificados. Sem novo endpoint/schema ou fila.
Motivo do novo caminho usa `context` com `reason` vazio; caminho legado continua
compatível. Datas de tarefa aparecem no contexto; captura recém-aplicada é
identificada como entrada recente e não como sentimento permanente.

João acrescentou ajuda por perguntas para tarefas sem clareza e pendências antigas.
Entrada no painel abre chat com pedido de uma pergunta útil por vez; instruções
orientam esclarecer antes de alterar e aplicar apenas o que a resposta sustenta.
Não há interrogatório para preencher todos os campos nem revisão forçada por idade.

Verificação: Go test/vet completos;191 testes frontend, check e build; quatro
cenários Chromium com banco temporário. Autofixer sem issues nos novos componentes
e estado; sugestões preexistentes de Maps locais e efeitos de navegação/sincronização
foram revisadas e mantidas. Detector sem achados; capturas mobile/desktop inspecionadas.
Avaliação com provedor real em `evaluations/2026-10-29-contextual-actions.json`
preserva o primeiro resultado que ignorou o ritual e a repetição após ajuste,
além de falha transitória do provedor durante esclarecimento. Não é garantia
de qualidade geral nem aceite no celular físico de João.

O plano detalhado abaixo é o contrato da entrega, não uma lista ainda pendente.

#### Critérios e dados incompletos — proposta explicada em 29/10

João pediu esclarecer prioridade versus energia/tempo/sentimento e o caso de
cadastros incompletos antes de prosseguir. O contrato abaixo é a proposta de
implementação, não evidência de que a qualidade dessas escolhas já foi validada.

- Prioridade orienta compromissos do dia; sugestão escolhe a ação viável agora.
  Sugerir uma ação de Não planejado não altera seu nível de prioridade.
- Energia da tarefa é esforço esperado; pique do dia/contexto atual representa
  a disponibilidade relatada de João. Sentimento da tarefa e humor do momento
  também são distintos. Não inferir um a partir do outro nem perpetuar um relato
  transitório como estado atual.
- Campos ausentes são desconhecidos, sem impedir execução nem criar punição
  automática. A IA pode interpretar a descrição para sugerir, mas não apresentar
  duração/energia inferidas como fatos preenchidos, inventar sentimento ou salvar
  metadados durante uma consulta de próxima ação. Captura/decomposição já podem
  propor estimativas aproximadas e constituem operações diferentes.
- Considerar compromissos/prazos, bloqueios, próximos passos, disponibilidade e
  pedido atual em conjunto. Ações curtas podem ajudar a iniciar trabalho relevante;
  menor duração não é ordenação universal. Metadados da tarefa não se transferem
  mecanicamente a cada passo. Tempo registrado não comprova trabalho concluído.
- Recomendar uma ação existente com justificativa curta. Tarefa concreta funciona
  sem passos ou metadados; tarefa vaga pode precisar da organização já existente.
  A recomendação não deve fingir ter criado um passo apenas por descrevê-lo.

Escopo técnico: ampliar o contexto/instruções e contratos do planejador Go já
existente; validação de IDs, elegibilidade e comandos permanece no serviço.
Frontend controla solicitação, validade da resposta e apresentação, sem pontuação
local própria. Inicialmente não há necessidade demonstrada de novo schema ou fila.
Antes da entrega, verificar casos com todos os metadados ausentes, parcialidade,
conflito entre prioridade e esforço, relato atual versus histórico e falha da IA.

Plano preparado após o fechamento dos cartões, em 29/10. É a próxima fatia do
objetivo de ajudar a começar; os refinamentos visuais acima não a substituem.

1. **Expor uma escolha concreta na própria aba Ações.** Reaproveitar
   `Planner.NextAction` e `/whatnow`, que já consideram tarefas elegíveis dentro
   e fora do plano. Mostrar ação, tarefa/contexto, estimativa e justificativa curta,
   com Começar e Outra sugestão, sem exigir que João repita o pedido pelo chat.
   Se há execução ativa, dar continuidade a ela antes de propor uma troca.
2. **Manter contexto e continuidade.** Conferir o percurso captura → sugestão
   para o pedido de destravar chegar à seleção. Atualizar sugestões que perderam
   validade após conclusão/edição, sem reescolher a cada tick ou apagar a escolha
   enquanto João interage. Usar as operações/fatos Go existentes; qualquer lacuna
   de contexto deve ser resolvida na fronteira responsável, não com ranking local
   duplicado. Não adotar menor duração como regra universal nem garantir que um
   ritual específico sempre vença.
3. **Verificar o ciclo completo em dados descartáveis.** Cobrir ritual não
   planejado com tarefas prioritárias concorrentes, tarefa já concreta sem passos,
   foco ativo, ação bloqueada/concluída, falha do provedor e atualização após
   concluir. A lista manual continua utilizável durante falha da sugestão.
   Diferenciar contratos com modelo simulado de qualidade com provedor real.
4. **Entregar uma fatia revisável.** Changelog e commits coesos, estado de
   implementação/validação explícito e próximo recorte de uso curto. Preservar
   timers, dados, undo e sincronização; não criar nova fila ou planejamento paralelo.

Primeira ação técnica: consolidar o ciclo de vida da sugestão hoje local ao
`+page.svelte` e definir como a entrada em Ações pede uma escolha sem chamadas
repetidas. O plano acima ainda não foi implementado nesta rodada de alinhamento.

Atualizado em 27/10/2026. Fonte: feedback direto de João nesta conversa e leitura do código. O estado atual está na seção abaixo; as seções de diagnóstico e entregas posteriores preservam o histórico e suas evidências, sem transformar antigas próximas ações em pendências atuais.

## Estado atual — voltar ao produto

**Nova direção autorizada:** João escolheu Ações/Foco primeiro, seguida da
verificação de instalação PWA, que não impede o uso. Isso substitui o adiamento
da fatia visual correspondente e a pergunta de prioridade registrada abaixo.
Primeira fatia entregue: passo destacado, contexto pai recolhível, metadados
compactos compartilhados e timer dominante no Foco. Detalhes no changelog.

**PWA:** diagnóstico de fonte confirmado: falta manifesto vinculado e ícones
de instalação. Não tratar ausência de service worker como bloqueio suficiente:
o Chrome retirou a exigência de fetch handler para instalação pelo menu no
Android108. Referência: https://developer.chrome.com/blog/update-install-criteria .
O próximo ajuste técnico seria declarar nome, ícones, start_url, scope e display;
isso não implica cache offline nem alarmes em segundo plano. Não foi implementado
nesta fatia, cuja prioridade é o aceite de Ações/Foco.

### Revisão de fechamento do fluxo principal — retomada de 27/10

Esta tabela prevalece sobre pendências históricas narradas nas seções abaixo.
Não é autorização para reabrir frentes adiadas nem afirmação de que todo uso
possível foi validado.

| Contrato | Evidência atual | Situação |
| --- | --- | --- |
| Organizar passos na tarefa existente e iniciar pelo chat | Feedback direto de João para organização, botão e início correto no Foco; cenários sintéticos com provedor real em evaluations/2026-10-27-shared-organization.json | Aceito no recorte de uso relatado |
| Timer ativo governa seleção; troca por play preserva tempo | Regra confirmada por João; focusTarget, Service.StartFocus e testes de troca/concorrência | Implementado; escolha explícita com timer parado corrigida em4cfd756 |
| Conclusão explícita, próximo parado e undo sem retroceder tempo | TestFocusClockSwitchIdempotencyAndUndo, TestCompleteStepWithoutStartingAndKeepActionTimeAcrossSegments e testes de representação | Verificado por testes, não novo aceite físico |
| Recuperação/sincronização sem descartar comandos | execution-state.test.ts e today.test.ts; relato anterior de sincronização entre dispositivos | Testes de corrida, suspensão e falha; sem promessa de funcionamento offline |
| Entrada web consistente e independente das regras Go | focus-navigation.ts usado por quatro entradas;177 testes frontend, check/autofixer/build; API REST e TUI usam serviço Go | Consolidado ema1dcb2d; erros de play/pausa visíveis |
| Uso mobile e entrega atual | Aceite de overflow por João; just dev no tmux; HTTPS7001 retornando200 | Versão em desenvolvimento disponível |

Revisão final do backend: `go test ./...` e `go vet ./...` passaram, incluindo
serviço, API, planner, agente e cliente terminal. Não houve mudança no terminal,
novas integrações, schema ou dados reais nesta revisão.

Próxima escolha enviada a João: retomar refinamento coeso de Ações/Foco,
priorizar instalação Android/PWA ou manter apenas correções orientadas pelo uso.
A migração ampla e a calibração histórica continuam adiadas enquanto não houver
nova direção. O objetivo autônomo permanece ativo; não usar testes verdes como
justificativa para acrescentar funcionalidades sem necessidade.

**Diagnóstico de instalação Android, somente leitura:** `web/static` contém
apenas robots.txt; não há manifesto, ícones de instalação ou service worker no
fonte. O HTML servido por HTTPS7001 não contém link de manifesto. HTTPS respondeu
200. Isso identifica uma lacuna de implementação, não uma falha comprovada do
aparelho. Caso PWA seja priorizado, separar instalação/abertura em janela própria
de operação offline: os comandos continuam dependendo da API Go. Nenhum cache,
worker, permissão ou promessa de alarmes em segundo plano foi adicionado.

**Retomada e entrega corrigida — 27/10:** João confirmou o overflow horizontal
corrigido no Android. A correção está registrada no [changelog](../CHANGELOG.md).
A conferência anterior abaixo que declarava a versão atual em uso estava errada:
o processo ainda executava `~/go/bin/whatnow` antigo. Nesta retomada ele foi
substituído por `just dev` no tmux, usando o banco existente. O Vite atende em
5173 e a API recompilada em7777 agora retorna `progress`. HTTPS Tailscale7001
encaminha para o Vite e respondeu200. Alterações frontend usam HMR; alterações
Go precisam de reinício da receita, que não contém watcher.

A organização pela IA já recebeu feedback positivo de João. Nesta retomada ele
também confirmou o botão “Começar esta ação” e a abertura/início correto no Foco.
O caminho chat → ação → foco está aceito em uso. A revisão de continuidade
reproduziu e corrigiu a seleção explícita ignorada após concluir um passo:
outro passo da mesma tarefa e execução direta agora respeitam a escolha,
preservando a prioridade de um timer rodando e o próximo passo parado.
Teste de regressão falhou antes da correção;174 testes frontend, check Svelte
e testes Go direcionados de execução passaram. Não retomar a migração ampla da UI
nem calibração histórica como substitutos desse ciclo. Manter detalhes técnicos
e mudanças neste repositório, com commits por fatia.

**Aviso da estimativa — 27/10:** o cliente acompanha a mesma estimativa e tempo
acumulado usados no Foco, inclusive segmentos anteriores. Ao atingir o tempo
previsto de uma ação que estava sendo acompanhada, avisa uma vez e oferece Ver
foco. O aviso usa as notificações já permitidas pelo navegador, com o toast
existente como alternativa (inclusive quando o construtor de Notification não é
suportado no celular). Não pede nova permissão, não toca áudio, não conclui,
pausa ou inicia ações. Abrir uma página com um timer já excedido não repete um
aviso antigo. A aba suspensa pode avisar quando volta a executar; com a página
fechada não há processo novo para alarmes em segundo plano.
Três testes relevantes cobrem soma dos segmentos, aviso único/ausência de
estimativa e recuperação da aba sem comandos de timer. Os173 testes frontend,
Svelte check, autofixer e build passaram. Detector Impeccable sem achados.
Notificações/gestos em aparelho físico e QA visual não foram executados.

**Continuidade entre dispositivos — retomada autônoma 27/10:** João autorizou
avançar nas partes de maior impacto enquanto testa e descansa, sem transformar
exemplos em regras rígidas. Foram reproduzidas e corrigidas duas condições de
corrida no cliente: avisos SSE recebidos durante uma gravação de foco eram
descartados, e uma leitura antiga do Dia podia substituir uma leitura nova.
Agora uma atualização pendente é buscada após a gravação e respostas antigas são
ignoradas. Erro de comando continua visível mesmo se uma leitura seguinte funciona;
atualizar dados não equivale a confirmar uma gravação que perdeu a resposta.
Três testes de regressão reproduziram as falhas antes da correção. Os170 testes,
check Svelte e build passaram; banco, API e timer não foram alterados.

**Conferência de entrega — 27/10:** o serviço local em7777 já serve HTML idêntico
ao build atual e sua API expõe progresso compartilhado e estimativas do servidor.
Verificação somente por leitura; não foi necessário reiniciar ou alterar dados.
Foram relidos os cenários de preservação de foco/tempo/undo, validação REST/SSE,
falha da IA com execução manual e referências de ação no histórico. Os ensaios
reais fictícios continuam delimitados ao registro abaixo. O próximo aceite é
o uso de João: pedir organização, conferir os passos salvos e iniciar a ação.
O reinício recomendado anteriormente já não é necessário para esse teste.
Não confundir entrega técnica com aprovação visual ou eficácia confirmada no uso
cotidiano. A UI ampla e a calibração histórica permanecem adiadas.

**Fatos na apresentação — 27/10:** cartões, metadados e foco agora leem a
estimativa agregada de `ExecutionState.estimates`, calculada no Go. Não refazem
a regra do total a partir de listas de tarefas/passos que podem atualizar em
momentos distintos. Ausência do snapshot permanece desconhecida. O cálculo local
foi delimitado como `draftTaskEstimate`, usado apenas na prévia do editor enquanto
o valor ainda é um rascunho; o servidor continua validando a gravação. O contador
por segundo segue derivado dos horários confirmados, como antes.

Verificação: 167 testes frontend, check Svelte sem erros/avisos, autofixer oficial
dos cinco componentes/módulos tocados sem achados e build passaram. O teste de
recuperação existente foi ampliado para confirmar o total/indicador parcial do
servidor mesmo com a lista de passos ainda vazia e depois de uma falha de API.
Nenhum novo teste cosmético, mudança visual ou mecanismo de sincronização.

**Fatia de fatos compartilhados — 27/10:** `facts.Execution` calcula estimativa
total, estimativas pendentes conhecidas, passos concluídos, tempo registrado por
ação, foco em execução e próxima referência. A API, planejamento, captura,
consulta de tarefas pelo chat e resumo do dia usam essa projeção. O modelo não
precisa mais agregar o histórico bruto de sessões. O resumo distingue tempo
acumulado de tempo atribuído ao dia; calibração histórica continua posterior.
Testes dos pacotes facts, prompt, service, cleanup, planner, agent e httpapi
passaram. Commit `5cdaf85`.

**Organização compartilhada — 27/10:** captura, chat e botão de decompor agora
usam `Service.OrganizeSteps`: sequência com referências a IDs existentes e ações
novas; omissões são preservadas depois da sequência. A operação é atômica,
reversível e não apaga, conclui, reabre ou troca foco. Preparar antes de trabalhar
não exige substituir um passo sem relação. Repetição deliberada é uma ocorrência
nova, preservando a concluída. Edições e divisão real continuam operações próprias.
O resultado traz progresso aplicado e próxima referência, calculados no serviço.

A suíte Go completa passou. [Avaliação sintética com o provedor real](evaluations/2026-10-27-shared-organization.json)
cobriu pedido implícito, desabafo, captura continuando tarefa existente e botão de
decompor. No caso final, a preparação foi salva primeiro com estimativas; tarefas,
passos anteriores e undo foram preservados. As primeiras amostras e refinamentos
estão registrados no arquivo: isso não demonstra confiabilidade em toda conversa.
Commit `7d59ab7`.

**Chat → execução — 27/10:** resultados bem-sucedidos de organização e E agora
carregam `ToolUse.action` (task_id/step_id) e um resultado organizado/sugerido. A
referência persiste no JSON já usado pelo histórico e chega por SSE, sem migração
de banco. O chat oferece “Começar esta ação” para um alvo ainda disponível; ao
clicar, usa a operação existente de foco e abre a tela correspondente. Nada inicia
por uma resposta do modelo. Erros, passos concluídos/removidos e alterações
desfeitas não oferecem o início de um alvo inválido. Não foi criada fila extra.
O contexto inclui também o foco selecionado quando pausado, além do em execução.

Verificação final: suíte Go completa; 167 testes frontend; Svelte check sem erros
ou avisos; build frontend e executável Go passaram. O autofixer oficial aprovou
ChatAction; em ChatDock, a colocação de const foi corrigida e as sugestões sobre
scrollTo no efeito/bind:this preexistentes foram revisadas e mantidas por sincronizarem
o DOM, sem reatribuir estado. QA visual e teste em aparelho físico não realizados,
conforme combinado com João. Nenhum serviço reiniciado ou dado real alterado.

Os cinco ensaios com provedor real usaram apenas dados fictícios: organização
implícita, desabafo, captura existente, decomposição pelo botão e tarefa já concreta.
No último, E agora escolheu a própria tarefa sem novos passos ou timer automático,
devolvendo sua referência ao histórico. As amostras são delimitadas, não garantia
de qualidade geral. UI mais ampla, calibração histórica e atribuição de tempo por
dia permanecem no plano. Próximo aceite: reiniciar na nova compilação e testar um
pedido natural de organização, verificando os passos salvos e o botão de começar.

**Atualização após o teste de uso:** a passagem de contexto de `147777b` é uma
correção parcial. João mostrou que o chat ainda respondeu só verbalmente e só
organizou os passos após cobrança. Pediu revisar a arquitetura antes de continuar
com correções pontuais. A [auditoria de integração](execution-architecture-audit.md)
registra os módulos inspecionados, lacunas, ensaios reais fictícios e a proposta
de sequência. Próximo trabalho recomendado: fatos de execução compartilhados,
contrato de organização verificável e ligação consistente entre sugestão e ação.
Nenhum desses três itens foi implementado durante a auditoria; o progresso posterior
está registrado acima. UI e calibração
histórica continuam guardadas para depois. Graphify fica fora do fluxo normal;
João esclareceu que era um experimento e questionou o custo das buscas laterais.

João considera a UI suficiente por enquanto e pediu guardar o restante para retomar
em breve. **A migração visual fica estacionada**, sem iniciar agora a próxima
fatia do editor. O objetivo continua ajudar a decidir, começar e continuar; a
biblioteca de componentes é um meio de reduzir atrito e manutenção.

| Frente | Estado em 27/10 | O que ainda precisa de atenção |
| --- | --- | --- |
| Planejamento e IA | Captura/chat trabalham com passos e minutos, podem continuar uma intenção existente e E agora pode escolher um passo. Integração implementada e avaliada nos cenários fictícios com provedor real descritos abaixo. | Qualidade no uso cotidiano: reconhecer contexto, evitar fragmentação e sugerir uma ação que ajude a começar. Os cenários exercitados não garantem isso em toda conversa. |
| Execução | Tarefa direta ou passos, estimativas, tempo real persistente, conclusão explícita e próxima ação apresentada parada. João confirmou sincronização fluida entre computador e celular. | Acompanhar o uso e corrigir quebras reais de continuidade; não reconstruir o timer nem introduzir duração de sessão separada. |
| Apresentação | Prioridades e Ações são duas perspectivas do mesmo plano; tarefa/passos compartilham controles. Ações filtra elegibilidade, mantém a ordem do plano e apresenta o foco ativo primeiro. | Ações ainda não é uma sequência contextual sugerida pela IA. Essa escolha contextual ocorre no E agora; a relação entre escolher, executar e continuar é o próximo ponto de discussão de produto. |
| Fundação da UI | Bun, Tailwind/tokens e shadcn-svelte/Bits UI em controles compartilhados, menus e metadados; fechamento do chat confirmado por João. Commits locais até `14e11c4`. | Trabalho restante registrado em [ui-migration.md](ui-migration.md#remaining-work--parked-on-2026-10-27). Retomar em fatias próprias, sem exigir uma migração completa antes de usar o app. |

### Próxima frente de produto — em alinhamento

Revisar o ciclo **“estou travado → recebo uma ação viável → começo → concluo ou
me distraio → consigo continuar”**. A base técnica já existe; o próximo ganho deve
vir de conectar melhor a ajuda da IA à execução, conforme o uso mostrar onde
esse ciclo falha. Não há decisão de criar fila nova, rotinas fixas, ordenação por
menor duração ou outra página.

João confirmou essa direção e esclareceu que a granularidade depende da dificuldade
de iniciar naquele momento. A próxima fatia conecta esse contexto à sugestão;
o escopo e os limites estão abaixo. Feedback de aparência permanece no backlog;
perda de dados, quebra de execução ou de sincronização continua justificando uma
correção imediata. Voz permanece um bônus posterior ao fluxo principal.

### Granularidade e contexto de ajuda — 27/10/2026

**Direção confirmada:** não existe quantidade ideal universal de passos. Quando
a ação já é executável, manter unidades significativas. Quando há dificuldade,
detalhar o próximo movimento observável, sem transformar a tarefa inteira em
uma checklist exaustiva. Estimar a ação real proposta; diminuir só o número de
minutos não resolve uma ação vaga. Passos continuam em um nível, com identidade,
ordem, estimativa opcional e conclusão explícita. Energia e sentimento são
independentes. Não há escala clínica, rotina fixa ou classificação permanente
do usuário sendo introduzida.

**Lacuna verificada:** o agente recebe o histórico da conversa, mas delegava
`what_now` sem argumentos ao planejador. Esse planejador recebia estado do dia,
notas, tarefas/passos e registros de execução, sem a mensagem corrente nem o
histórico do chat. Portanto o pedido “estou travado” podia se perder nessa
fronteira mesmo quando as instruções do modelo falavam em ganhar movimento.

**Implementado nesta fatia:** `planner.ActionRequest` representa o pedido atual
com contexto e exclusão opcional de tarefa. O agente preserva a mensagem original
na chamada delegada e pode fornecer contexto anterior relevante pelo argumento
estruturado `what_now.context`. O endpoint `/whatnow` aceita `context` opcional;
o botão comum continua funcionando sem ele. A necessidade é válida para aquela
requisição, não grava automaticamente estado de bloqueio ou sentimento no pai.
Fallback e validação de tarefa/passo permanecem. Instruções de captura,
decomposição, chat e sugestão explicam granularidade adaptável, sem impor número
de passos ou minutos. Isso é orientação ao modelo, não garantia de qualidade da
decomposição nem um mecanismo novo que cria passos automaticamente ao clicar
E agora. Nenhuma alteração de timer, plano, schema SQLite ou dados em uso.

Verificação desta fatia: suíte Go completa passou; os novos cenários verificam
chat → planejador com mensagem original e resumo opcional, API com contexto e
ausência de gravação desse estado transitório no sentimento, além de ausência de
vazamento para outro pedido. São testes com modelo simulado, não avaliação nova
de qualidade com provedor real. Os 167 testes frontend e check Svelte sem
erros/avisos passaram. O processo Go em uso não foi reiniciado; a alteração de
backend será carregada no próximo reinício habitual.

**Limite atual:** o E agora aberto diretamente pela UI ainda não recebe todo o
histórico do chat. A nova passagem resolve a delegação pelo agente e fornece um
campo explícito de requisição; não declarar que todas as superfícies agora
compartilham um estado persistente de apoio. Decidir essa persistência somente
quando houver uma necessidade concreta de continuidade entre pedidos.

### Histórico para melhorar sugestões — proposta posterior

O modelo atual já guarda registros de foco por tarefa/passo com horários e tempo
acumulado, correções explícitas e atividades antes/depois das edições. Porém
`Step` não tem horário próprio de conclusão, `FocusSession` não guarda uma cópia
da descrição/estimativa no início e os trechos retomados são acumulados; não é um
histórico completo de intervalos para atribuir tempo a cada dia. O planejador
recebe o JSON bruto da execução; isso não é calibração histórica implementada.

Antes de prometer aprendizado, a próxima etapa de dados deve preservar a
estimativa vigente na ocasião e a identidade/versão da ação comparada, respeitando
divisões/edições e correções. Reaproveitar o histórico existente quando suficiente;
adicionar somente fatos cuja ausência impeça a comparação desejada. Derivar totais
e comparações no Go e oferecer ao modelo um resumo limitado e estruturado.
Poucas observações ou casos não comparáveis devem continuar incerteza explícita.
Tempo com timer ligado é tempo registrado, não prova de execução, distração ou
conclusão; não inferir esses fatos nem alterar tempo real automaticamente.

Sequência sugerida: contexto/granularidade agora; continuidade da próxima ação
conforme uso; histórico confiável e calibração depois. Não construir relatórios,
treinamento de modelo, categorias universais de comportamento ou uma nova fila
de planejamento como pré-requisito para essa primeira entrega.

## Norte

Ajudar João a começar, decidir o próximo movimento e continuar. Exemplos como beber água antes de trabalhar são referências de contexto, não rotinas obrigatórias. Evitar transformar cada sugestão em regra rígida ou cada reclamação visual em uma nova prioridade. Preferir a solução mais simples que ajude a execução.

## Direção confirmada e escolhas abertas

- Tarefa mantém intenção, contexto e continuidade; passos ajudam a executar, com um nível. Reconhecer uma intenção existente antes de criar outra tarefa parecida.
- Minutos substituem S/M/L como medida de tamanho na experiência. Energia e sentimento continuam independentes. Não converter automaticamente dados antigos S/M/L em minutos fictícios; compatibilidade técnica é decisão de implementação.
- A estimativa total considera os passos e não deve ficar abaixo da soma das estimativas conhecidas. Tempo real permanece separado. Estimativas faltantes são desconhecidas, não zero. A fórmula para estimativa manual maior que a soma e a apresentação de totais incompletos devem permanecer simples e explícitas.
- Prioridades e ajuda para começar são perspectivas complementares. Uma visão de sequência é uma hipótese útil, não uma exigência de página nova: decidir a forma depois de corrigir a escolha de ações. Duração, energia, sentimento, contexto, dependências e urgência informam sugestões; não ordenar mecanicamente pela menor duração nem criar rituais fixos.
- Decomposição deve responder à dificuldade de iniciar e às ações já descritas pelo usuário, sem exigir uma palavra-chave ou decompor toda captura. Evitar burocracia para tarefas já executáveis.

## Diagnóstico inicial verificado no código (antes da implementação abaixo)

1. **Captura não produz passos nem estimativas em minutos.** `internal/cleanup/schema.go` e `result.go` só representam os campos antigos. `instructions.go` pede tarefas/eventos/bloqueios e S/M/L. A situação em `cleanup.go` mostra títulos das tarefas abertas, sem passos. Existe instrução para não duplicar e união por ticket, mas isso não equivale a reconciliar intenções domésticas pelo contexto. O relato de tarefas duplicadas não foi reproduzido nem os dados reais alterados.
2. **Chat tem ferramentas, mas o contrato é inconsistente.** `internal/agent/steps.go` permite adicionar/editar/reordenar/dividir passos. `instructions.go` restringe decomposição a pedido e edição de estimativas a pedido explícito. `create_task` em `tools.go` não recebe estimativa; `update_task` recebe. Isso explica caminhos distintos, não prova incapacidade do modelo de usar as ferramentas existentes.
3. **Planejador recebe progresso, porém ainda decide com regras antigas.** `internal/planner/planner.go` inclui JSON de execução no contexto. Portanto não é correto afirmar que a IA nunca vê passos. Contudo, `internal/facts/facts.go` usa S/M/L para encaixe em intervalos; `internal/prompt/facts.go` apresenta size. Os significados de S/M/L diferem entre cleanup e create_task.
4. **E agora escolhe tarefa; código escolhe passo.** `internal/planner/whatnow.go`: schema da IA contém task_id/reason, e `pickAction` anexa o primeiro passo pendente. A IA não escolhe diretamente uma ação entre os passos. `plan.go` produz prioridades e ordem de tarefas; não cria passos, o que pode continuar sendo uma separação saudável se os outros fluxos forem integrados.
5. **Estimativa agregada não tem regra compartilhada.** `web/src/lib/execution.ts:taskEstimate` prioriza estimativa manual da tarefa, mesmo menor que a soma, e só soma passos quando ela está ausente. Serviços validam valores individuais em `tasks.go`/`steps.go`, sem garantir o piso agregado. Precisamos de significado consistente entre API, IA, planejamento e UI.

## Ordem de entrega proposta

1. **Integrar significado e operações:** estimativa coerente no servidor/contexto, captura/chat capazes de continuar tarefa existente e criar/ajustar passos quando útil. Preservar origem, undo, progresso e foco real. Não fundir tarefas reais silenciosamente nem reescrever todo o pipeline sem necessidade.
2. **Conectar sugestão à execução:** E agora considera ações concretas e progresso, preserva elegibilidade e fallback. Só então decidir a menor apresentação útil de próximas ações; a sequência é flexível e não uma agenda rígida.
3. **Consolidar UI como conjunto:** tarefa/passos, cartões, formulário, controles e tempos com componentes/tokens consistentes. Avaliar shadcn-svelte/Tailwind como base; não foi instalado nem é pré-requisito da etapa 1. Salvar/overflow entra aqui como achado conhecido, sem virar uma entrega isolada.

Feedback novo entra na etapa correspondente. Só perda de dados, quebra de execução ou de sincronização justifica interromper a ordem. QA de uso continua com João, sem exigir repetição de relatos conhecidos.

## Verificação pertinente

Usar os padrões existentes e cenários ligados ao problema: atualizar intenção existente em vez de duplicá-la; decompor tarefa difícil sem fragmentar em cartões; três passos de 15 min implicam ao menos 45 min; estimativa ausente não vira zero; edição/planejamento preservam concluídos e tempo real; sugestão respeita bloqueios e oferece ação executável. Distinguir testes com modelo falso de avaliação com provedor real.

Nesta revisão foram consultados o grafo como índice e os arquivos citados. Não houve alteração de código, migração, chamada de IA real, teste automatizado ou deploy. Próximo trabalho: implementar a primeira etapa como um fluxo integrado, usando estes cenários para limitar o escopo.

## Implementação da integração — 27/10/2026

Na branch `codex/execution-direction`, sem worktree:

- Captura estruturada aceita minutos, passos opcionais e referência explícita a uma tarefa aberta existente. O contexto inclui notas, sentimento, energia e passos concluídos/pendentes. A continuação preserva os dados/progresso existentes, anexa a captura às notas e pode acrescentar passos; tudo integra o lote de desfazer. IDs inválidos/fechados e autorreferência são rejeitados. Descrições iguais normalizadas não viram passos duplicados. Isso habilita reconciliação pela IA, não garante reconhecimento semântico perfeito.
- Chat recebe orientações para reconhecer intenção existente e decompor quando ajuda a começar, sem exigir palavra-chave. `create_task` aceita minutos e os resultados de ferramentas retornam a estimativa explícita. Estimativas reais e correções de tempo continuam separadas.
- `domain.TaskEstimate` calcula o maior entre estimativa explícita e soma dos passos vivos, incluindo concluídos no total original e marcando informação incompleta. A API de execução expõe `estimates`; planejamento usa esse cálculo e minutos em lugar do encaixe por S/M/L. A apresentação existente usa o mesmo significado. Não há migração destrutiva, conversão fictícia dos S/M/L antigos ou mudança no tempo real.
- E agora aceita um passo escolhido pelo modelo, valida pertencimento/estado e mantém fallback para primeiro pendente. Instruções pedem contexto, urgência, energia, sentimento e facilidade para iniciar, sem ordenação rígida por duração.

Verificação: suíte Go completa passou com permissão para portas locais dos servidores de teste; o primeiro ensaio restrito falhou por bind não permitido. Testes direcionados cobrem captura com passos, continuação de intenção com progresso/undo, rollback de captura inválida, piso de estimativa, minutos na criação por chat e rejeição de passo concluído/estrangeiro na sugestão. Check Svelte sem erros/avisos. Testes usam respostas simuladas, não um provedor real. Nenhum dado real, deploy ou processo do aplicativo foi alterado.

Ainda aberto: avaliação qualitativa com IA real; substituição visual de S/M/L por minutos; apresentação simples de próximas ações; consolidação de cartões/formulários/controles. As falhas visuais relatadas continuam na etapa 3, sem repriorização automática. Próxima ação: avaliar os cenários de intenção/estimativa/decomposição com o provedor real em dados descartáveis antes de considerar esta integração validada em uso.


## Avaliação com provedor real — 27/10/2026

Executadas duas rodadas delimitadas com o provedor OpenAI-compatible já configurado (light `gpt-6-luna`, main `gpt-6.1-sol`), em SQLite temporário criado do zero. Somente configuração de provedor/modelos foi lida; não foram copiados tarefas, mensagens, perfil pessoal ou integrações. Cenários fictícios, sem servidor de preview/deploy. O banco temporário foi removido ao final.

- Continuação de arrumação: nova captura dizendo que recolher o chão ficou urgente foi associada à tarefa existente; permaneceu um cartão e os dois passos originais, incluindo o concluído.
- Captura de apresentação com dificuldade para iniciar: criou três passos com 15 minutos cada, total de 45 minutos.
- Chat pedindo primeiro movimento menor: usou `list_steps` e `split_step`, preservou slides/ensaio e não criou outra tarefa. Primeira ação virou abrir uma nota e escrever uma frase (2 minutos).
- E agora selecionou esse passo de 2 minutos, com justificativa ligada ao prazo de amanhã; não utilizou fallback.

A inspeção identificou que o contexto principal omitia notas e captura original. `internal/prompt/facts.go` agora inclui ambas, com teste de regressão; os três pacotes afetados (prompt/planner/agent) passaram. A segunda rodada com o código corrigido manteve os resultados acima. Isso é evidência qualitativa limitada aos cenários exercitados, não garantia geral de classificação, deduplicação ou adequação clínica.

A integração da IA está implementada e esses cenários reais foram avaliados. Próxima etapa permanece apresentação da execução (minutos no lugar de S/M/L, contexto tarefa/passos e próximos movimentos), seguida de consolidação visual, sem transformar comentários isolados em novas prioridades. O aplicativo em uso não foi atualizado nesta rodada.

## Apresentação da execução — 27/10/2026

Primeira parte da consolidação: `InlineEdit` compartilhado salva em Enter/saída e cancela em Escape, mantém o rascunho e informa falhas. Descrição e minutos dos passos usam esse controle, sem botão Salvar. `ExecutionSteps` continua compartilhado entre Dia e Foco, com disposição adaptável. `TaskMetadata` e `TaskEditor` mostram minutos em lugar de S/M/L; os dados antigos permanecem preservados. Cartões mostram próxima ação (ou passo em execução) com tempo próprio e total da tarefa em nível secundário.

Validação: 165 testes frontend passaram; Svelte check sem erros/avisos; build estático passou; autofixer oficial sem pendências nos cinco componentes; detector Impeccable sem achados determinísticos. Não houve QA de navegador ou celular físico, conforme combinado com João; isso não declara aceite visual. Nenhuma biblioteca de componentes nova foi instalada e nenhum processo/deploy foi atualizado.

Próxima unidade: revisar o caminho completo Dia → sugestão → execução e consolidar a composição visual, mantendo estas mudanças junto da integração de IA na mesma branch. A alternativa de sequência continua hipótese de apresentação; não duplicar planejamento nem impor ordem fixa.

### Entrada pela sugestão

Corrigido `WhatNow.svelte`: “Começar” antes somente navegava; com outro foco rodando, a representação autoritativa mantinha a ação anterior. Agora inicia o alvo sugerido pela operação existente do servidor, que pausa o anterior preservando tempo, e só navega após confirmação. Falha mantém a sugestão e apresenta erro; cliques repetidos ficam bloqueados enquanto aguarda. Sem passo, a URL marca execução direta para não escolher um filho diferente silenciosamente. Falha de navegação informa que o foco já iniciou.

Check Svelte sem erros/avisos, autofixer sem pendências, 16 testes existentes de representação/estado passaram. Verificação direcionada do servidor cobre troca, idempotência, undo, elegibilidade e inícios concorrentes. Sem novo teste que apenas espelhe o handler e sem QA visual automatizado.

### Composição da tela de foco

A superfície de foco agora aproveita a largura disponível, em vez do limite fixo de 900px. Ordem preservada: tarefa/metadados → timer/controle → passos/contexto → correções. Timer e play/pause recebem ênfase; `PlayToggle` aceita variáveis CSS de tamanho e cores, mantendo o mesmo componente e interação das linhas. O painel do timer continua sticky, com fundo da superfície para não competir com ela. Em contêineres largos (1000px ou mais), notas, referências e captura original ocupam uma coluna secundária; em espaços menores seguem os passos. Sem contexto, não se reserva uma coluna vazia. Correções continuam recolhidas abaixo.

Verificação estática: Svelte check sem erros/avisos, autofixer nos dois componentes sem pendências, build passou e detector de layout sem achados. Os tokens existentes e temas foram preservados. Não houve screenshot, QA visual ou teste em celular físico; essa validação continua com João. Esta rodada só altera apresentação, sem mudanças na persistência/sincronização.

## Pacote local pronto para testar — 27/10/2026

Estado revalidado: HTTPS Tailscale :7000 aponta para 127.0.0.1:7777; o processo atual foi iniciado por `just serve`/`go run` neste checkout e usa o banco padrão em `~/Library/Application Support/whatnow/whatnow.db`. Não há listener em7790 e não existe `.demo/focus-preview.db` neste checkout. As instruções históricas do preview de25/10 não descrevem o serviço atual.

Preparados localmente (ignorados pelo Git): `.demo/review-20261027/whatnow` com frontend do commit af78747 e `.demo/review-20261027/before-update.db`, backup via SQLite do banco real, diretório700/arquivo600; quick_check ok e zero erros de chaves estrangeiras. O backup contém dados privados, não deve entrar em Git/PR. O banco original não foi alterado.

Executável iniciado para smoke em7791 com banco descartável, depois encerrado. GET raiz HTTP200 serviu o mesmo SHA256 do build: `0baa624cfb289e82c197b4f07dfde9dd23a2edf62e0bf217d3e7091c6e3f9da9`. GET /api/v1/execution retornou sessões/passos/estimativas e horário do servidor. Primeiro bind restrito foi negado; nova execução com permissão local passou. Nenhum teste visual. O processo principal segue intacto; seu HTML ainda é o build anterior. Trocar esse processo deve ser confirmado com João, pois agora é seu serviço principal, diferente do preview isolado anteriormente autorizado.

## Auditoria de entrega antes da troca do serviço

| Parte do fluxo | Evidência disponível | Estado |
| --- | --- | --- |
| Tarefa direta, CRUD/reordenação/divisão de passos, undo | Serviço único em steps.go/focus.go; TestStepCRUDReorderSplitUndo e TestCompleteStepWithoutStartingAndKeepActionTimeAcrossSegments; componente ExecutionSteps compartilhado | Implementado e coberto por testes lógicos; descoberta/gesto touch ainda depende de uso |
| Tempo persistente e conclusão explícita | TestFocusClockSwitchIdempotencyAndUndo verifica 90s preservados na troca, 105s após retomada, pause idempotente e undo sem retroceder tempo ou concluir pai; actionClock cobre excedente | Implementado; sessão independente rejeitada continua removida |
| Atualização entre dispositivos | live.svelte.ts invalida execução por SSE; Execution.startLive e teste de polling visível/desconectado; João relatou sincronização fluida na versão anterior | Base preservada; nova versão ainda não validada em aparelho |
| IA e continuidade do plano | TestWhatNowNextActionAndReplanPreserveExecution verifica mesmo foco ativo, 70000ms e passo concluído após replan; avaliação real fictícia documentada acima | Implementado e verificado nos cenários limitados |
| Interface e entrega em uso | Componentes/build verificados, pacote HTTP isolado confirmado, backup consistente pronto | Aceite visual aberto; atualização do serviço principal aguarda confirmação |

A entrega não está declarada concluída. O próximo ato dependente é substituir o processo que atende7777/HTTPS7000, autorizado anteriormente apenas para um preview distinto. Pergunta enviada a João; continuação automática não é consentimento. Nenhuma nova feature ou teste cosmético será criado só para preencher essa espera. Após autorização, refazer backup recente, trocar com possibilidade de voltar ao executável anterior, verificar saúde/HTML e preservação dos dados. O teste de uso mais útil é capturar uma intenção, entrar pelo E agora e acompanhar um passo no celular.

## Consolidação do editor — 27/10/2026

Após João confirmar que reiniciou e gostou das mudanças, autorizou continuar a UI enquanto testa. O editor de tarefa deixou de empilhar passos depois de todo o formulário: dados principais e passos ficam em duas colunas em desktop largo; celular usa diálogo em tela cheia com coluna única. Datas/dependências ficam recolhíveis (abertas quando já preenchidas), captura original recolhida, ações de salvar/cancelar e cabeçalho permanecem acessíveis durante rolagem. A estimativa agregada aparece junto aos passos. Os passos continuam usando ExecutionSteps e salvando imediatamente; o formulário da tarefa mantém confirmação explícita, com indicação dessa diferença.

EditSheet usa dialog.showModal, com foco modal e Escape nativos. Toasts são renderizados dentro do diálogo enquanto aberto para que erros e Undo não fiquem atrás da camada modal; fora dele continuam no layout raiz. Sem duplicação simultânea. Falha de salvar tarefa mantém o formulário e também informa erro nele. Não alterados IA, tempo, banco ou sincronização. Nenhuma biblioteca nova foi necessária para esta etapa.

Validação: 165 testes existentes passaram, Svelte check zero erros/avisos, build passou. Autofixer dos componentes modificados sem problemas; sugestões sobre efeitos preexistentes do layout foram inspecionadas: applyTheme escreve DOM/storage e welcomeDue é função pura, navegação continua efeito legítimo. Detector Impeccable sem achados. Não foi executado QA de browser/celular, conforme escolha de João. Servidor em uso não foi reiniciado. Próximo aceite útil: abrir uma tarefa com passos e editar/reordenar sem percorrer o formulário inteiro.

## Fundação shadcn-svelte — 27/10/2026

João instalou shadcn-svelte, Tailwind e a skill em `web/.agents`, depois autorizou
a migração inicial. Essa decisão substitui o adiamento anterior da biblioteca.
Skill exposta na raiz por symlink relativo, com instruções em ambos os níveis.
`docs/ui-migration.md` descreve tokens, componentes e a sequência incremental.

Button/Input oficiais adicionados pela CLI com Bun e adaptados nos próprios fontes
ao tema Whatnow. PlayToggle e InlineEdit agora os compõem. A ponte semântica em
layout.css usa os tokens existentes para cores, bordas, raios, fontes e sombras;
evita introduzir um tema neutro paralelo. Menus e editor completo são as próximas
fatias. Extras permanece referência para interações ausentes, sem dependência
adicionada nesta etapa. Persistência, IA e SSE não foram alterados.

Verificação: 167 testes existentes passaram, Svelte check sem erros/avisos e build
estático passou. Autofixer aprovou os componentes de produto; sugestões sobre
bind:this nos primitivos oficiais foram mantidas por integrarem a API pública de
ref. Detector encontrou apenas o gradiente preexistente do tema Vaporwave,
preservado por ser uma opção visual estabelecida. Sem QA de navegador ou aparelho.
