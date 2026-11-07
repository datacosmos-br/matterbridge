<<<<<<< HEAD
=======
> [!NOTE]
> This is a community-maintained fork [42wim/matterbridge](https://github.com/42wim/matterbridge), because
> there is not a lot of activity upstream and lots of bugfixes were not merged.
>
> **Get involved now and [become a maintainer for your favorite network/protocol](https://github.com/matterbridge-org/community/issues/1)!**

>>>>>>> 52bf0eeb559d1bef0e2ca19b8208f3e4867572a3
<div align="center">

# matterbridge

![Matterbridge Logo](img/matterbridge-notext.gif)<br />
**A simple chat bridge**<br />
<<<<<<< HEAD
Letting people be where they want to be.<br />
<sub>Bridges between a growing number of protocols. Click below to demo or join the development chat.</sub>

   <sup>

[Discord][mb-discord] |
[Gitter][mb-gitter] |
[IRC][mb-irc] |
[Keybase][mb-keybase] |
[Matrix][mb-matrix] |
[Mattermost][mb-mattermost] |
[MSTeams][mb-msteams] |
[Rocket.Chat][mb-rocketchat] |
[Slack][mb-slack] |
[Telegram][mb-telegram] |
[Twitch][mb-twitch] |
[WhatsApp][mb-whatsapp] |
[XMPP][mb-xmpp] |
[Zulip][mb-zulip] |
And more...
</sup>

---

[![Download stable](https://img.shields.io/github/release/42wim/matterbridge.svg?label=download%20stable)](https://github.com/42wim/matterbridge/releases/latest)
[![Maintainability](https://api.codeclimate.com/v1/badges/82dff70ef2ba85a6173a/maintainability)](https://codeclimate.com/github/42wim/matterbridge/maintainability)
[![Test Coverage](https://api.codeclimate.com/v1/badges/82dff70ef2ba85a6173a/test_coverage)](https://codeclimate.com/github/42wim/matterbridge/test_coverage)<br />

  <hr />
</div>
<div align="right"><sup>

**Note:** Matter<em>most</em> isn't required to run matter<em>bridge</em>.</sup></div>

<p>
  <a href="https://www.digitalocean.com/">
    <img src="https://opensource.nyc3.cdn.digitaloceanspaces.com/attribution/assets/PoweredByDO/DO_Powered_by_Badge_blue.svg" width="201px">
  </a>
</p>

## Table of Contents

- [matterbridge](#matterbridge)
- [Table of Contents](#table-of-contents)
  - [Features](#features)
    - [Natively supported](#natively-supported)
    - [3rd party via matterbridge api](#3rd-party-via-matterbridge-api)
    - [API](#api)
  - [Chat with us](#chat-with-us)
  - [Screenshots](#screenshots)
  - [Installing / upgrading](#installing--upgrading)
    - [Binaries](#binaries)
    - [Packages](#packages)
  - [Building](#building)
  - [Building with whatsapp (beta) multidevice support](#building-with-whatsapp-beta-multidevice-support)
  - [Configuration](#configuration)
    - [Basic configuration](#basic-configuration)
    - [Settings](#settings)
    - [Advanced configuration](#advanced-configuration)
      - [Examples](#examples)
        - [Bridge mattermost (off-topic) - irc (#testing)](#bridge-mattermost-off-topic---irc-testing)
        - [Bridge slack (#general) - discord (general)](#bridge-slack-general---discord-general)
    - [Running](#running)
      - [Docker](#docker)
      - [Systemd](#systemd)
    - [Changelog](#changelog)
    - [FAQ](#faq)
    - [Related projects](#related-projects)
    - [Articles / Tutorials](#articles--tutorials)
    - [Thanks](#thanks)

## Features

- [Support bridging between any protocols](https://github.com/42wim/matterbridge/wiki/Features#support-bridging-between-any-protocols)
- [Support multiple gateways(bridges) for your protocols](https://github.com/42wim/matterbridge/wiki/Features#support-multiple-gatewaysbridges-for-your-protocols)
- [Message edits and deletes](https://github.com/42wim/matterbridge/wiki/Features#message-edits-and-deletes)
- Preserves threading when possible
- [Attachment / files handling](https://github.com/42wim/matterbridge/wiki/Features#attachment--files-handling)
- [Username and avatar spoofing](https://github.com/42wim/matterbridge/wiki/Features#username-and-avatar-spoofing)
- [Private groups](https://github.com/42wim/matterbridge/wiki/Features#private-groups)
- [API](https://github.com/42wim/matterbridge/wiki/Features#api)

### Natively supported

- [Discord](https://discordapp.com)
- [Gitter](https://gitter.im)
- [Harmony](https://harmonyapp.io)
- [IRC](http://www.mirc.com/servers.html)
- [Keybase](https://keybase.io)
- [Matrix](https://matrix.org)
- [Mattermost](https://github.com/mattermost/mattermost-server/)
- [Microsoft Teams](https://teams.microsoft.com)
- [Mumble](https://www.mumble.info/)
- [Nextcloud Talk](https://nextcloud.com/talk/)
- [Rocket.chat](https://rocket.chat)
- [Slack](https://slack.com)
- [Ssh-chat](https://github.com/shazow/ssh-chat)
- ~~[Steam](https://store.steampowered.com/)~~
  - Not supported anymore, see [here](https://github.com/Philipp15b/go-steam/issues/94) for more info.
- [Telegram](https://telegram.org)
- [Twitch](https://twitch.tv)
- [VK](https://vk.com/)
- [WhatsApp](https://www.whatsapp.com/)
  - Whatsapp legacy is natively supported
  - Whatsapp multidevice beta is natively supported but you need to build yourself, see [here](#building-with-whatsapp-beta-multidevice-support)
- [XMPP](https://xmpp.org)
- [Zulip](https://zulipchat.com)

### 3rd party via matterbridge api

- [Delta Chat](https://github.com/deltachat-bot/matterdelta)
- [Minecraft](https://github.com/raws/mattercraft)
- [Minecraft](https://gitlab.com/Programie/MatterBukkit)

#### Past 3rd party projects

- [Discourse](https://github.com/DeclanHoare/matterbabble)
- [Facebook messenger](https://github.com/powerjungle/fbridge-asyncio)
- [Facebook messenger](https://github.com/VictorNine/fbridge)
- [Minecraft](https://github.com/elytra/MatterLink)
- [Reddit](https://github.com/bonehurtingjuice/mattereddit)
- [MatterAMXX](https://github.com/andrewlindberg/MatterAMXX): [Counter-Strike, half-life and more](https://forums.alliedmods.net/showthread.php?t=319430)
- [Vintage Story](https://github.com/NikkyAI/vs-matterbridge)
- [Ultima Online Emulator](https://github.com/kuoushi/ServUO-Matterbridge)
- [Teamspeak](https://github.com/Archeb/ts-matterbridge)

### API

The API is basic at the moment.
More info and examples on the [wiki](https://github.com/42wim/matterbridge/wiki/Api).

Used by the projects below. Feel free to make a PR to add your project to this list.

- [MatterLink](https://github.com/elytra/MatterLink) (Matterbridge link for Minecraft Forge server chat, archived)
- [MatterCraft](https://github.com/raws/mattercraft) (Matterbridge link for Minecraft Forge server chat)
- [MatterBukkit](https://gitlab.com/Programie/MatterBukkit) (Matterbridge link for Minecraft Bukkit/Spigot server chat)
- [pyCord](https://github.com/NikkyAI/pyCord) (crossplatform chatbot)
- [Mattereddit](https://github.com/bonehurtingjuice/mattereddit) (Reddit chat support)
- [fbridge-asyncio](https://github.com/powerjungle/fbridge-asyncio) (Facebook messenger support)
- [fbridge](https://github.com/VictorNine/fbridge) (Facebook messenger support)
- [matterbabble](https://github.com/DeclanHoare/matterbabble) (Discourse support)
- [MatterAMXX](https://forums.alliedmods.net/showthread.php?t=319430) (Counter-Strike, half-life and more via AMXX mod)
- [Vintage Story](https://github.com/NikkyAI/vs-matterbridge)
- [ServUO-matterbridge](https://github.com/kuoushi/ServUO-Matterbridge) (A matterbridge connector for ServUO servers)
- [ts-matterbridge](https://github.com/Archeb/ts-matterbridge) (Integrate teamspeak chat with matterbridge)
- [beerchat](https://github.com/mt-mods/beerchat) (Matterbridge link for minetest)

## Chat with us

Questions or want to test on your favorite platform? Join below:

- [Discord][mb-discord]
- [Gitter][mb-gitter]
- [IRC][mb-irc]
- [Keybase][mb-keybase]
- [Matrix][mb-matrix]
- [Mattermost][mb-mattermost]
- [Rocket.Chat][mb-rocketchat]
- [Slack][mb-slack]
- [Telegram][mb-telegram]
- [Twitch][mb-twitch]
- [XMPP][mb-xmpp] (<matterbridge@conference.jabber.de>)
- [Zulip][mb-zulip]

## Screenshots

See <https://github.com/42wim/matterbridge/wiki>

## Installing / upgrading

### Binaries

- Latest stable release [v1.26.0](https://github.com/42wim/matterbridge/releases/latest)
- Development releases (follows master) can be downloaded [here](https://github.com/42wim/matterbridge/actions) selecting the latest green build and then artifacts.

To install or upgrade just download the latest [binary](https://github.com/42wim/matterbridge/releases/latest). On \*nix platforms you may need to make the binary executable - you can do this by running `chmod a+x` on the binary (example: `chmod a+x matterbridge-1.24.1-linux-64bit`). After downloading (and making the binary executable, if necessary), follow the instructions on the [howto](https://github.com/42wim/matterbridge/wiki/How-to-create-your-config) for a step by step walkthrough for creating your configuration.

### Packages

- [Overview](https://repology.org/metapackage/matterbridge/versions)
- [snap](https://snapcraft.io/matterbridge)
- [scoop](https://github.com/42wim/scoop-bucket)

## Building

Most people just want to use binaries, you can find those [here](https://github.com/42wim/matterbridge/releases/latest)

If you really want to build from source, follow these instructions:
Go 1.18+ is required. Make sure you have [Go](https://golang.org/doc/install) properly installed.

Building the binary with **all** the bridges enabled needs about 3GB RAM to compile.
You can reduce this memory requirement to 0,5GB RAM by adding the `nomsteams` tag if you don't need/use the Microsoft Teams bridge.

Matterbridge can be build without gcc/c-compiler: If you're running on windows first run `set CGO_ENABLED=0` on other platforms you prepend `CGO_ENABLED=0` to the `go build` command. (eg `CGO_ENABLED=0 go install github.com/42wim/matterbridge`)

To install the latest stable run:

```bash
go install github.com/42wim/matterbridge
```

To install the latest dev run:

```bash
go install github.com/42wim/matterbridge@master
```

To install the latest stable run without msteams or zulip bridge:

```bash
go install -tags nomsteams,nozulip github.com/42wim/matterbridge
```

You should now have matterbridge binary in the ~/go/bin directory:

```bash
$ ls ~/go/bin/
matterbridge
```

## Configuration

### Basic configuration

See [howto](https://github.com/42wim/matterbridge/wiki/How-to-create-your-config) for a step by step walkthrough for creating your configuration.

### Settings

All possible [settings](https://github.com/42wim/matterbridge/wiki/Settings) for each bridge.

### Advanced configuration

- [matterbridge.toml.sample](https://github.com/42wim/matterbridge/blob/master/matterbridge.toml.sample) for documentation and an example.

### Examples

#### Bridge mattermost (off-topic) - irc (#testing)

```toml
[irc]
    [irc.libera]
    Server="irc.libera.chat:6667"
    Nick="yourbotname"

[mattermost]
    [mattermost.work]
    Server="yourmattermostserver.tld"
    Team="yourteam"
    Login="yourlogin"
    Password="yourpass"
    PrefixMessagesWithNick=true
    RemoteNickFormat="[{PROTOCOL}] <{NICK}> "

[[gateway]]
name="mygateway"
enable=true
    [[gateway.inout]]
    account="irc.libera"
    channel="#testing"

    [[gateway.inout]]
    account="mattermost.work"
    channel="off-topic"
```

#### Bridge slack (#general) - discord (general)

```toml
[slack]
[slack.test]
Token="yourslacktoken"
PrefixMessagesWithNick=true

[discord]
[discord.test]
Token="yourdiscordtoken"
Server="yourdiscordservername"

[general]
RemoteNickFormat="[{PROTOCOL}/{BRIDGE}] <{NICK}> "
=======
Letting people be where they want to be with the magic of [interoperability](https://en.wikipedia.org/wiki/Interoperability).<br />

---

[![Download stable](https://img.shields.io/github/release/42wim/matterbridge.svg?label=download%20stable)](https://github.com/42wim/matterbridge/releases/latest)

---

</div>

matterbridge is a solution to connect users on different platforms/protocols, allowing them to chat in the best conditions possible. Despite the name, Matter<em>most</em> isn't required to run matter<em>bridge</em>.

> [!WARNING]
> This fork has edited the history of the project to remove ~100MB of vendoring.
> This will faciliate review of new PRs, see [community/#5](https://github.com/matterbridge-org/community/issues/5)
> for the reasons why we did this and how to make sure we didn't introduce backdoors in this process.

## Features

Many features are available, but not all of them are supported on all protocols:

- [x] Bridge many rooms from supported protocols
- [x] Message threads/topics, and replies
- [x] Attachments, file uploads, and inline images
- [x] Transparent bridging with spoofed usernames and avatars
- [x] Private groups

**The complete and up-to-date list of supported protocols is in [docs/protocols/](docs/protocols/).** Additionally, we have an [API for 3rd party integration](https://github.com/42wim/matterbridge/wiki/Features#api) if you'd like to add a custom bridge without implementing it in this codebase.

![Screenshot of users discussing from different networks using matterbridge](https://user-images.githubusercontent.com/849975/52647227-9c3a5300-2ee4-11e9-9c57-ea096473aba8.png)

## Getting started with matterbridge

Get matterbridge up-and-running in a few minutes in 3 simple steps:

- [Setting up matterbridge](docs/setup.md) (see [docs/compiling.md](docs/compiling.md) for compiling from source)
- [Configuring matterbridge](docs/config.md)
- [Running matterbridge](docs/running.md) (CLI or as a systemd service)

## Documentation

See [docs/](docs/) folder in this repository.

## Contributing

You are welcome to submit pull requests, report bugs and request new features. matterbridge is a volunteer-run project and you are expected to behave with respect for the maintainers and other users. In particular, harassment and hate speech are not welcome.

For more development guidelines, see [docs/development/](docs/development/).

This project is licensed under the [GNU AGPLv3 license](LICENSE) since after commit `20988f6446c6ad3ea416044712e634d3ed85ee53`. It was relicensed following discussion in [commmunity/#10](https://github.com/matterbridge-org/community/issues/10). Apart from the obvious advantages of copyleft to promote innovation and cooperation, in very practical terms, we had to use either `GPL` or `AGPL` to include the `whatsappmulti` bridge in official builds and deprecate the broken legacy `whatsapp` bridge. When contributing to matterbridge development, you agree that your contributions will be published under that license.

Commits up-to `20988f6446c6ad3ea416044712e634d3ed85ee53` remain available under the looser [Apache License 2.0](LICENSE.old). 

### Chat with us

Questions or want to test on your favorite platform? Join us on:

- federated networks: [Jabber/XMPP][mb-xmpp], [Matrix][mb-matrix]
- non-free centralized networks: [Discord][mb-discord], [Keybase][mb-keybase], [Slack][mb-slack], [Telegram][mb-telegram], [Twitch][mb-twitch]
- self-hostable centralized networks: [IRC][mb-irc], [Mattermost][mb-mattermost], [Rocket.Chat][mb-rocketchat], [Zulip][mb-zulip]

## Related projects

- [jwflory/ansible-role-matterbridge](https://galaxy.ansible.com/jwflory/matterbridge) (Ansible role to simplify deploying Matterbridge)
- [matterbridge autoconfig](https://github.com/patcon/matterbridge-autoconfig)
- [matterbridge config viewer](https://github.com/patcon/matterbridge-heroku-viewer)
- [matterbridge-heroku](https://github.com/cadecairos/matterbridge-heroku)
- [mattermost-plugin](https://github.com/matterbridge/mattermost-plugin) - Run matterbridge as a plugin in mattermost
- [isla](https://github.com/alphachung/isla) (Bot for Discord-Telegram groups used alongside matterbridge)

## Thanks

Matterbridge wouldn't exist without amazing libraries, without [@42wim](https://github.com/42wim) who started the project, and without the 100+ contributors who participated in this adventure. See [docs/credits.md](docs/credits.md) for more complete credits.

<!-- Links -->

[mb-discord]: https://discord.gg/AkKPtrQ
[mb-gitter]: https://gitter.im/42wim/matterbridge
[mb-irc]: https://web.libera.chat/#matterbridge
[mb-keybase]: https://keybase.io/team/matterbridge
[mb-matrix]: https://riot.im/app/#/room/#matterbridge:matrix.org
[mb-mattermost]: https://framateam.org/signup_user_complete/?id=tfqm33ggop8x3qgu4boeieta6e
[mb-msteams]: https://teams.microsoft.com/join/hj92x75gd3y7
[mb-rocketchat]: https://open.rocket.chat/channel/matterbridge
[mb-slack]: https://join.slack.com/t/matterbridgechat/shared_invite/zt-2ourq2h2-7YvyYBq2WFGC~~zEzA68_Q
[mb-telegram]: https://t.me/Matterbridge
[mb-twitch]: https://www.twitch.tv/matterbridge
[mb-whatsapp]: https://www.whatsapp.com/
[mb-xmpp]: xmpp:matterbridge@conference.jabber.de?join
[mb-zulip]: https://matterbridge.zulipchat.com/register/
>>>>>>> 52bf0eeb559d1bef0e2ca19b8208f3e4867572a3

[[gateway]]
    name = "mygateway"
    enable=true

    [[gateway.inout]]
    account = "discord.test"
    channel="general"

    [[gateway.inout]]
    account ="slack.test"
    channel = "general"
```

## Running

See [howto](https://github.com/42wim/matterbridge/wiki/How-to-create-your-config) for a step by step walkthrough for creating your configuration.

```bash
Usage of ./matterbridge:
  -conf string
        config file (default "matterbridge.toml")
  -debug
        enable debug
  -gops
        enable gops agent
  -version
        show version
```

### Docker

Please take a look at the [Docker Wiki page](https://github.com/42wim/matterbridge/wiki/Deploy:-Docker) for more information.

### Systemd

Please take a look at the [Service Files page](https://github.com/42wim/matterbridge/wiki/Service-files) for more information.

## Changelog

See [changelog.md](https://github.com/42wim/matterbridge/blob/master/changelog.md)

## FAQ

See [FAQ](https://github.com/42wim/matterbridge/wiki/FAQ)

## Related projects

- [jwflory/ansible-role-matterbridge](https://galaxy.ansible.com/jwflory/matterbridge) (Ansible role to simplify deploying Matterbridge)
- [matterbridge autoconfig](https://github.com/patcon/matterbridge-autoconfig)
- [matterbridge config viewer](https://github.com/patcon/matterbridge-heroku-viewer)
- [matterbridge-heroku](https://github.com/cadecairos/matterbridge-heroku)
- [mattereddit](https://github.com/bonehurtingjuice/mattereddit)
- [matterlink](https://github.com/elytra/MatterLink)
- [mattermost-plugin](https://github.com/matterbridge/mattermost-plugin) - Run matterbridge as a plugin in mattermost
- [pyCord](https://github.com/NikkyAI/pyCord) (crossplatform chatbot)
- [fbridge](https://github.com/VictorNine/fbridge) (Facebook messenger support)
- [isla](https://github.com/alphachung/isla) (Bot for Discord-Telegram groups used alongside matterbridge)
- [matterbabble](https://github.com/DeclanHoare/matterbabble) (Connect Discourse threads to Matterbridge)
- [nextcloud talk](https://github.com/nextcloud/talk_matterbridge) (Integrates matterbridge in Nextcloud Talk)
- [mattercraft](https://github.com/raws/mattercraft) (Minecraft bridge)
- [vs-matterbridge](https://github.com/NikkyAI/vs-matterbridge) (Vintage Story bridge)
- [ServUO-matterbridge](https://github.com/kuoushi/ServUO-Matterbridge) (A matterbridge connector for ServUO servers)
- [ts-matterbridge](https://github.com/Archeb/ts-matterbridge) (Integrate teamspeak chat with matterbridge)

## Articles / Tutorials

- [matterbridge on kubernetes](https://medium.freecodecamp.org/using-kubernetes-to-deploy-a-chat-gateway-or-when-technology-works-like-its-supposed-to-a169a8cd69a3)
- <https://mattermost.com/blog/connect-irc-to-mattermost/>
- <https://blog.valvin.fr/2016/09/17/mattermost-et-un-channel-irc-cest-possible/>
- <https://blog.brightscout.com/top-10-mattermost-integrations/>
- <https://www.algoo.fr/blog/2018/01/19/recouvrez-votre-liberte-en-quittant-slack-pour-un-mattermost-auto-heberge/>
- <https://kopano.com/blog/matterbridge-bridging-mattermost-chat/>
- <https://www.stitcher.com/s/?eid=52382713>
- <https://daniele.tech/2019/02/how-to-use-matterbridge-to-connect-2-different-slack-workspaces/>
- <https://userlinux.net/mattermost-and-matterbridge.html>
- <https://nextcloud.com/blog/bridging-chat-services-in-talk/>
- <https://minecraftchest1.wordpress.com/2021/06/05/how-to-install-and-setup-matterbridge/>
- Youtube: [whatsapp - telegram bridging](https://www.youtube.com/watch?v=W-VXISoKtNc)

## Thanks

<p>This project is supported by:</p>
<p>
  <a href="https://www.digitalocean.com/">
    <img src="https://opensource.nyc3.cdn.digitaloceanspaces.com/attribution/assets/SVG/DO_Logo_horizontal_blue.svg" width="201px">
  </a>
</p>

Matterbridge wouldn't exist without these libraries:

- discord - <https://github.com/bwmarrin/discordgo>
- echo - <https://github.com/labstack/echo>
- gops - <https://github.com/google/gops>
- gozulipbot - <https://github.com/ifo/gozulipbot>
- gumble - <https://github.com/layeh/gumble>
- harmony - <https://github.com/harmony-development/shibshib>
- irc - <https://github.com/lrstanley/girc>
- keybase - <https://github.com/keybase/go-keybase-chat-bot>
- matrix - <https://github.com/matrix-org/gomatrix>
- mattermost - <https://github.com/mattermost/mattermost-server>
- msgraph.go - <https://github.com/yaegashi/msgraph.go>
- mumble - <https://github.com/layeh/gumble>
- nctalk - <https://github.com/gary-kim/go-nc-talk>
- rocketchat - <https://github.com/RocketChat/Rocket.Chat.Go.SDK>
- slack - <https://github.com/nlopes/slack>
- sshchat - <https://github.com/shazow/ssh-chat>
- steam - <https://github.com/Philipp15b/go-steam>
- telegram - <https://github.com/go-telegram-bot-api/telegram-bot-api>
- tengo - <https://github.com/d5/tengo>
- vk - <https://github.com/SevereCloud/vksdk>
- whatsapp - <https://github.com/Rhymen/go-whatsapp>
- whatsapp - <https://github.com/tulir/whatsmeow>
- xmpp - <https://github.com/mattn/go-xmpp>
- zulip - <https://github.com/ifo/gozulipbot>

<!-- Links -->

[mb-discord]: https://discord.gg/AkKPtrQ
[mb-gitter]: https://gitter.im/42wim/matterbridge
[mb-irc]: https://web.libera.chat/#matterbridge
[mb-keybase]: https://keybase.io/team/matterbridge
[mb-matrix]: https://riot.im/app/#/room/#matterbridge:matrix.org
[mb-mattermost]: https://framateam.org/signup_user_complete/?id=tfqm33ggop8x3qgu4boeieta6e
[mb-msteams]: https://teams.microsoft.com/join/hj92x75gd3y7
[mb-rocketchat]: https://open.rocket.chat/channel/matterbridge
[mb-slack]: https://join.slack.com/t/matterbridgechat/shared_invite/zt-2ourq2h2-7YvyYBq2WFGC~~zEzA68_Q
[mb-telegram]: https://t.me/Matterbridge
[mb-twitch]: https://www.twitch.tv/matterbridge
[mb-whatsapp]: https://www.whatsapp.com/
[mb-xmpp]: https://inverse.chat/
[mb-zulip]: https://matterbridge.zulipchat.com/register/
