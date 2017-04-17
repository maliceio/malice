# Installation on OSX

#### Install [Homebrew](http://brew.sh)

```bash
$ /usr/bin/ruby -e "$(curl -fsSL https://raw.githubusercontent.com/Homebrew/install/master/install)"
```

## Download Pre-Compiled Binary

#### Install Malice

```bash
$ wget https://github.com/maliceio/malice/releases/download/v0.2.0-alpha/malice_0.2.0-alpha_linux_amd64.zip -O /tmp/malice.zip
$ sudo unzip /tmp/malice.zip -d /usr/local/bin/
```

#### Uninstall Malice  

```bash
$ rm /usr/local/bin/malice
$ rm -rf ~/.malice
```

> **NOTE:** We are removing the `.malice` config folder in your home directory which also contains gzipped versions of all the files you have scanned.

## Install With Homebrew

#### Install Malice

```bash
$ brew install https://raw.githubusercontent.com/maliceio/malice/master/contrib/homebrew/Formula/malice.rb
```

> **NOTE:** Included with the homebrew install are zsh completions :sunglasses:

#### Uninstall Malice  

```bash
$ brew uninstall malice
$ rm -rf ~/.malice
$ brew prune
$ brew cleanup
```

> **NOTE:** We are removing the `.malice` config folder in your home directory which also contains gzipped versions of all the files you have scanned.
