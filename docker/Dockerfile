FROM gliderlabs/alpine:edge

MAINTAINER blacktop, https://github.com/blacktop

RUN apk-install file
RUN apk-install -t build-deps go git mercurial build-base file-dev \
  && set -x \
  && echo "Building info Go binary..." \
  && git clone https://github.com/maliceio/malice.git /go/src/github.com/maliceio/malice \
  && cd /go/src/github.com/maliceio/malice \
  && export GOPATH=/go \
  && go version \
  && go get \
  && go build -ldflags "-X main.Version=$(cat VERSION) -X main.BuildTime=$(date -u +%Y%m%d)" -o /bin/malice \
  && echo "Copy malice config files..." \
  && cp /go/src/github.com/maliceio/malice/data/config.toml /malice/config \
  && cp /go/src/github.com/maliceio/malice/data/plugins.toml /malice/config \
  && rm -rf /go \
  && apk del --purge build-deps

ENV MALICE_STORAGE_PATH /malice

VOLUME ["/malice/config"]
VOLUME ["/malice/samples"]

EXPOSE 80 443 9200 9300

WORKDIR /malware/samples

ENTRYPOINT ["/bin/malice"]

CMD ["--help"]
