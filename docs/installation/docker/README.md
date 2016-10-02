# Run Malice inside Docker (no binary)

[![Docker Stars](https://img.shields.io/docker/stars/malice/engine.svg)][hub]
[![Docker Pulls](https://img.shields.io/docker/pulls/malice/engine.svg)][hub]
[![Docker Image](https://img.shields.io/badge/docker image-37.04 MB-blue.svg)][hub]

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

[hub]: https://hub.docker.com/r/malice/engine/