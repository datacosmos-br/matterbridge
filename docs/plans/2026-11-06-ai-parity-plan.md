> Historical evidence (recovered 2026-11-07 from Kilo session ses_e4c45b055ffeYsn4s3WflJwUWV, kilo.db part prt_1b3cf6d66001sOZC4HFp5T78Ot; originally written 2026-11-06 to .kilo/plans/1794016694186-matterbridge-ai-parity-plan.md and lost to a git clean).
> Superseded as execution authority by 2026-11-07-upstream-adoption.md; kept verbatim for its product decisions.

# Matterbridge — Modernização WhatsApp + Paridade de Bridges + IA (Dify) — Plano Implementation-Ready

## Estado atual (evidência)
- Workspace `/home/marlonsc/gpt/matterbridge` é fork datacosmos-br — `datacosmos 4ac0248d` 23 commits à frente de `master d7a5c142`, inclui merges de `origin/master`, `remoterelease`, `whatsappmulti`.
- Branches locais divergentes: `mspgeek 27c1045d` (197 commits, reescrita de import path, unificação WhatsApp), `master d7a5c142` (14 commits atrás de `origin/master 25c62e28`), `origin/master` puro 42cbf5d2. `datacosmos` não é upstream puro.
- WhatsApp hoje: `bridge/whatsapp/` (legacy Rhymen/go-whatsapp, deprecado) + `bridge/whatsappmulti/` (whatsmeow, build tag `whatsappmulti`). `bridge_config.go` decide via tag. CI não builda tag; `Dockerfile` builda sem tag, `Dockerfile_whatsappmulti` com tag. Sem testes WhatsApp.
- Bridges: 19 dirs em `bridge/`; paridade WhatsApp acima de Discord/Slack. Razão `bridge/*/whatsapp.go` vs `bridge/*/handlers.go` mostra WhatsApp como outlier.

## Decisões travadas
- **Base**: `datacosmos` como fonte; `origin/master` + `mspgeek` via cherry-pick auditado; descartar reescrita de import path do mspgeek.
- **WhatsApp**: unificar em whatsmeow único em `bridge/whatsapp`, sem build tag (igual mspgeek); migrar `SessionFile` SQLite (`.data/whatsapp/`); alinhar `bridge_config.go`, `bridge.go`, Dockerfile(s) e CI para build/testar sempre whatsmeow; remover legacy.
- **Gaps prioritários WhatsApp**: LID→PN, reactions/poll/sticker/status, corrigir avatar/mention; broadcast/newsletter como follow-up.
- **Paridade**: todos os bridges em baseline AI/search (texto, reply/quote, arquivo, delete/edit, threads, captions, reactions, presença) — Discord, Slack, Teams, Meet, Mattermost, Telegram, Matrix, IRC, XMPP e demais. Manutenção prioriza Discord/Slack/WhatsApp/Teams/Meet.
- **Teams**: chat + meetings (Graph/webhooks); **Meet (Google)**: spike de chat via Workspace API com fallback obrigatório para transcrição/resumo pós-reunião para Dify (API Meet não expõe chat geral).
- **IA / Dify (pesquisa Datacosmos)**: Dify como sidecar via `datacosmos-br/dify-helm` / `cosmos-charts/dify 0.36.1` (api/worker/sandbox, redis, postgres, object storage, vector DB). Matterbridge não embute LLM; publica eventos e impõe permissão; Dify faz RAG/workflows/agentes. Knowledge via plugin datasource custom Dify 1.9+ para Jira/Confluence. Documentação Hub como sistema vivo integrado ao Dify (padrão DataBoar/docs PLANS_HUB como índice, mas operação viva).
- **Identidade/permissão**: tabela phone→user (CRUD + sync LDAP/IdP) para vínculo WhatsApp↔corporativo; segredo Jira/Confluence como conta de serviço única (K8s Secret) com pré-filtro por usuário no matterbridge (ACL efetiva por usuário, credencial única). Retenção total auditável v1; LGPD avançada e modo offline fora de escopo v1.
- **Contrato Matterbridge→Dify**: REST API (Workflow/Chatflow + Knowledge Base) com Bearer, idempotência por MessageID; Dify expõe MCP depois se precisar. Jira/Confluence via plugin datasource custom.
- **Contexto grupos MS**: resumo + triagem (resumo diário/semanal do grupo, intenção/urgência, roteamento para Mattermost/Slack) antes de moderação/automação.
- **Deployment**: K8s Helm (matterbridge + Dify via datacosmos dify-helm/cosmos-charts; values para Jira/Confluence/Mattermost/Slack/Discord).
- **Rollout**: big-bang direto no próximo deploy, com plano de rollback por backup de SessionFile/DB (shadow desconsiderado por decisão).
- **Validação**: testes unitários por bridge + build/test com `-tags whatsappmulti` (enquanto tag existir) e depois build único whatsmeow no CI; harness golden de mensagens.
- **Stack LLM/vector**: interface desacoplada (adaptador) — escolha concreta após pesquisa rápida; default de trabalho: API OpenAI-compatível + pgvector/Qdrant em sidecar via Dify.

## Plano de execução (ordem)
1. **Fase 0 — Cherry-pick auditado**: rebase `datacosmos` sobre `origin/master`; cherry-pick de `mspgeek` (cache persistente, gateway de comandos, fixes whatsmeow) sem import path; fechar `pr/*` (upload, threads, captions, socket-mode, matrix/xmpp/irc, API proxy) com testes.
2. **Fase 1 — Unificação WhatsApp**: mover `whatsappmulti` para `bridge/whatsapp`, remover legacy, atualizar `bridge_config.go`/`bridge.go`/Dockerfile/CI; migração SessionFile; LID→PN, reactions/poll/sticker/status; harness golden WhatsApp↔Mattermost/Slack/Discord.
3. **Fase 2 — Paridade bridges**: aplicar matriz de paridade (texto/reply/arquivo/delete/edit/threads/captions/reactions) para Discord, Slack, Teams (chat+meetings), Telegram, Matrix, IRC, XMPP; spike Meet + fallback transcrição; Documentation Hub vivo + índice de planos.
4. **Fase 3 — IA Dify**: deploy K8s Helm Dify (datacosmos helm); contrato REST idempotente; tabela phone→user; plugin datasource custom Jira/Confluence; workflows de resumo/triagem; retenção total auditável; admissão de que Google Meet chat via API não existe como bridging geral (fallback transcrição).
5. **Fase 4 — Validação & big-bang**: `golangci-lint`/CI com tag/único, smoke Mattermost/Slack/Discord/Teams/Meet, backup SessionFile/DB, cutover big-bang com rollback documentado.

## Fora de escopo v1
- LGPD avançada, modo offline, conta bot por usuário (optou-se por conta de serviço única), broadcast/newsletter completo, moderação/automação completa.

## Riscos
- Conflito `mspgeek` vs `datacosmos` em spins (remover vs manter legacy) — mitigado por cherry-pick sem reescrita.
- `modernc.org/sqlite` pinado (ver #2797) — ver migração SQLite WhatsApp.
- Google Meet sem chat geral — fallback transcrição obrigatório.
- LLM/vector a definir — começar desacoplado via adaptador.

