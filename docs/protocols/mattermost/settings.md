# Mattermost settings

> [!TIP]
> This page contains the details about mattermost settings. More general information about mattermost support in matterbridge can be found in [README.md](README.md).

## Login

login of your bot.
Use a dedicated user for this and not your own!

- Setting: **REQUIRED** (when not using webhooks)
- Format: *string*
- Example:
<<<<<<< HEAD

=======
>>>>>>> 52bf0eeb559d1bef0e2ca19b8208f3e4867572a3
  ```toml
  Login="yourlogin"
  ```

## NoTLS

Enable this to make a http connection (instead of https) to your mattermost.

- Setting: **OPTIONAL**
- Format: *boolean*
- Example:
<<<<<<< HEAD

=======
>>>>>>> 52bf0eeb559d1bef0e2ca19b8208f3e4867572a3
  ```toml
  NoTLS=true
  ```

## Password

login of your bot.
Use a dedicated user for this and not your own!

- Setting: **REQUIRED** (when not using webhooks)
- Format: *string*
- Example:
<<<<<<< HEAD

=======
>>>>>>> 52bf0eeb559d1bef0e2ca19b8208f3e4867572a3
  ```toml
  Password="yourpass"
  ```

## Server

The mattermost hostname. (do not prefix it with http or https)

- Setting: **REQUIRED**
- Format: *string*
- Example:
<<<<<<< HEAD

=======
>>>>>>> 52bf0eeb559d1bef0e2ca19b8208f3e4867572a3
  ```toml
  Server="yourmattermostserver.domain"
  ```

## SkipVersionCheck

Skip the Mattermost server version checks that are normally done when connecting.
<<<<<<< HEAD
The usage scenario for this feature would be when the Mattermost instance is hosted behind a
=======
The usage scenario for this feature would be when the Mattermost instance is hosted behind a 
>>>>>>> 52bf0eeb559d1bef0e2ca19b8208f3e4867572a3
reverse proxy that suppresses "non-standard" response headers in flight.

- Setting: **OPTIONAL**
- Format: *boolean*
- Example:
<<<<<<< HEAD

=======
>>>>>>> 52bf0eeb559d1bef0e2ca19b8208f3e4867572a3
  ```toml
  SkipVersionCheck=true
  ```

## Team

Your team on mattermost.

- Setting: **REQUIRED** (when not using webhooks)
- Format: *string*
- Example:
<<<<<<< HEAD

=======
>>>>>>> 52bf0eeb559d1bef0e2ca19b8208f3e4867572a3
  ```toml
  Team="yourteam"
  ```

## Token

personal access token of the bot.
<<<<<<< HEAD
new feature since mattermost 4.1. See <https://docs.mattermost.com/developer/personal-access-tokens.html>
=======
new feature since mattermost 4.1. See https://docs.mattermost.com/developer/personal-access-tokens.html
>>>>>>> 52bf0eeb559d1bef0e2ca19b8208f3e4867572a3
you can use token instead of login/password.

- Setting: **OPTIONAL** (when not using webhooks)
- Format: *string*
- Example:
<<<<<<< HEAD

=======
>>>>>>> 52bf0eeb559d1bef0e2ca19b8208f3e4867572a3
  ```toml
  Token="abcdefghijklm"
  ```

### WebhookURL

> [!WARNING]
> NOT RECOMMENDED TO USE INCOMING/OUTGOING WEBHOOK.
> USE DEDICATED BOT USER WHEN POSSIBLE!

Url is your incoming webhook url as specified in mattermost.
See account settings - integrations - incoming webhooks on mattermost.
If specified, messages will be sent to mattermost using this URL

- Setting: **OPTIONAL**
- Format: *string*
<<<<<<< HEAD
- Example:

=======
- Example: 
>>>>>>> 52bf0eeb559d1bef0e2ca19b8208f3e4867572a3
  ```toml
  WebhookURL="https://yourdomain/hooks/yourhookkey"`
  ```

### WebhookBindAddress

> [!WARNING]
> NOT RECOMMENDED TO USE INCOMING/OUTGOING WEBHOOK.
> USE DEDICATED BOT USER WHEN POSSIBLE!

Address to listen on for outgoing webhook requests from mattermost.
See account settings - integrations - outgoing webhooks on mattermost.
If specified, messages will be received from mattermost on this ip:port
(this will only work if `WebhookURL` is also configured)

- Setting: **OPTIONAL**
- Format: *string*
- Example:
<<<<<<< HEAD

=======
>>>>>>> 52bf0eeb559d1bef0e2ca19b8208f3e4867572a3
  ```toml
  WebhookBindAddress="0.0.0.0:9999"
  ```

### IconURL

> [!WARNING]
> NOT RECOMMENDED TO USE INCOMING/OUTGOING WEBHOOK.
> USE DEDICATED BOT USER WHEN POSSIBLE!

Icon that will be showed in mattermost.
This only works when `WebhookURL` is configured.

- Setting: **OPTIONAL**
- Format: *string*
- Example:
<<<<<<< HEAD

=======
>>>>>>> 52bf0eeb559d1bef0e2ca19b8208f3e4867572a3
  ```toml
  IconURL="http://youricon.png"
  ```
