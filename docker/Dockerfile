FROM gliderlabs/alpine:3.4

MAINTAINER blacktop, https://github.com/blacktop

RUN apk-install file tini
RUN apk-install -t build-deps go git mercurial build-base file-dev \
  && set -x \
  && echo "Grab gosu for easy step-down from root..." \
  && apk-install -t .gosu-deps \
                    dpkg \
                    gnupg \
                    openssl \
  && dpkgArch="$(dpkg --print-architecture | awk -F- '{ print $NF }')" \
  && wget -O /usr/local/bin/gosu "https://github.com/tianon/gosu/releases/download/$GOSU_VERSION/gosu-$dpkgArch" \
  && wget -O /usr/local/bin/gosu.asc "https://github.com/tianon/gosu/releases/download/$GOSU_VERSION/gosu-$dpkgArch.asc" \
  && export GNUPGHOME="$(mktemp -d)" \
  && gpg --keyserver ha.pool.sks-keyservers.net --recv-keys B42F6819007F00F88E364FD4036A9C25BF357DD4 \
  && gpg --batch --verify /usr/local/bin/gosu.asc /usr/local/bin/gosu \
  && rm -r "$GNUPGHOME" /usr/local/bin/gosu.asc \
  && chmod +x /usr/local/bin/gosu \
  && gosu nobody true \
  && apk del .gosu-deps \  
  && echo "Building malice Go binary..." \
  && git clone https://github.com/maliceio/malice.git /go/src/github.com/maliceio/malice \
  && cd /go/src/github.com/maliceio/malice \
  && export GOPATH=/go \
  && go version \
  && go get \
  && go build -ldflags "-X main.Version=$(cat VERSION) -X main.BuildTime=$(date -u +%Y%m%d)" -o /bin/malice \
  && echo "Copy malice config files..." \
  && mkdir /malice \
  && cp /go/src/github.com/maliceio/malice/data/config.toml /malice \
  && cp /go/src/github.com/maliceio/malice/data/plugins.toml /malice \
  && rm -rf /go \
  && apk del --purge build-deps

ENV MALICE_STORAGE_PATH /malice

VOLUME ["/malice/config"]
VOLUME ["/malice/samples"]

EXPOSE 80 443

WORKDIR /malice/samples

ENTRYPOINT ["gosu","malice","/sbin/tini","--","malice"]

CMD ["--help"]
