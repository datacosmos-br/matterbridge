# Matterbridge API settings

> [!TIP]
> This page contains the details about matterbridge API settings. More general information about the matterbridge API can be found in [README.md](README.md).

<<<<<<< HEAD
## BindAddress
=======
## BindAddress 
>>>>>>> 52bf0eeb559d1bef0e2ca19b8208f3e4867572a3

Address to listen on for API

- Setting: **REQUIRED**
- Format: *string*
- Example:
<<<<<<< HEAD

=======
>>>>>>> 52bf0eeb559d1bef0e2ca19b8208f3e4867572a3
  ```toml
  BindAddress="127.0.0.1:4242"
  ```

## Buffer

Amount of messages to keep in memory

- Setting: **OPTIONAL**, **RELOADABLE**
- Format: *int*
- Default: *10*
- Example:
<<<<<<< HEAD

=======
>>>>>>> 52bf0eeb559d1bef0e2ca19b8208f3e4867572a3
  ```toml
  Buffer=1000
  ```

## Token

HTTP Bearer token used for authentication. If unset, no authentication
will be applied at all and anyone who can reach the API will be able
to control your matterbridge instance.

- Setting: **OPTIONAL**, **RELOADABLE**
- Format: *string*
- Example:
<<<<<<< HEAD

=======
>>>>>>> 52bf0eeb559d1bef0e2ca19b8208f3e4867572a3
  ```toml
  Token="mytoken"
  ```
