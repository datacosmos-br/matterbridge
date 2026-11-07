# IRC settings

> [!TIP]
> This page contains the details about IRC settings. More general information about IRC support in matterbridge can be found in [README.md](README.md).

## Charset

If you know your charset, you can specify it manually. Otherwise it tries to detect this automatically.
<<<<<<< HEAD

=======
 
>>>>>>> 52bf0eeb559d1bef0e2ca19b8208f3e4867572a3
The selected charset will be converted to utf-8 when sent to other bridges.

- Setting: **OPTIONAL**
- Format: *string*
- Example:
<<<<<<< HEAD

=======
>>>>>>> 52bf0eeb559d1bef0e2ca19b8208f3e4867572a3
  ```toml
  Charset="utf-8"
  ```

## ColorNicks

ColorNicks will show each nickname in a different color.
Only works in IRC right now.

- Setting: **OPTIONAL**, **RELOADABLE**
- Format: *boolean*
- Example:
<<<<<<< HEAD

=======
>>>>>>> 52bf0eeb559d1bef0e2ca19b8208f3e4867572a3
  ```toml
  ColorNicks=true
  ```

## DebugLevel

Debug log verbosity.

- Setting: **OPTIONAL**
- Default: *0*
- Format: *int*
- Example:
<<<<<<< HEAD

=======
>>>>>>> 52bf0eeb559d1bef0e2ca19b8208f3e4867572a3
  ```toml
  DebugLevel=1
  ```

## JoinDelay

Delay in milliseconds between channel joins.
Only useful when you have a LOT of channels to join.

- Setting: **OPTIONAL**, **RELOADABLE**
- Default: *0*
- Format: *int*
- Example:
<<<<<<< HEAD

=======
>>>>>>> 52bf0eeb559d1bef0e2ca19b8208f3e4867572a3
  ```toml
  JoinDelay=1000
  ```

## MessageDelay

Flood control.
Delay in milliseconds between each message send to the IRC server.

- Setting: **OPTIONAL**, **RELOADABLE**
- Default: *1300*
- Format: *int*
- Example:
<<<<<<< HEAD

=======
>>>>>>> 52bf0eeb559d1bef0e2ca19b8208f3e4867572a3
  ```toml
  MessageDelay=1300
  ```

## MessageLength

Maximum length of message sent to irc server. If it exceeds
`<message clipped>` will be add to the message.

- Setting: **OPTIONAL**, **RELOADABLE**
- Default: *400*
- Format: *int*
- Example:
<<<<<<< HEAD

=======
>>>>>>> 52bf0eeb559d1bef0e2ca19b8208f3e4867572a3
  ```toml
  MessageLength=400
  ```

## MessageQueue

<<<<<<< HEAD
Maximum amount of messages to hold in queue. If queue is full
messages will be dropped.
=======
Maximum amount of messages to hold in queue. If queue is full 
messages will be dropped. 
>>>>>>> 52bf0eeb559d1bef0e2ca19b8208f3e4867572a3
`<message clipped>` will be add to the message that fills the queue.

- Setting: **OPTIONAL**, **RELOADABLE**
- Default: *30*
- Format: *int*
- Example:
<<<<<<< HEAD

=======
>>>>>>> 52bf0eeb559d1bef0e2ca19b8208f3e4867572a3
  ```toml
  MessageQueue=30
  ```

## MessageSplit
<<<<<<< HEAD

=======
>>>>>>> 52bf0eeb559d1bef0e2ca19b8208f3e4867572a3
Split messages on `MessageLength` instead of showing the `<message clipped>`
WARNING: this could lead to flooding

- Setting: **OPTIONAL**, **RELOADABLE**
- Format: *boolean*
- Example:
<<<<<<< HEAD

=======
>>>>>>> 52bf0eeb559d1bef0e2ca19b8208f3e4867572a3
  ```toml
  MessageSplit=true
  ```

## Nick *
<<<<<<< HEAD

Your nick on irc.
=======
Your nick on irc. 
>>>>>>> 52bf0eeb559d1bef0e2ca19b8208f3e4867572a3

- Setting: REQUIRED
- Format: *string*
- Example:
<<<<<<< HEAD

=======
>>>>>>> 52bf0eeb559d1bef0e2ca19b8208f3e4867572a3
  ```toml
  Nick="matterbot"
  ```

## NickServNick
<<<<<<< HEAD

If you registered your bot with a service like Nickserv on freenode.
Also being used when `UseSASL=true`
Note: when `UseSASL=true`, this is the name of *your* account.
Note: if you want do to quakenet auth, set NickServNick="<Q@CServe.quakenet.org>"
=======
If you registered your bot with a service like Nickserv on freenode.
Also being used when `UseSASL=true`
Note: when `UseSASL=true`, this is the name of *your* account.
Note: if you want do to quakenet auth, set NickServNick="Q@CServe.quakenet.org"
>>>>>>> 52bf0eeb559d1bef0e2ca19b8208f3e4867572a3

- Setting: **OPTIONAL**
- Format: *string*
- Example:
<<<<<<< HEAD

=======
>>>>>>> 52bf0eeb559d1bef0e2ca19b8208f3e4867572a3
  ```toml
  NickServNick="nickserv"
  ```

## NickServPassword
<<<<<<< HEAD

=======
>>>>>>> 52bf0eeb559d1bef0e2ca19b8208f3e4867572a3
The password you use if you registered your bot with a service like Nickserv on freenode.
Also being used when `UseSASL=true`
Also see `NickServNick`

- Setting: **OPTIONAL**
- Format: *string*
- Example:
<<<<<<< HEAD

=======
>>>>>>> 52bf0eeb559d1bef0e2ca19b8208f3e4867572a3
  ```toml
  NickServPassword="secret"
  ```

