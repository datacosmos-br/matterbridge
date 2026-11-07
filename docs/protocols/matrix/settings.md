# Matrix settings

> [!TIP]
> This page contains the details about matrix settings. More general information about matrix support in matterbridge can be found in [README.md](README.md).

## DeviceID

The device id use when logging in with MxID.

Unless this option is set, the Matrix client is unencrypted and MxID based login won't work.

- Setting: **OPTIONAL**
- Format: *string*
- Example:
<<<<<<< HEAD

=======
>>>>>>> 52bf0eeb559d1bef0e2ca19b8208f3e4867572a3
  ```toml
  DeviceID="yourdeviceid"
  ```

## HTMLDisable

Whether to disable sending of HTML content to matrix
<<<<<<< HEAD
See <https://github.com/42wim/matterbridge/issues/1022>
=======
See https://github.com/42wim/matterbridge/issues/1022
>>>>>>> 52bf0eeb559d1bef0e2ca19b8208f3e4867572a3

- Setting: **OPTIONAL**, **RELOADABLE**
- Format: *boolean*
- Example:
<<<<<<< HEAD

=======
>>>>>>> 52bf0eeb559d1bef0e2ca19b8208f3e4867572a3
  ```toml
  HTMLDisable=true
  ```

## Login

login of your bot.
Use a dedicated user for this and not your own!
Messages sent from this user will not be relayed to avoid loops.

- Setting: **REQUIRED**
- Format: *string*
- Example:
<<<<<<< HEAD

=======
>>>>>>> 52bf0eeb559d1bef0e2ca19b8208f3e4867572a3
  ```toml
  Login="yourlogin"
  ```

## MxID

MxID of your bot.
Use a dedicated user for this and not your own!
Messages sent from this user will not be relayed to avoid loops.

- Setting: **REQUIRED**
- Format: *string*
- Example:
<<<<<<< HEAD

=======
>>>>>>> 52bf0eeb559d1bef0e2ca19b8208f3e4867572a3
  ```toml
  MxID="@yourbot:example.net"
  ```

## NoHomeServerSuffix

Whether to send the homeserver suffix. eg ":matrix.org" in @username:matrix.org
to other bridges, or only send "username".(true only sends username)

- Setting: **OPTIONAL**, **RELOADABLE**
- Format: *boolean*
- Example:
<<<<<<< HEAD

=======
>>>>>>> 52bf0eeb559d1bef0e2ca19b8208f3e4867572a3
  ```toml
  NoHomeServerSuffix=true
  ```

## Password

password of your bot.

- Setting: **REQUIRED**
- Format: *string*
- Example:
<<<<<<< HEAD

=======
>>>>>>> 52bf0eeb559d1bef0e2ca19b8208f3e4867572a3
  ```toml
  Password="yourpass"
  ```

## PickleKey

The key to use when accessing E2EE encryption in an encryption database.

Unless this option is set, the Matrix client is unencrypted.

- Setting: **OPTIONAL**
- Format: *string*
- Example:
<<<<<<< HEAD

=======
>>>>>>> 52bf0eeb559d1bef0e2ca19b8208f3e4867572a3
  ```toml
  Password="yourpicklekey"
  ```

## RecoveryKey

The key to use when accessing E2EE encryption in an encryption database.

Unless this option is set, the Matrix client won't be verified for encryption.

- Setting: **OPTIONAL**
- Format: *string*
- Example:
<<<<<<< HEAD

=======
>>>>>>> 52bf0eeb559d1bef0e2ca19b8208f3e4867572a3
  ```toml
  RecoveryKey="yourrecoverykey"
  ```

## Server

<<<<<<< HEAD
Server is your homeserver (eg <https://matrix.org>)
=======
Server is your homeserver (eg https://matrix.org)
>>>>>>> 52bf0eeb559d1bef0e2ca19b8208f3e4867572a3

- Setting: **REQUIRED**
- Format: *string*
- Example:
<<<<<<< HEAD

=======
>>>>>>> 52bf0eeb559d1bef0e2ca19b8208f3e4867572a3
  ```toml
  Server="https://matrix.org"
  ```

## SessionFile

The database file to use when accessing E2EE encryption in an encryption database.

Unless this option is set, the Matrix client is unencrypted.

- Setting: **OPTIONAL**
- Format: *string*
- Example:
<<<<<<< HEAD

=======
>>>>>>> 52bf0eeb559d1bef0e2ca19b8208f3e4867572a3
  ```toml
  SessionFile="yourdatabasefile.db"
  ```

## UseUserName

Shows the username instead of the displayname

- Setting: **OPTIONAL**, **RELOADABLE**
- Format: *boolean*
- Example:
<<<<<<< HEAD

=======
>>>>>>> 52bf0eeb559d1bef0e2ca19b8208f3e4867572a3
  ```toml
  UseUserName=true
  ```
