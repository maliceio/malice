FROM malice/alpine:tini

MAINTAINER {{ creator }}

COPY . /go/src/{{ plugin_repo }}
RUN apk-install -t build-deps go git mercurial build-base \
  && set -x \
  && echo "Install {{ plugin_name }}..." \
  && echo "Building scan Go binary..." \
  && cd /go/src/{{ plugin_repo }} \
  && export GOPATH=/go \
  && go version \
  && go get \
  && go build -ldflags "-X main.Version=$(cat VERSION) -X main.BuildTime=$(date -u +%Y%m%d)" -o /bin/scan \
  && rm -rf /go /tmp/* \
  && apk del --purge build-deps

WORKDIR /malware

ENTRYPOINT ["gosu","malice","/sbin/tini","--","scan"]

CMD ["--help"]