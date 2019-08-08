FROM golang:1.12.7-alpine AS binary

ENV WORKDIR /go/src/github.com/EthVM/docker-templer

WORKDIR ${WORKDIR}

RUN apk -U add openssl git

ADD . ${WORKDIR}

RUN go install

FROM alpine:3.10

COPY --from=binary /go/bin/docker-templer /usr/local/bin

ENTRYPOINT ["docker-templer"]
CMD ["--help"]
