FROM gliderlabs/alpine:3.4

MAINTAINER blacktop, https://github.com/blacktop

RUN apk-install file tini
RUN apk-install -t build-deps go git mercurial build-base file-dev \
  && set -x \ 
  && echo "Building malice Go binary..." \
  && git clone https://github.com/maliceio/malice.git /go/src/github.com/maliceio/malice \
  && cd /go/src/github.com/maliceio/malice \
  && export GOPATH=/go \
  && go version \
  && go get \
  && go build -ldflags "-X main.Version=$(cat VERSION) -X main.BuildTime=$(date -u +%Y%m%d)" -o /bin/malice \
  && echo "Copy malice config files..." \
  && mkdir /malice \
  && cp /go/src/github.com/maliceio/malice/config/config.toml /malice \
  && cp /go/src/github.com/maliceio/malice/plugins/plugins.toml /malice \
  && rm -rf /go \
  && apk del --purge build-deps

ENV MALICE_STORAGE_PATH /malice

VOLUME ["/malice/config"]
VOLUME ["/malice/samples"]

EXPOSE 80 443

WORKDIR /malice/samples

ENTRYPOINT ["/sbin/tini","--","malice"]

CMD ["--help"]
