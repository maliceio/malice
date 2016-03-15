# magicmime [![Build Status](http://img.shields.io/travis/rakyll/magicmime/master.svg?style=flat)](https://travis-ci.org/rakyll/magicmime) [![GoDoc](https://godoc.org/github.com/rakyll/magicmime?status.svg)](https://godoc.org/github.com/rakyll/magicmime)

`magicmime` is a Go package which allows you to discover a file's mimetype by looking for magic numbers in its content. It could be used as a supplementary for Go's [`mime`](http://golang.org/pkg/mime/) package which only interprets the file extension to detect mimetypes. Internally, it implements [libmagic(3)](http://linux.die.net/man/3/libmagic) bindings.

Tested on Linux and Mac OS X, should be working on BSD. You could be able to build and make it working with Cygwin on Windows.

## Prerequisites
This uses the `libmagic` library from the `file` tool, so you might need to install the development package for `libmagic`. On Debian or Ubuntu: `apt-get install libmagic-dev`. On RHEL, CentOS or Fedora: `yum install file-devel`. On Mac OS X: `brew install libmagic`. If you don't have the required package, compilation will be terminated by an error saying `magic.h` cannot be found.


## Usage
In order to start, go get this repository:

``` sh
go get github.com/rakyll/magicmime
```

### Examples

See [godoc examples](https://godoc.org/pkg/github.com/rakyll/magicmime/#pkg-examples).

## Docs

https://godoc.org/github.com/rakyll/magicmime

## License
    Copyright 2013 Google Inc. All Rights Reserved.
    Licensed under the Apache License, Version 2.0 (the "License");
    you may not use this file except in compliance with the License.
    You may obtain a copy of the License at
         http://www.apache.org/licenses/LICENSE-2.0
    Unless required by applicable law or agreed to in writing, software
    distributed under the License is distributed on an "AS IS" BASIS,
    WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
    See the License for the specific language governing permissions and
    limitations under the License.
