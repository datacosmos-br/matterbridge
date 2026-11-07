# Discord settings

> [!TIP]
> This page contains the details about discord settings. More general information about discord support in matterbridge can be found in [README.md](README.md).

## Server

Name or uid of your server/guild.

For example, If you want to bridge `https://discord.com/channels/AAAA/BBBB`, then this should be `AAAA`.

- Setting: **REQUIRED**
- Format: *string*
- Example:
  ```toml
  Server="yourservername"
  ```

## Token

Token to connect with Discord API. See [account.md](account.md) for
instructions to generate your token. If you want roles/groups mentions to be
shown with names instead of ID, you'll need to give your bot the "Manage Roles"
permission.

- Setting: **REQUIRED**
- Format: *string*
- Example:
  ```toml
  Token="YOUR_TOKEN_HERE"
  ```

## CustomStatus

Custom activity status message to set when connecting the bot to Discord. This
will be shown in member lists on Discord, as well as the bot's profile.

- Setting: **OPTIONAL**
- Format: *string*
- Example:
  ```toml
  CustomStatus="Bridging your chats"
  ```

## AllowMention

AllowMention controls which mentions are allowed.

If not specified, all mentions are allowed. Note that even when a mention is
not allowed, it will still be displayed nicely and be clickable. It just
prevents the ping/notification.

- Setting: **OPTIONAL**
- Format: *List[string]*
- Possible values:
  - `"everyone"` allows `@everyone` and `@here` mentions
  - `"roles"` allows `@role` mentions
  - `"users"` allows `@user` mentions
- Example:
  ```toml
  AllowMention=["everyone", "roles", "users"]
  ```

## ShowEmbeds

Shows title, description and URL of embedded messages (sent by other bots)

- Setting: **OPTIONAL**, **RELOADABLE**
- Format: *boolean*
- Example:
  ```toml
  ShowEmbeds=true
  ```

## EmbedFormat

Sets the format of the embeds attached to Discord messages sent by users (e.g. YouTube videos).
The only available placeholders are `{TITLE}`, `{URL}` (a URL), `{DESCRIPTION}`, `{TIME}`, `{FOOTER}`, `{IMAGE}` (a URL), `{THUMBNAIL}` (a URL), `{VIDEO}` (a URL), `{PROVIDER}` (a URL), `{AUTHOR}`, and `{AUTHORURL}` (a URL). If a parameter isn't in the originating embed it is simply omitted. (e.g. `Hello {AUTHOR}!` becomes `Hello !` if there is no author)

- Setting: **OPTIONAL** (default `" embed: {TITLE} - {DESCRIPTION} - {URL}"`), **RELOADABLE**
- Format: *string*
- Example:
  ```toml
  EmbedFormat="\n# {TITLE} ({URL}) by {AUTHOR} ({AUTHORURL})\n{DESCRIPTION}\n{TIME}\n{FOOTER}"
  ```

## UseUserName

Shows the username instead of the server nickname

- Setting: **OPTIONAL**, **RELOADABLE**
- Format: *boolean*
- Example:
  ```toml
  UseUserName=true
  ```

## UseDiscriminator

Show `#xxxx` discriminator with `UseUserName`

- Setting: **OPTIONAL**, **RELOADABLE**
- Format: *boolean*
- Example:
  ```toml
  UseDiscriminator=true
  ```

## WebhookURL

Specify WebhookURL. If given, will relay messages using the Webhook, which
gives a better look to messages.

If you have multiple discord channels, it is recommended to use the
`AutoWebHooks` setting. Alternatively, you'll have to specify the webhook URL
for each gateway and not in the account configuration.

- Setting: **OPTIONAL**
- Format: *string*
- Example:
  ```toml
  WebhookURL="Yourwebhooktokenhere"
  ```

## ForumID

Bridge a thread of a forum channel (a forum post). Set `channel` to the thread
ID and `ForumID` to the ID of the forum channel containing it. Webhooks belong
to the forum, so messages are sent and edited through the forum's webhook (from
`WebhookURL` or `AutoWebhooks`) inside the thread. Deletes use that webhook while
the message ID is still in the in-memory record of webhook-sent messages; after a
restart or eviction they go through the bot user, which needs `Manage Messages`. Receiving
messages from the thread works with the thread ID alone.

Do not add `ID:` at the beginning of the ForumID. A `WebhookURL` containing
`?thread_id=` is not supported: use the forum's webhook URL and `ForumID`.

This is a gateway channel option, set under `[gateway.inout.options]`.

- Setting: **OPTIONAL**
- Format: *string*
- Example:
  ```toml
  [[gateway.inout]]
  account="discord.game"
  channel="ID:1234567890"

    [gateway.inout.options]
    ForumID="987654321"
  ```

## AutoWebhooks

Relay messages using the Webhook, which gives a better look to messages. Needs "Manage Webhooks" permission to function.

AutoWebhooks automatically configures message sending in the style of puppets.
This is an easier alternative to manually configuring "WebhookURL" for each gateway,
as turning this on will automatically load or create webhooks for each channel.
This feature requires the "Manage Webhooks" permission (either globally or as per-channel).


Setting: OPTIONAL \
Format: boolean \
Example:

  ```toml
  AutoWebhooks=true
  ```

## QuoteDisable

Disable quotes in reply messages. Disable if your destination bridges understand native replies.

- Setting: **OPTIONAL**, **RELOADABLE**
- Format: *boolean*
- Example:
  ```toml
  QuoteDisable=true
  ```

## QuoteFormat

Format quoted/reply messages

- Setting: **OPTIONAL**, **RELOADABLE**
- Format: *string*
- Default: `"{MESSAGE} (re @{QUOTENICK}: {QUOTEMESSAGE})"`
- Example:
  ```toml
  QuoteFormat="{@{QUOTENICK}, {MESSAGE}}"
  ```

## QuoteLengthLimit

Limits {QUOTEMESSAGE} into integer characters specified in this.

TODO: does this cut graphemes or bytes?

- Setting: **OPTIONAL**
- Format: *int*
- Example:
  ```toml
  QuoteLengthLimit=46
  ```
