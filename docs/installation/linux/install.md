# Installation on Linux (**Ubuntu 14.04.5**)

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

## Download Pre-Compiled Binary

#### Install Malice

```bash
$ sudo apt-get install libmagic-dev build-essential unzip
$ wget https://github.com/maliceio/malice/releases/download/v0.2.0-alpha/malice_0.2.0-alpha_linux_amd64.zip -O /tmp/malice.zip
$ sudo unzip /tmp/malice.zip -d /usr/local/bin/
```

## Install From Source

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

#### Install Malice

```bash
$ sudo apt-get install libmagic-dev build-essential
$ go get -v github.com/maliceio/malice
```
