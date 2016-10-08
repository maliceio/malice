FROM malice/alpine:tini

MAINTAINER blacktop, https://github.com/blacktop

COPY . /go/src/github.com/maliceio/malice-floss
RUN apk-install python py-setuptools
RUN apk-install -t build-deps go git mercurial build-base py-pip python-dev \
  && set -x \
  && echo "Install FLOSS..." \
  && pip install https://github.com/williballenthin/vivisect/zipball/master \
  && pip install https://github.com/fireeye/flare-floss/zipball/master \
  && echo "Building scan Go binary..." \
  && cd /go/src/github.com/maliceio/malice-floss \
  && export GOPATH=/go \
  && go version \
  && go get \
  && go build -ldflags "-X main.Version=$(cat VERSION) -X main.BuildTime=$(date -u +%Y%m%d)" -o /bin/scan \
  && rm -rf /go /tmp/* \
  && apk del --purge build-deps

WORKDIR /malware

ENTRYPOINT ["gosu","malice","/sbin/tini","--","scan"]

CMD ["--help"]