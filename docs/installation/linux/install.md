Installation on Linux (**Ubuntu 18.10**\)
===========================================

#### Install Docker

```bash
$ sudo apt-get install apt-transport-https ca-certificates curl software-properties-common
$ curl -fsSL https://download.docker.com/linux/ubuntu/gpg | sudo apt-key add - # Adding GPG Docker's keys
$ sudo apt-key fingerprint 0EBFCD88
$ sudo add-apt-repository "deb [arch=amd64] https://download.docker.com/linux/ubuntu $(lsb_release -cs) stable" # Adding Docker's repository in Ubuntu repository list
$ sudo apt-get update
$ sudo apt-get install docker docker.io
$ sudo usermod -aG docker $USER # Ading Docker as a user
$ sudo reboot now # Ubuntu will restart, save your work
```

#### Fix Elasticsearch

```bash
echo "vm.max_map_count=262144" | sudo tee -a /etc/sysctl.conf
sudo sysctl -w vm.max_map_count=262144
sudo docker rm -f malice
```

Download Pre-Compiled Binary
----------------------------

#### Install Malice

> **NOTE:** Grab the **latest** release [here](https://github.com/maliceio/malice/releases/latest)

```bash
$ wget https://github.com/maliceio/malice/releases/download/v0.3.11/malice_0.3.11_linux_amd64.tar.gz -O /tmp/malice.tar.gz
$ sudo tar -xzf /tmp/malice.tar.gz -C /usr/local/bin/
```
#### Update Malice
```bash
$ malice plugin update --all # Updating plugins
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
$ GO_VERSION=1.8.3
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
