FROM alpine:3.6

LABEL maintainer "https://github.com/blacktop"

RUN apk add --no-cache file tini
RUN apk add --no-cache -t build-deps go git mercurial build-base gcc file-dev libc-dev \
  && set -x \
  && echo "Building malice Go binary..." \
  && git clone https://github.com/maliceio/malice.git /go/src/github.com/maliceio/malice \
  && cd /go/src/github.com/maliceio/malice \
  # fix git redirects
  && git config --global http.https://gopkg.in.followRedirects true \
  && export GOPATH=/go \
  && go version \
  && CGO_ENABLED=0 go build -i -a -installsuffix cgo -ldflags "-X main.Version=$(cat VERSION) -X main.BuildTime=$(date -u +%Y%m%d)" -o /bin/malice \
  && echo "Copy malice config files..." \
  && mkdir /malice \
  && cp /go/src/github.com/maliceio/malice/config/config.toml /malice \
  && cp /go/src/github.com/maliceio/malice/plugins/plugins.toml /malice \
  && rm -rf /go \
  && apk del --purge build-deps

ENV MALICE_STORAGE_PATH /malice
ENV MALICE_IN_DOCKER true

VOLUME ["/malice/config"]
VOLUME ["/malice/samples"]

EXPOSE 80 443

WORKDIR /malice/samples

ENTRYPOINT ["/sbin/tini","--","malice"]
CMD ["--help"]
