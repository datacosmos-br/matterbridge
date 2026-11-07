# WhatsApp Cutover Runbook (big-bang)

1. Backup `.data/whatsapp/` SessionFile/DB.
2. Deploy new binary (whatsmeow unified, no build tag).
3. Validate QR login and one round-trip message per bridge.
4. Rollback: restore backup + previous binary.
