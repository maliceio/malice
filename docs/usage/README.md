Using Malice
============

-	[Download All Malice Plugins](#download-all-malice-plugins)
-	[Run Malice](#run-malice)
-	[Lookup a Hash](#lookup-a-hash)
-	[Scan Some Malware](#scan-some-malware)
-	[Watch a Folder](#watch-a-folder)

Download All Malice Plugins
---------------------------

```bash
$ malice plugin update --all
```

> **NOTE:** pulling down all of the plugins can take a long time depending on your network speed.

Run Malice
----------

```bash
$ export MALICE_VT_API=<YOUR API KEY>
$ malice
```

> **NOTE:** Malice has just created a `.malice` folder in your home directory. This is used to store the `config.toml/plugins.toml` that you can change.

```bash
Usage: malice [OPTIONS] COMMAND [arg...]

Open Source Malware Analysis Framework

Version: 0.1.0-alpha, build HEAD

Author:
  blacktop - <https://github.com/blacktop>

Options:
  --debug, -D  	Enable debug mode [$MALICE_DEBUG]
  --help, -h   	show help
  --version, -v	print the version

Commands:
  scan		Scan a file
  watch		Watch a folder
  lookup	Look up a file hash
  elk		Start an ELK docker container
  plugin	List, Install or Remove Plugins
  help		Shows a list of commands or help for one command

Run 'malice COMMAND --help' for more information on a command.
```

Lookup a Hash
-------------

```bash
$ malice lookup 6fe80e56ad4de610304bab1675ce84d16ab6988e
```

See [Lookup Output](https://github.com/maliceio/malice/blob/master/docs/examples/lookup.md)

Scan Some Malware
-----------------

```bash
$ malice scan befb88b89c2eb401900a68e9f5b78764203f2b48264fcc3f7121bf04a57fd408
```

See [Scan Output](https://github.com/maliceio/malice/blob/master/docs/examples/scan.md)

Watch a Folder
--------------

```bash
$ malice watch .
```

```bash
INFO[0000] Malice watching folder: .                     env=development
```
