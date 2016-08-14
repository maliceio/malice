![malice logo](https://raw.githubusercontent.com/maliceio/malice/master/docs/logo/malice.png)
=============================================================================================

[![CircleCI](https://circleci.com/gh/maliceio/malice.png?style=shield)](https://circleci.com/gh/maliceio/malice) [![License](https://img.shields.io/badge/licence-Apache%202.0-blue.svg)](http://www.apache.org/licenses/LICENSE-2.0) [![Docker Stars](https://img.shields.io/docker/stars/malice/engine.svg)](https://hub.docker.com/r/malice/engine/) [![Docker Pulls](https://img.shields.io/docker/pulls/malice/engine.svg)](https://hub.docker.com/r/malice/engine/) [![Docker Image](https://img.shields.io/badge/docker image-29.56 MB-blue.svg)](https://hub.docker.com/r/malice/engine/)

This repository contains a **Dockerfile** of the [Malice Engine](https://github.com/maliceio/malice).

### Dependencies

-	[gliderlabs/alpine:3.4](https://index.docker.io/_/gliderlabs/alpine/)

### Image Tags

```bash
REPOSITORY          TAG                 SIZE
malice/engine       latest              29.56 MB
malice/engine       0.1                 29.56 MB
```

### Installation

1.	Install [Docker](https://docs.docker.com).
2.	Install [docker-compose](https://docs.docker.com/compose/install/)
3.	Download [trusted build](https://hub.docker.com/r/malice/engine/) from public [Docker Registry](https://hub.docker.com/): `docker pull malice/engine`

### Getting Started

#### Install/Update all Plugins

```bash
docker run --rm -v /var/run/docker.sock:/var/run/docker.sock malice/engine plugin update --all
```

#### Scan a file

```bash
docker run --rm -v /var/run/docker.sock:/var/run/docker.sock \
                -v `pwd`:/malice/samples \
                -e MALICE_VT_API=$MALICE_VT_API \
                malice/engine scan SAMPLE
```

### Documentation

#### Usage

#### Plugins

-	[Plugins List (and growing)](https://github.com/maliceio/malice/blob/master/docs/plugins.md)

#### Examples

-	[Examples](https://github.com/maliceio/malice/blob/master/docs/example.md)

#### Tips and Tricks

### Known Issues

### Issues

Find a bug? Want more features? Find something missing in the documentation? Let me know! Please don't hesitate to [file an issue](https://github.com/maliceio/malice/issues/new) and I'll get right on it.

### Todo

### CHANGELOG

See [`CHANGELOG.md`](https://github.com/maliceio/malice/blob/master/CHANGELOG.md)

### Contributing

[See all contributors on GitHub](https://github.com/maliceio/malice/graphs/contributors).

Please update the [CHANGELOG.md](https://github.com/maliceio/malice/blob/master/CHANGELOG.md) and submit a [Pull Request on GitHub](https://help.github.com/articles/using-pull-requests/).

### License

Apache License (Version 2.0)  
Copyright (c) 2013 - 2016 **blacktop** Joshua Maine
