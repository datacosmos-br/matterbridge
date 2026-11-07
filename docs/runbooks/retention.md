# Retention Policy (v1)

Target policy: full, auditable retention of messages and events from bridged WhatsApp groups and
other bridges, with no TTL and exports on demand. Advanced LGPD handling is out of scope for v1.

Current state: matterbridge does **not** archive message content and offers no export. With
`PersistentMessageStorePath` it only keeps the gateway message-ID map (bounded by
`PersistentMessageStoreSize`) so edits, deletes and replies survive restarts. Retention and
exports must come from the destination platforms or a dedicated archival system until such a
subsystem exists; do not rely on matterbridge for compliance.
