# Using Malice

TODO(blacktop)

This directory will hold documentation around using malice and how it works.

## Download All Malice Plugins

```bash
$ malice plugin update --all
```

> **NOTE:** pulling down all of the plugins can take a long time depending on your network speed.

## Run Malice

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