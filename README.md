![malice logo](https://raw.githubusercontent.com/maliceio/malice/master/docs/logo/malice.png)

malice
======

[![CircleCI](https://circleci.com/gh/maliceio/malice.png?style=shield)](https://circleci.com/gh/maliceio/malice)
[![License][license]](http://www.apache.org/licenses/LICENSE-2.0)
[![GoDoc](https://godoc.org/github.com/maliceio/malice?status.svg)](https://godoc.org/github.com/maliceio/malice)
[![Gitter Chat][gitter-badge]][gitter-link]

Malice's mission is to be a free open source version of VirusTotal that anyone can use at any scale from an independent researcher to a fortune 500 company.

### Setup Docker

To Run on OSX - Install [Docker for Mac](https://docs.docker.com/docker-for-mac/)

Or install with [homebrew](http://brew.sh).

```bash
$ brew install caskroom/cask/brew-cask
$ brew cask install virtualbox
$ brew install docker
$ brew install docker-machine
$ docker-machine create --driver virtualbox --engine-storage-driver overlay malice
$ eval $(docker-machine env malice)
```

### Getting Started

#### Install  

```bash
$ brew install https://raw.githubusercontent.com/maliceio/malice/master/contrib/homebrew/Formula/malice.rb
```

```
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
  web		Start, Stop Web services
  plugin	List, Install or Remove Plugins
  help		Shows a list of commands or help for one command

Run 'malice COMMAND --help' for more information on a command.
```

### Usage (Docker in Docker)

[![Docker Stars](https://img.shields.io/docker/stars/malice/engine.svg)][hub]
[![Docker Pulls](https://img.shields.io/docker/pulls/malice/engine.svg)][hub]
[![Docker Image](https://img.shields.io/badge/docker image-29.56 MB-blue.svg)][hub]

Install/Update all Plugins

```bash
docker run --rm -v /var/run/docker.sock:/var/run/docker.sock malice/engine plugin update --all
```

Scan a file

```bash
docker run --rm -v /var/run/docker.sock:/var/run/docker.sock \
                -v `pwd`:/malice/samples \
                -e MALICE_VT_API=$MALICE_VT_API \
                malice/engine scan SAMPLE
```

### Documentation

#### Install

```bash
$ go get github.com/maliceio/malice
```
> Malice will have binary releases for all platforms soon.  

To install on Linux:

-	[Ubuntu](https://github.com/maliceio/malice/blob/master/docs/ubuntu.md)

#### Plugins

-	[Plugins List (and growing)](https://github.com/maliceio/malice/blob/master/docs/plugins.md)

#### Examples

-	[Example Output](https://github.com/maliceio/malice/blob/master/docs/example.md)

#### Tips and Tricks

If you have have [zsh](http://www.zsh.org/) installed you can install the zsh-completions:

```bash
$ cd $GOPATH/src/github.com/maliceio/malice/contrib/completion/zsh
$ ./install.sh
```

### Issues

Find a bug? Want more features? Find something missing in the documentation? Let me know! Please don't hesitate to [file an issue](https://github.com/maliceio/malice/issues/new) and I'll get right on it.

### MVP

Minimum Viable Product

> To be able to scan malware on OSX via cli and have the results either sent to stdout as Markdown tables or store results in ELK with an arbitrary amount of registered Malice plugins.

### TODO

-	[ ] Figure out how to do Windows AV ? :confounded:

### CHANGELOG

See [`CHANGELOG.md`](https://github.com/maliceio/malice/blob/master/CHANGELOG.md)

### Contributing

[See all contributors on GitHub](https://github.com/maliceio/malice/graphs/contributors).

Please update the [CHANGELOG.md](https://github.com/maliceio/malice/blob/master/CHANGELOG.md) and submit a [Pull Request on GitHub](https://help.github.com/articles/using-pull-requests/).

### License

Apache License (Version 2.0)  
Copyright (c) 2013 - 2016 **blacktop** Joshua Maine

[malice-logo]: https://raw.githubusercontent.com/maliceio/malice/master/docs/logo/malice.png
[travis-badge]: https://travis-ci.org/maliceio/malice.svg?branch=master
[gitter-badge]: https://badges.gitter.im/maliceio/malice.svg
[gitter-link]: https://gitter.im/maliceio/malice
[license]: https://img.shields.io/badge/licence-Apache%202.0-blue.svg
[hub]: https://hub.docker.com/r/malice/engine/
