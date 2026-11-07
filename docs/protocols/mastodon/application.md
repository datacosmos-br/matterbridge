# Add application to mastodon

1. Create an mastodon Application

<<<<<<< HEAD
   Go to this url, change the domain to your mastodon website: <https://mastodon.social/settings/applications/new>

2. Set the following values:

=======
   Go to this url, change the domain to your mastodon website: https://mastodon.social/settings/applications/new
   
2. Set the following values:
   
>>>>>>> 52bf0eeb559d1bef0e2ca19b8208f3e4867572a3
   **Application name:** `MatterBridge`

   **Redirect URL:** `urn:ietf:wg:oauth:2.0:oob`

   **Scopes:**

   Check the following: `read`, `profile`, `write:conversations`, `write:statuses`

<<<<<<< HEAD
   Then, save changes.
=======
    Then, save changes.
>>>>>>> 52bf0eeb559d1bef0e2ca19b8208f3e4867572a3

3. Copy tokens to matterbridge.toml

   ```toml
   [mastodon]
   [mastodon.mymastodon]
   Server = "https://mastodon.social"
   ClientID = "<Client key>"
   ClientSecret = "<Client secret>"
   AccessToken = "<Your access token>"
   ```
