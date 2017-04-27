Installation on Linux (**Ubuntu 14.04.5**\)
===========================================

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

#### Fix Elasticsearch

```bash
echo "vm.max_map_count=262144" | sudo tee -a /etc/sysctl.conf
sudo sysctl -w vm.max_map_count=262144
```

Download Pre-Compiled Binary
----------------------------

#### Install Malice

```bash
$ wget https://github.com/maliceio/malice/releases/download/v0.3.2/malice_linux_amd64.tar.gz -O /tmp/malice.tar.gz
$ sudo tar -xzf /tmp/malice.tar.gz -C /usr/local/bin/
```

#### Uninstall Malice

```bash
$ rm /usr/local/bin/malice
$ rm -rf ~/.malice
```

> **NOTE:** We are removing the `.malice` config folder in your home directory which also contains gzipped versions of all the files you have scanned.

Install From Source
-------------------

#### Install Go

```bash
$ GO_VERSION=1.7.1
$ ARCH="$(dpkg --print-architecture)"
$ wget https://storage.googleapis.com/golang/go$GO_VERSION.linux-$ARCH.tar.gz -O /tmp/go.tar.gz
$ tar -C /usr/local -xzf /tmp/go.tar.gz
# You should add these two lines to you .bashrc file.
$ export PATH=$PATH:/usr/local/go/bin
$ export GOPATH=$HOME/go
$ export PATH=$PATH:$GOPATH/bin
```

#### Install Malice

```bash
$ go get -v github.com/maliceio/malice
```

#### Uninstall Malice

```bash
$ rm $GOPATH/src/github.com/maliceio/malice
$ rm $GOPATH/bin/malice
$ rm -rf ~/.malice
```

> **NOTE:** We are removing the `.malice` config folder in your home directory which also contains gzipped versions of all the files you have scanned.
