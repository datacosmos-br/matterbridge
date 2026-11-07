<<<<<<< HEAD
Matterbridge is not going to implement it's own "mediaserver" instead we make use of other tools that are good at this sort of stuff.
This mediaserver will be used to upload media to services that don't have support for uploading images/video/files.
=======
Matterbridge is not going to implement it's own "mediaserver" instead we make use of other tools that are good at this sort of stuff. 
This mediaserver will be used to upload media to services that don't have support for uploading images/video/files.   
>>>>>>> 52bf0eeb559d1bef0e2ca19b8208f3e4867572a3
At this moment this is xmpp and irc

> [!INFO]
> The `MediaServerUpload` option has been deprecated. If you are using it and would like to
> help reimplement and document it, please open an issue or a pull request.

Running the media server requires a web server which publicly serves files
in a given directory, where matterbridge can write the files.

# Use local download

In this case we're using matterbridge to download to a local path your webserver has read access to and matterbridge has write access to.
Matterbridge is running on this same server.

## matterbridge configuration

In this example the local path matterbridge has write access to is `/var/www/matterbridge`

Your server (apache, nginx, ...) exposes this on `http://yourserver.com/matterbridge` (nginx, apache configuration is out of scope)

configuration needs to happen in `[general]`
<<<<<<< HEAD

```text
=======
```
>>>>>>> 52bf0eeb559d1bef0e2ca19b8208f3e4867572a3
[general]
MediaDownloadPath="/var/www/matterbridge"
MediaServerDownload="https://yourserver.com/matterbridge"
```

<<<<<<< HEAD
When using the local download configuration, matterbridge does not clean up any of the content it downloads to the Mediaserver path.

## Sidenote

If you run into issues with the amount of storage available, then it is advised to do an automated cleanup which is to be done externally (i.e. via cron). An example of a clean up script and two examples of cron jobs are provided below. These represent the minimal amount of effort needed to handle this and don't take into account any ability to customize much.

cleanup.sh:

```text
=======
When using the local download configuration, matterbridge does not clean up any of the content it downloads to the Mediaserver path. 
## Sidenote
If you run into issues with the amount of storage available, then it is advised to do an automated cleanup which is to be done externally (i.e. via cron). An example of a clean up script and two examples of cron jobs are provided below. These represent the minimal amount of effort needed to handle this and don't take into account any ability to customize much.

cleanup.sh:
```
>>>>>>> 52bf0eeb559d1bef0e2ca19b8208f3e4867572a3
#!/bin/bash
find /path/to/matterbridge/media -mindepth 1 -mtime +30 -delete
```

This will delete all downloaded content that is more than 30 days old (but not the media directory itself, due to the use of `-mindepth 1`). You should adjust the path and the max age to suit your own needs. It may be helpful to look at the [`find` manual page](https://www.gnu.org/software/findutils/manual/html_node/find_html/index.html).

To run the script as the user running matterbridge, execute `crontab -e` and add the following line to the bottom of the file:
<<<<<<< HEAD

```text
@daily /path/to/cleanup.sh
```

If you want to run it as root, you probably shouldn't (fix your file permissions instead). However, you can add the script to /etc/cron.daily:
`cp /path/to/cleanup.sh /etc/cron.daily`. This will execute it daily automatically in most Redhat and Debian based Linux distros.
=======
```
@daily /path/to/cleanup.sh
```

If you want to run it as root, you probably shouldn't (fix your file permissions instead).  However, you can add the script to /etc/cron.daily:
`cp /path/to/cleanup.sh /etc/cron.daily`. This will execute it daily automatically in most Redhat and Debian based Linux distros.
>>>>>>> 52bf0eeb559d1bef0e2ca19b8208f3e4867572a3
