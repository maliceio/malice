FROM gliderlabs/alpine

MAINTAINER blacktop, https://github.com/blacktop

RUN apk-install -t build-deps go git mercurial build-base file-dev \
  && set -x \
  && echo "Building info Go binary..." \
  && git clone https://github.com/maliceio/malice.git /go/src/github.com/maliceio/malice \
  && cd /go/src/github.com/maliceio/malice \
  && export GOPATH=/go \
  && go version \
  && go get \
  && go build -ldflags "-X main.Version=$(cat VERSION) -X main.BuildTime=$(date -u +%Y%m%d)" -o /bin/malice \
  && rm -rf /go \
  && apk del --purge build-deps

VOLUME ["/malice/data"]

EXPOSE 80 443 9200 9300

WORKDIR /malware

ENTRYPOINT ["/bin/malice"]

CMD ["--help"]

# CMD ["/usr/bin/supervisord"]
