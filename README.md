![malice logo][malice-logo]
# malice
[![Build Status][travis-badge]](https://travis-ci.org/maliceio/malice)
[![GoDoc](https://godoc.org/github.com/maliceio/malice?status.svg)](https://godoc.org/github.com/maliceio/malice)
[![Gitter Chat][gitter-badge]][gitter-link]
[![License][license]](http://www.apache.org/licenses/LICENSE-2.0)

Malice's mission is to be a free open source version of VirusTotal that anyone can use at any scale from an independent researcher to a fortune 500 company.

### Install

Malice will have binary releases for all platforms.

**If you are building from source, please note that Malice requires Go v1.5 or above!**

```bash
$ go get -u github.com/maliceio/malice
```
If you have have [zsh](http://www.zsh.org/) and installed it via [homebrew](http://brew.sh) you can do this as well:
```bash
$ cd $GOPATH/src/github.com/maliceio/malice/contrib/zsh-completion
$ bash install.sh
```

### Setup
To Run on OSX
 - Install [Homebrew](http://brew.sh)

```bash
$ brew install caskroom/cask/brew-cask
$ brew cask install virtualbox
$ brew install docker
$ brew install docker-machine
$ docker-machine create --driver virtualbox --engine-storage-driver overlay malice
$ eval $(docker-machine env malice)
```

### Usage
```
Usage: malice [OPTIONS] COMMAND [arg...]

Open Source Malware Analysis Framework

Version: 0.1.0-alpha, build HEAD

Author:
  blacktop - <https://github.com/blacktop>

Options:
  --debug, -D	Enable debug mode [$MALICE_DEBUG]
  --help, -h	show help
  --version, -v	print the version

Commands:
  scan		Scan a file
  elk		Start an ELK docker container
  web		Start, Stop Web services
  plugin	List, Install or Remove Plugins
  help		Shows a list of commands or help for one command

Run 'malice COMMAND --help' for more information on a command.
```

### Documentation
#### Plugins
 - [Plugins List (and growing)](https://github.com/maliceio/malice/blob/master/docs/plugins.md)

#### Examples
 - [Examples](https://github.com/maliceio/malice/blob/master/docs/example.md)

### Issues

Find a bug? Want more features? Find something missing in the documentation? Let me know! Please don't hesitate to [file an issue](https://github.com/maliceio/malice/issues/new) and I'll get right on it.

### MVP
Minimum Viable Product

> To be able to scan malware on OSX via cli and have the results either sent to stdout as Markdown tables or store results in ELK with an arbitrary amount of registered Malice plugins.

### TODO
 - [ ] Figure out how to do Windows AV ? :confounded:

### Credits
I want to give a big shout out to [jordan-wright](http://jordan-wright.com), I am using his program [gophish](https://github.com/jordan-wright/gophish) as a template to get me started with the Malice rewrite.  Jordan has an amazing sense of style and his code is elegant and beautiful.  I aspire to be on his level someday.

### License
Apache License (Version 2.0)  
Copyright (c) 2013 - 2016 **blacktop** Joshua Maine

<!-- Links -->
[malice-logo]: https://raw.githubusercontent.com/maliceio/malice/master/docs/logo/malice.png
[travis-badge]: https://travis-ci.org/maliceio/malice.svg?branch=master
[gitter-badge]: https://badges.gitter.im/maliceio/malice.svg
[gitter-link]: https://gitter.im/maliceio/malice
[license]: https://img.shields.io/badge/licence-Apache%202.0-blue.svg
<!-- [license]: https://img.shields.io/github/license/maliceio/malice.svg -->
