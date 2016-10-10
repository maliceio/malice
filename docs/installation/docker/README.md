Run Malice inside Docker (no binary)
====================================

[![Docker Stars](https://img.shields.io/docker/stars/malice/engine.svg)](https://hub.docker.com/r/malice/engine/) [![Docker Pulls](https://img.shields.io/docker/pulls/malice/engine.svg)](https://hub.docker.com/r/malice/engine/) [![Docker Image](https://img.shields.io/badge/docker image-37.04 MB-blue.svg)](https://hub.docker.com/r/malice/engine/)

-	[Install/Update all Plugins](#installupdate-all-plugins)
-	[Scan a file](#scan-a-file)
-	[Lookup a hash](#lookup-a-hash)
-	[Start ELK](#start-elk)
-	[Watch a folder](#watch-a-folder)
-	[Use **malice/engine** like a host binary](#use-maliceengine-like-a-host-binary)

Install/Update all Plugins
--------------------------

```bash
docker run --rm -v /var/run/docker.sock:/var/run/docker.sock malice/engine plugin update --all
```

Scan a file
-----------

```bash
docker run --rm -v /var/run/docker.sock:/var/run/docker.sock \
                -v `pwd`:/malice/samples \
                -e MALICE_VT_API=$MALICE_VT_API \
                malice/engine scan SAMPLE
```

Lookup a hash
-------------

```bash
docker run --rm -v /var/run/docker.sock:/var/run/docker.sock \
                -e MALICE_VT_API=$MALICE_VT_API \
                malice/engine lookup 6fe80e56ad4de610304bab1675ce84d16ab6988e
```

Start ELK
---------

```bash
docker run --rm -v /var/run/docker.sock:/var/run/docker.sock \
                malice/engine elk
```

Watch a folder
--------------

```bash
docker run --rm -v /var/run/docker.sock:/var/run/docker.sock \
                -v `pwd`:/malice/samples \
                -e MALICE_VT_API=$MALICE_VT_API \
                malice/engine watch .
```

Use **malice/engine** like a host binary
----------------------------------------

Add the following to your bash or zsh profile

```bash
$ alias malice='docker run --rm -v /var/run/docker.sock:/var/run/docker.sock \
                -v `pwd`:/malice/samples \
                -e MALICE_VT_API=$MALICE_VT_API \
                malice/engine $@'
```