## NickServUsername
<<<<<<< HEAD

Only used for quakenet auth.
See <https://github.com/42wim/matterbridge/issues/263> for more info
=======
Only used for quakenet auth.
See https://github.com/42wim/matterbridge/issues/263 for more info
>>>>>>> 52bf0eeb559d1bef0e2ca19b8208f3e4867572a3

- Setting: **OPTIONAL**
- Format: *string*
- Example:
<<<<<<< HEAD

=======
>>>>>>> 52bf0eeb559d1bef0e2ca19b8208f3e4867572a3
  ```toml
  NickServUsername="username"
  ```

## NoSendJoinPart
<<<<<<< HEAD

=======
>>>>>>> 52bf0eeb559d1bef0e2ca19b8208f3e4867572a3
Do not send joins/parts to other bridges
Currently works for messages from the following bridges: irc, mattermost, slack

- Setting: **OPTIONAL**
- Format: *boolean*
- Example:
<<<<<<< HEAD

=======
>>>>>>> 52bf0eeb559d1bef0e2ca19b8208f3e4867572a3
  ```toml
  NoSendJoinPart=true
  ```

<<<<<<< HEAD
## Password
=======
## Password 
>>>>>>> 52bf0eeb559d1bef0e2ca19b8208f3e4867572a3

Password for irc server (if necessary)

- Setting: **OPTIONAL**
- Format: *string*
- Example:
<<<<<<< HEAD

=======
>>>>>>> 52bf0eeb559d1bef0e2ca19b8208f3e4867572a3
  ```toml
  Password="s3cret"
  ```

## Pingdelay

PingDelay specifies how long to wait to send a ping to the irc server.
You can use s for second, m for minute

- Setting: **OPTIONAL**, **RELOADABLE**
Default: "1m"
- Example:
<<<<<<< HEAD

=======
>>>>>>> 52bf0eeb559d1bef0e2ca19b8208f3e4867572a3
  ```toml
  PingDelay="1m"
  ```

## RejoinDelay
<<<<<<< HEAD

=======
>>>>>>> 52bf0eeb559d1bef0e2ca19b8208f3e4867572a3
Delay in seconds to rejoin a channel when kicked

- Setting: **OPTIONAL**, **RELOADABLE**
Default: 0
- Format: *int*
- Example:
<<<<<<< HEAD

=======
>>>>>>> 52bf0eeb559d1bef0e2ca19b8208f3e4867572a3
  ```toml
  RejoinDelay=2
  ```

## RunCommands

RunCommands allows you to send RAW irc commands after connection

- Setting: **OPTIONAL**, **RELOADABLE**
- Format: *List<string>*
- Example:
<<<<<<< HEAD

=======
>>>>>>> 52bf0eeb559d1bef0e2ca19b8208f3e4867572a3
  ```toml
  RunCommands=["PRIVMSG user :hello","PRIVMSG chanserv :something"]
  ```

## Server

<<<<<<< HEAD
irc server to connect to.
=======
irc server to connect to. 
>>>>>>> 52bf0eeb559d1bef0e2ca19b8208f3e4867572a3

- Setting: REQUIRED
- Format: *string* (hostname:port)
- Example:
<<<<<<< HEAD

=======
>>>>>>> 52bf0eeb559d1bef0e2ca19b8208f3e4867572a3
  ```toml
  Server="irc.freenode.net:6667"
  ```

## StripMarkdown

Strips Markdown from messages

- Setting: **OPTIONAL**
- Format: *boolean*
- Example:
<<<<<<< HEAD

=======
>>>>>>> 52bf0eeb559d1bef0e2ca19b8208f3e4867572a3
  ```toml
  StripMarkdown=true
  ```

## UseSASL

Enable SASL (PLAIN) authentication. (freenode requires this from eg AWS hosts)
It uses `NickServNick` and `NickServPassword` as login and password

- Setting: **OPTIONAL**
- Format: *boolean*
- Example:
<<<<<<< HEAD

=======
>>>>>>> 52bf0eeb559d1bef0e2ca19b8208f3e4867572a3
  ```toml
  UseSASL=true
  ```

## UseTLS

<<<<<<< HEAD
Enable to use TLS connection to your irc server.
=======
Enable to use TLS connection to your irc server. 
>>>>>>> 52bf0eeb559d1bef0e2ca19b8208f3e4867572a3

- Setting: **OPTIONAL**
- Format: *boolean*
- Example:
<<<<<<< HEAD

=======
>>>>>>> 52bf0eeb559d1bef0e2ca19b8208f3e4867572a3
  ```toml
  UseTLS=true
  ```

## UseRelayMsg

Enable to replace bot's nick with user's nick.
<<<<<<< HEAD

- `RemoteNickFormat` has to contain `/`.
The server has to support RELAYMSG.
Bot may need to be channel operator to use RELAYMSG.
- Setting: OPTIIONAL
- Format: *boolean*
- Example:

=======
- `RemoteNickFormat` has to contain `/`.
The server has to support RELAYMSG.
Bot may need to be channel operator to use RELAYMSG.

- Setting: OPTIIONAL
- Format: *boolean*
- Example:
>>>>>>> 52bf0eeb559d1bef0e2ca19b8208f3e4867572a3
  ```toml
  UseRelayMsg=true
  ```

## VerboseJoinPart

Enable to show verbose users joins/parts (ident@host) from other bridges
Currently works for messages from the following bridges: irc

- Setting: **OPTIONAL **
- Format: *boolean*
- Example:
<<<<<<< HEAD

=======
>>>>>>> 52bf0eeb559d1bef0e2ca19b8208f3e4867572a3
  ```toml
  VerboseJoinPart=true
  ```
