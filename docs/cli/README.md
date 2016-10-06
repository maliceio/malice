Commands
========

| Command           | Description                                       |
|-------------------|---------------------------------------------------|
| [scan](#scan)     | Scan a file.                                      |
| [watch](#watch)   | Watch a folder.                                   |
| [lookup](#lookup) | Look up a file hash.                              |
| [elk](#elk)       | Start an ELK docker container.                    |
| [web](#web)       | Start, Stop Web services. :construction:          |
| [plugin](#plugin) | List, Install or Remove Plugins.                  |
| [help](#help)     | Shows a list of commands or help for one command. |

scan
----

```bash
Usage: malice scan [OPTIONS] [arg...]
Scan a file

Description:
   File to be scanned.

Options:

   --logs	Display the Logs of the Plugin containers
```

watch
-----

```bash
Usage: malice watch [OPTIONS] [arg...]
Watch a folder

Description:
   Folder to be watched.

Options:

   --logs	Display the Logs of the Plugin containers
```

lookup
------

```bash
Usage: malice lookup [OPTIONS] [arg...]
Look up a file hash

Description:
   Hash to be queried.

Options:

   --logs	Display the Logs of the Plugin containers
```

elk
---

Start an ELK docker container.

web
---

> **NOTE:** the api/web ui is not done yet.

plugin
------

```bash
NAME:
   malice plugin - List, Install or Remove Plugins

USAGE:
   malice plugin command [command options] [arguments...]

COMMANDS:
     list	list enabled installed plugins
     install	install plugin
     remove	remove plugin
     update	update plugin

OPTIONS:
   --help, -h	show help
```

help
----

Shows a list of commands or help for one command.
