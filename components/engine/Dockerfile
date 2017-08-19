#############################################
## [golang builder]  ########################
#############################################
FROM golang:1.8.3 as builder

COPY . /go/src/github.com/maliceio/engine
WORKDIR /go/src/github.com/maliceio/engine/

RUN go get -v -d
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo \
  -ldflags "-X main.Version=$(cat VERSION) -X main.BuildTime=$(date -u +%Y%m%d)" -o app .

#############################################
## [malice image] ###########################
#############################################
FROM alpine:3.6

LABEL maintainer "https://github.com/blacktop"

ENV MALICE_STORAGE_PATH /malice
ENV MALICE_IN_DOCKER true

RUN apk --no-cache add ca-certificates

COPY --from=builder /go/src/github.com/maliceio/engine/app /bin/malice
WORKDIR /malice/samples

VOLUME ["/malice/config"]
VOLUME ["/malice/samples"]

EXPOSE 80 443

ENTRYPOINT ["malice"]
CMD ["--help"]

#############################################
#############################################
#############################################
