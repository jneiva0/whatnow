# Execução assistida — direção e diagnóstico

Atualizado em 27/10/2026. Fonte: feedback direto de João nesta conversa e leitura estática do código. Este documento organiza a próxima entrega; não declara implementação ou validação em uso.

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
