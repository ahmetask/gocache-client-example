FROM golang:1.14 AS builder

ADD . /go/src/gocache-client-example

WORKDIR /go/src/gocache-client-example
ENV GO111MODULE=on

RUN go build -mod=mod -o main

FROM debian:9.9-slim

ENV LANG C.UTF-8
ENV GOPATH /go

COPY --from=builder /go/src/gocache-client-example/main   /app/main
WORKDIR /app

RUN chmod +x main

ENTRYPOINT ["./main"]
EXPOSE 8080