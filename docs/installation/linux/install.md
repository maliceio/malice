# Installation on Linux (**Ubuntu 14.04.3**)

#### Install Go

```bash
$ GO_VERSION=1.7.1
$ ARCH="$(dpkg --print-architecture)" \
$ wget https://storage.googleapis.com/golang/go$GO_VERSION.linux-$ARCH.tar.gz -O /tmp/go.tar.gz \
$ tar -C /usr/local -xzf /tmp/go.tar.gz \
# You should add these two lines to you .bashrc file.
$ export PATH=$PATH:/usr/local/go/bin
$ export GOPATH=$HOME/go
$ export PATH=$PATH:$GOPATH/bin
```

#### Install Docker

```bash
$ sudo apt-get update
$ sudo apt-get install apt-transport-https ca-certificates
$ sudo apt-key adv --keyserver hkp://p80.pool.sks-keyservers.net:80 --recv-keys 58118E89F3A912897C070ADBF76221572C52609D
$ echo "deb https://apt.dockerproject.org/repo ubuntu-trusty main" > sudo tee -a /etc/apt/sources.list.d/docker.list
$ sudo apt-get update
$ sudo apt-get install docker-engine
$ sudo usermod -aG docker $USER  # You might have to logout for change to take effect
```

#### Install Malice

```bash
$ sudo apt-get install libmagic-dev build-essential
$ go get -v github.com/maliceio/malice
```

#### Download All Malice Plugins

```bash
$ malice plugin update --all
```

> **NOTE:** pulling down all of the plugins can take a long time depending on your network speed.

#### Run Malice

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
  web		Start, Stop Web services
  plugin	List, Install or Remove Plugins
  help		Shows a list of commands or help for one command

Run 'malice COMMAND --help' for more information on a command.
```
