#############################################
## [golang builder]  ########################
#############################################
FROM golang:1.8.3 as builder

ARG VERSION
ARG GITCOMMIT

COPY . /go/src/github.com/maliceio/engine
WORKDIR /go/src/github.com/maliceio/engine/

RUN hack/build/binary

#############################################
## [malice image] ###########################
#############################################
FROM alpine:3.6

LABEL maintainer "https://github.com/blacktop"

ENV MALICE_STORAGE_PATH /malice
ENV MALICE_IN_DOCKER true

RUN apk --no-cache add ca-certificates

COPY --from=builder /go/src/github.com/maliceio/engine/cmd/malice/build/malice /bin/malice
WORKDIR /malice/samples

VOLUME ["/malice/config"]
VOLUME ["/malice/samples"]

EXPOSE 80 443

ENTRYPOINT ["malice"]
CMD ["--help"]

#############################################
#############################################
#############################################
