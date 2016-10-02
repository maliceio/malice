# Installation

Directions for installing Malice on all supported platforms is provided here. Please note you *must* have Go 1.5 or higher installed. Also, if using Go 1.5 you must set `GO15VENDOREXPERIMENT=1` before attempting to install.

- [OSX](osx/install.md)
- [Linux](linux/install.md)
- [Windows](windows/install.md)
- [Docker](docker)

#### Tips and Tricks

If you have have [zsh](http://www.zsh.org/) installed you can install the zsh-completions:

```bash
$ cd $GOPATH/src/github.com/maliceio/malice/contrib/completion/zsh
$ ./install.sh
```