# Telegram settings

> [!TIP]
> This page contains the details about telegram settings. More general information about telegram support in matterbridge can be found in [README.md](README.md).

## MediaConvertTgs

Convert Tgs (Telegram animated sticker) images to some other file format before upload. See FAQ for setup instructions.

- Setting: **OPTIONAL**, **RELOADABLE**
- Format: *string*
- Possible values:
<<<<<<< HEAD
  - `"png"` (still-image)
  - `"webp"` (animated webp)
- Example:

=======
    - `"png"` (still-image)
    - `"webp"` (animated webp)
- Example:
>>>>>>> 52bf0eeb559d1bef0e2ca19b8208f3e4867572a3
  ```toml
  MediaConvertTgs="png"
  ```

## MediaConvertWebPToPNG

Convert WebP images to PNG before upload to the media server or other bridges (for compatibility).

- Setting: **OPTIONAL**, **RELOADABLE**
- Format: *boolean*
- Example:
<<<<<<< HEAD

=======
>>>>>>> 52bf0eeb559d1bef0e2ca19b8208f3e4867572a3
  ```toml
  MediaConvertWebPToPNG=true
  ```

## MessageFormat

TODO: is this only for outgoing messages?

- Setting: **OPTIONAL**, **RELOADABLE**
- Format: *string*
- Possible values:
  - [`"HTML"`](https://core.telegram.org/bots/api#html-style)
  - [`"Markdown"`](https://core.telegram.org/bots/api#markdown-style). Deprecated, does not display links with underscores `_` correctly.
  - [`"MarkdownV2"`](https://core.telegram.org/bots/api#markdownv2-style)
  - `"HTMLNick"`. This only allows HTML for the nick, the message itself will be html-escaped.
- Default: `""`
- Example:
<<<<<<< HEAD

=======
>>>>>>> 52bf0eeb559d1bef0e2ca19b8208f3e4867572a3
  ```toml
  MessageFormat="MarkdownV2"
  ```

## QuoteDisable

Disable quoted/reply messages

- Setting: **OPTIONAL**, **RELOADABLE**
- Format: *boolean*
- Example:
<<<<<<< HEAD

=======
>>>>>>> 52bf0eeb559d1bef0e2ca19b8208f3e4867572a3
  ```toml
  QuoteDisable=true
  ```

## QuoteFormat

Format quoted/reply messages

- Setting: **OPTIONAL**, **RELOADABLE**
- Format: *string*
- Default: `"{MESSAGE} (re @{QUOTENICK}: {QUOTEMESSAGE})"`
<<<<<<< HEAD
- Example:

=======
- Example: 
>>>>>>> 52bf0eeb559d1bef0e2ca19b8208f3e4867572a3
  ```toml
  QuoteFormat="{@{QUOTENICK}, {MESSAGE}}"
  ```

## QuoteLengthLimit

Limits {QUOTEMESSAGE} into integer characters specified in this.

TODO: does this cut graphemes or bytes?

- Setting: **OPTIONAL**
- Format: *int*
- Example:
<<<<<<< HEAD

=======
>>>>>>> 52bf0eeb559d1bef0e2ca19b8208f3e4867572a3
  ```toml
  QuoteLengthLimit=46
  ```

## DisableWebPagePreview

Disables link previews for links in messages

- Setting: **OPTIONAL**, **RELOADABLE**
- Format: *boolean*
- Example:
<<<<<<< HEAD

=======
>>>>>>> 52bf0eeb559d1bef0e2ca19b8208f3e4867572a3
  ```toml
  DisableWebPagePreview=true
  ```

## Token

Token to connect with telegram API.

- Setting: **REQUIRED**
- Format: *string*
- Example:
<<<<<<< HEAD

=======
>>>>>>> 52bf0eeb559d1bef0e2ca19b8208f3e4867572a3
  ```toml
  Token="Yourtokenhere"
  ```

## UseFirstName

If enabled use the "First Name" as username. If this is empty use the Username
If disabled use the "Username" as username. If this is empty use the First Name
If all names are empty, username will be "unknown"

- Setting: **OPTIONAL**, **RELOADABLE**
- Format: *boolean*
- Example:
<<<<<<< HEAD

=======
>>>>>>> 52bf0eeb559d1bef0e2ca19b8208f3e4867572a3
  ```toml
  UseFirstName=true
  ```

## UseInsecureURL

> [!WARNING]
> If enabled this will relay GIF/stickers/documents and other attachments as URLs
> Those URLs will contain your bot-token. This may not be what you want.
> For now there is no secure way to relay GIF/stickers/documents without seeing your token.

- Setting: **OPTIONAL**, **RELOADABLE**
- Format: *boolean*
- Example:
<<<<<<< HEAD

=======
>>>>>>> 52bf0eeb559d1bef0e2ca19b8208f3e4867572a3
  ```toml
  UseInsecureURL=true
  ```
