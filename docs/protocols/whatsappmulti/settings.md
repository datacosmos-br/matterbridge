# WhatsApp settings

> [!TIP]
> This page contains the details about whatsapp settings. More general information about whatsapp support in matterbridge can be found in [README.md](README.md).

## Number

Number you will use as a relay bot. Tip: Get some disposable sim card, don't rely on your own number.

- Setting: **REQUIRED**
- Format: *string*
- Example:
<<<<<<< HEAD

=======
>>>>>>> 52bf0eeb559d1bef0e2ca19b8208f3e4867572a3
  ```toml
  Number="+48111222333"
  ```

## SessionFile

First time that you login you will need to scan QR code, then credentials will
be saved in a session file. If you do not set a `SessionFile`, you will need
to scan your QR code on every restart.

Unless this option is set, the Whatsapp login session is stored only in memory
until restarting matterbridge.

- Setting: **OPTIONAL**
- Format: *string*
- Example:
<<<<<<< HEAD

=======
>>>>>>> 52bf0eeb559d1bef0e2ca19b8208f3e4867572a3
  ```toml
  SessionFile="/usr/share/matterbridge/session-48111222333.gob"
  ```
