![malice logo](https://raw.githubusercontent.com/maliceio/malice/master/docs/images/logo/malice.png)
=============================================================================================

[![CircleCI](https://circleci.com/gh/maliceio/malice.png?style=shield)](https://circleci.com/gh/maliceio/malice) [![License](https://img.shields.io/badge/licence-Apache%202.0-blue.svg)](http://www.apache.org/licenses/LICENSE-2.0) [![Docker Stars](https://img.shields.io/docker/stars/malice/engine.svg)](https://hub.docker.com/r/malice/engine/) [![Docker Pulls](https://img.shields.io/docker/pulls/malice/engine.svg)](https://hub.docker.com/r/malice/engine/) [![Docker Image](https://img.shields.io/badge/docker%20image-30.6%20MB-blue.svg)](https://hub.docker.com/r/malice/engine/)

This repository contains a **Dockerfile** of the [Malice Engine](https://github.com/maliceio/malice).

### Dependencies

-	[alpine:3.5](https://hub.docker.com/r/_/alpine/)

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

## Documentation

- [Documentation](https://github.com/maliceio/malice/blob/master/docs/README.md)
- [Plugins](https://github.com/maliceio/malice/blob/master/docs/plugins)
- [Examples](https://github.com/maliceio/malice/blob/master/docs/examples)
- [Roadmap](https://github.com/maliceio/malice/blob/master/docs/roadmap)
- [Contributing](https://github.com/maliceio/malice/blob/master/CHANGELOG.md)

### CHANGELOG

See [`CHANGELOG.md`](https://github.com/maliceio/malice/blob/master/CHANGELOG.md)

### License

Apache License (Version 2.0)  
Copyright (c) 2013 - 2017 **blacktop** Joshua Maine
