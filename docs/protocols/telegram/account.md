TODO: is this information up to date?

See [here](https://core.telegram.org/bots#6-botfather) and [here](https://www.linkedin.com/pulse/telegram-bots-beginners-marco-frau).

**Sequence within a BotFather chat:**

> **You**: /setprivacy
> 
> **BotFather**: Choose a bot to change group messages settings.
> 
> **You**: @your_name_bot
> 
> **BotFather**: 'Enable' - your bot will only receive messages that either start with the '/' symbol or mention the bot by username.
>
> - 'Disable' - your bot will receive all messages that people send to groups.
> - Current status is: ENABLED
>
> **You**: Disable
> 
> **BotFather**: Success! The new status is: DISABLED. /help

<<<<<<< HEAD
(copied from [this stack overflow answer](https://stackoverflow.com/a/45441585))

## Tutorial

This basic example will guide you through getting started on your first Telegram bridge.

- Create your bot and retrieve your API token.
=======

(copied from [this stack overflow answer](https://stackoverflow.com/a/45441585))


## Tutorial

This basic example will guide you through getting started on your first Telegram bridge. 

- Create your bot and retrieve your API token. 
>>>>>>> 52bf0eeb559d1bef0e2ca19b8208f3e4867572a3
- Make sure to disable privacy mode.
- Add the bot to your Telegram group
- Set up the Matterbridge configuration.

<<<<<<< HEAD
### 1. Speak to the BotFather.
=======
### 1. Speak to the BotFather. 
>>>>>>> 52bf0eeb559d1bef0e2ca19b8208f3e4867572a3

Message the [BotFather](https://t.me/BotFather) bot and run the `/newbot` command to get started. Once you pick a name it will provide you with a `token` used to access the HTTP API.

**Disable bot privacy**

> **You**: /setprivacy
> 
> **BotFather**: Choose a bot to change group messages settings.
> 
> **You**: @your_name_bot
> 
> **BotFather**: 'Enable' - your bot will only receive messages that either start with the '/' symbol or mention the bot by username.
>
> - 'Disable' - your bot will receive all messages that people send to groups.
> - Current status is: ENABLED
>
> **You**: Disable
> 
> **BotFather**: Success! The new status is: DISABLED. /help

<<<<<<< HEAD
Now add this newly created bot to the Telegram chat you're trying to bridge.

### 2. Add your config to Matterbridge.
=======
Now add this newly created bot to the Telegram chat you're trying to bridge. 

### 2. Add your config to Matterbridge. 
>>>>>>> 52bf0eeb559d1bef0e2ca19b8208f3e4867572a3

Here is an example gateway bridging `Discord`<>`Telegram` which should be placed in your `matterbridge.toml` (or equivalent)

```toml
# / 
# SERVERS
# / 
[telegram]
    [telegram.mytelegram]
        Token="your_token_from_botfather" 
        RemoteNickFormat="<{NICK}> " # How your message will be formatted when bridged. 
        MessageFormat="HTMLNick :"
        QuoteFormat="{MESSAGE} (re @{QUOTENICK}: {QUOTEMESSAGE})"
        QuoteLengthLimit=46 # Truncuate long quotes to prevent spammy bridged messages
        IgnoreMessages="^/" # Don't bridge bot commands (as the responses will not be bridged)

# / 
# GATEWAYS
# / 
[[gateway]]
name="YourUniqueGateWayName"
enable=true

    [[gateway.inout]]
    account="telegram.mytelegram"
    channel="-100xxx"

    [[gateway.inout]]
    account="discord.mydiscord"
    channel="channel-name" 
```

> See all possible settings [here](https://github.com/42wim/matterbridge/wiki/Settings#telegram)

<<<<<<< HEAD
The easiest way to retrieve your channel id is to navigate to <https://web.telegram.org/> and simply copy the numbers at the end of the hyperlink into the `xxx` position in the example snippet above.

To add more nodes to this bridge, you simply need to add additional `[[gateway.inout]]` fields.

Make sure to test your config with the `-debug` flag after each change.

### Limitations

- The Telegram API does not report any changes **when messages are deleted** so Matterbridge is unable to remove any bridged messages after they've been sent (This will render many common spam solutions useless). Use regexp with the `IgnoreMessages=` field to remove any common spam messages. # Telegram

## Matterbridge is not relaying messages from Telegram

See

- <https://core.telegram.org/bots#privacy-mode>
- <https://github.com/yagop/node-telegram-bot-api/issues/174#issuecomment-244632667>
=======
The easiest way to retrieve your channel id is to navigate to https://web.telegram.org/ and simply copy the numbers at the end of the hyperlink into the `xxx` position in the example snippet above. 


To add more nodes to this bridge, you simply need to add additional `[[gateway.inout]]` fields.

Make sure to test your config with the `-debug` flag after each change. 

### Limitations

- The Telegram API does not report any changes **when messages are deleted** so Matterbridge is unable to remove any bridged messages after they've been sent (This will render many common spam solutions useless). Use regexp with the `IgnoreMessages=` field to remove any common spam messages. # Telegram 

## Matterbridge is not relaying messages from Telegram

See 
* https://core.telegram.org/bots#privacy-mode
* https://github.com/yagop/node-telegram-bot-api/issues/174#issuecomment-244632667
>>>>>>> 52bf0eeb559d1bef0e2ca19b8208f3e4867572a3

If your bot is not getting messages:

Disable privacy mode with @Botfather.

<<<<<<< HEAD
- Go to [BotFather](https://t.me/botfather) send /setprivacy.
- Select the username of the bot.
- Select Disable.
- Kick bot from chat if it's already in it.
- Invite bot to chat.
=======
* Go to [BotFather](https://t.me/botfather) send /setprivacy.
* Select the username of the bot.
* Select Disable.
* Kick bot from chat if it's already in it.
* Invite bot to chat.
>>>>>>> 52bf0eeb559d1bef0e2ca19b8208f3e4867572a3

The order is important.

## Matterbridge is not relay messages from other bots

> Bots talking to each other could potentially get stuck in unwelcome loops. To avoid this, we decided that bots will not be able to see messages from other bots regardless of mode.

<<<<<<< HEAD
<https://core.telegram.org/bots/faq#why-doesn-39t-my-bot-see-messages-from-other-bots>
=======
https://core.telegram.org/bots/faq#why-doesn-39t-my-bot-see-messages-from-other-bots
>>>>>>> 52bf0eeb559d1bef0e2ca19b8208f3e4867572a3

## Matterbridge is not relaying images/stickers/files from Telegram

Because images/stickers/files are from non-public url's, you'll need to setup a [mediaserver](Mediaserver-setup-(advanced))

## Matterbridge is not relaying messages to Telegram

Did you enable `MessageFormat="HTML"` in your config?  
You could be sending invalid HTML. Set it to `MessageFormat=""`

<<<<<<< HEAD
More info in <https://github.com/42wim/matterbridge/blob/master/matterbridge.toml.sample#L836-L838>
=======
More info in https://github.com/42wim/matterbridge/blob/master/matterbridge.toml.sample#L836-L838
>>>>>>> 52bf0eeb559d1bef0e2ca19b8208f3e4867572a3

## Matterbridge is not deleting messages from Telegram to other bridges

Telegram doesn't has "deleted messages" metadata, so we don't know which messages are deleted.

## Users are shown as "unknown user"

Telegram channels will always return unknown users. Telegram groups will show usernames if possible

## How to use MediaConvertTgs with lottie?

This requires the external dependency `lottie`, which can be installed like this:
`pip install lottie cairosvg`
<<<<<<< HEAD
<https://github.com/42wim/matterbridge/issues/874>

Note that if you insist on using an ancient Python version like 3.5, the pip installation is slightly more complicated. Matterbridge expects `lottie_convert.py` to be in your `$PATH`; if that's not already the case, try putting this into your `~/.profile`:

```text
=======
https://github.com/42wim/matterbridge/issues/874

Note that if you insist on using an ancient Python version like 3.5, the pip installation is slightly more complicated. Matterbridge expects `lottie_convert.py` to be in your `$PATH`; if that's not already the case, try putting this into your `~/.profile`:
```
>>>>>>> 52bf0eeb559d1bef0e2ca19b8208f3e4867572a3
PATH=$HOME/.local/bin:$PATH
export PATH
```

If you encounter bugs with this, try to extract the Telegram sticker file and run lottie on it like this:
`lottie_convert.py --input-format lottie file_1234_tgs.webp myoutput.webp`
This might give you additional information about what's going on.
<<<<<<< HEAD
=======

>>>>>>> 52bf0eeb559d1bef0e2ca19b8208f3e4867572a3
