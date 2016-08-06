### Ubuntu 14.04.3

#### Install Go

```bash
$ sudo add-apt-repository ppa:ubuntu-lxc/lxd-stable
$ sudo apt-get update
$ sudo apt-get install golang
# You should add these two lines to you .bashrc file.
$ export GOPATH=$HOME  
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
$ go get github.com/maliceio/malice
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

> **NOTE:** Malice has just created a `.malice` in your home directory. This is used to store the `config/plugin` info as well as to store the **samples** that you scan.

```bash
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
  lookup	Look up a file hash
  elk		Start an ELK docker container
  web		Start, Stop Web services
  plugin	List, Install or Remove Plugins
  help		Shows a list of commands or help for one command

Run 'malice COMMAND --help' for more information on a command.
```
