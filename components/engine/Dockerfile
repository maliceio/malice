FROM golang:1.8.3 as builder

COPY . /go/src/github.com/maliceio/engine
WORKDIR /go/src/github.com/maliceio/engine/

RUN go get -v -d
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo \
  -ldflags "-X main.Version=$(cat VERSION) -X main.BuildTime=$(date -u +%Y%m%d)" -o app .

FROM alpine:3.6

LABEL maintainer "https://github.com/blacktop"

RUN apk --no-cache add ca-certificates

WORKDIR /root/

COPY --from=builder /go/src/github.com/maliceio/engine/app /bin/malice 

ENTRYPOINT ["malice"]
CMD ["--help"]