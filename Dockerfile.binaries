FROM golang

MAINTAINER blacktop, https://github.com/blacktop

RUN apt-get update && apt-get install -y libmagic-dev zip

WORKDIR /go/src/github.com/maliceio/malice

ENTRYPOINT install/scripts/compile.sh
