# Matterbridge ↔ Dify Integration

- Matterbridge publishes events to Dify via REST `/v1/workflows/run` with Bearer and `X-Idempotency-Key` = MessageID.
- Permission: phone→user table in matterbridge (gateway/identity) filters before calling Dify.
- Jira/Confluence ingested via Dify plugin datasource custom (Dify 1.9+).
- Retention: total auditável (v1), LGPD avançada fora de escopo.
