### Ubuntu 14.04.3

#### Install Go
```bash
$ sudo add-apt-repository ppa:ubuntu-lxc/lxd-stable
$ sudo apt-get update
$ sudo apt-get install golang
```
#### Install Docker
```bash
$ sudo apt-get update
$ sudo apt-get install apt-transport-https ca-certificates
$ sudo apt-key adv --keyserver hkp://p80.pool.sks-keyservers.net:80 --recv-keys 58118E89F3A912897C070ADBF76221572C52609D
$ sudo echo "deb https://apt.dockerproject.org/repo ubuntu-trusty main
" > /etc/apt/sources.list.d/docker.list
$ sudo apt-get update
$ sudo apt-get install docker-engine
```
#### Install Malice
```bash
$ sudo apt-get install libmagic-dev build-essential
$ export GOPATH=$HOME
$ export PATH=$PATH:$GOPATH/bin
$ go get github.com/tools/godep
$ go get -u github.com/maliceio/malice
$ godep restore  # This might not be needed with Go 1.6+
$ cd ~/src/github.com/maliceio/malice
$ sudo -E go install
```
#### Download All Malice Plugins
```bash
$ cd $HOME 
$ sudo bin/malice plugin update --all
```
#### Run Malice
```bash
$ export MALICE_VT_API=<YOUR API KEY>
$ sudo bin/malice
```
> **NOTE:** Malice has just created a `.malice` in your home directory.  This is used to store the config/plugin info as well as to store the samples that you scan.

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
