# Group summaries through Dify

Matterbridge can answer a summary command with a summary of the recent
conversation of a gateway. It does not run any language model itself: it
sends the transcript to a [Dify](https://dify.ai) workflow and relays the
workflow output. Only users linked to a corporate user may request summaries.

## Settings

All settings live in the `[general]` section.

| Setting | Default | Description |
| --- | --- | --- |
| `DifyURL` | empty | Base URL of the Dify API (for example `https://dify.example.com`). When empty, summaries are disabled: no command is intercepted and no transcript is kept. |
| `DifyAPIKey` | none | API key of the Dify workflow application, sent as `Authorization: Bearer <key>`. Required with `DifyURL`. It is never logged. |
| `IdentityStorePath` | none | Path of the sqlite database linking chat users to corporate users. Required with `DifyURL`. |
| `DifyWorkflowTimeout` | `30` | Seconds to wait for the workflow run. |
| `DifySummaryCommand` | `!summary` | Message text that requests a summary. |
| `DifySummaryMessages` | `50` | Number of recent messages of the gateway sent as transcript. |

Matterbridge refuses to start when `DifyURL` is not an absolute http(s) URL,
when `DifyAPIKey` or `IdentityStorePath` is missing, when a numeric setting is
negative, or when the identity database cannot be opened.

```toml
[general]
DifyURL="https://dify.example.com"
DifyAPIKey="app-xxxxxxxx"
IdentityStorePath="/var/lib/matterbridge/identities.sqlite"
```

## Linking users

A chat user is identified by the account they write from (`protocol.name` of
the bridge) and their user ID on it, for example a WhatsApp JID. Link them to
their corporate user with:

```sh
matterbridge -conf matterbridge.toml -identity-link 'whatsapp.corp:5511999990000@s.whatsapp.net=alice@corp'
```

The flag writes the link to `IdentityStorePath` (replacing any previous link
of that user) and exits without starting the gateways. A running matterbridge
sees new links immediately.

## Transcript

Every gateway keeps the last `DifySummaryMessages` text messages received on
its sending channels in memory, after opt-out redaction and message
modifications. The transcript is lost on restart.

## Command

When a user sends exactly `DifySummaryCommand` in a channel of a gateway, the
command is not relayed and:

- if the bridge does not report user IDs, or the user is not linked, they get
  a reply saying so and Dify is not called;
- otherwise matterbridge calls `POST <DifyURL>/v1/workflows/run` with

  ```json
  {
    "inputs": {
      "channel": "<gateway name>",
      "transcript": "<username>: <text>\n...",
      "requested_by": "<corporate user>"
    },
    "response_mode": "blocking",
    "user": "<corporate user>"
  }
  ```

  and the header `Idempotency-Key: <gateway name> <protocol> <message ID>`.
  Bridges that do not report message IDs (the API bridge) use a digest of the
  account, channel, user ID and receive time of the command instead.

The workflow must succeed (`data.status` is `succeeded`) and return a non-empty
`data.outputs.text`, which is sent as a reply in the channel of the command.
Any other answer, an HTTP error status or a timeout is logged and the user gets
a short failure reply.

The Dify workflow therefore needs the input variables `channel`, `transcript`
and `requested_by`, and an output variable `text`.
