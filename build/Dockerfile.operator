FROM golang:1.12.10-alpine3.10

ENV GOPATH /go

COPY . /go/src/github.com/fusion-app/fusion-app

RUN go build -o /operator /go/src/github.com/fusion-app/fusion-app/cmd/manager

FROM alpine:3.10

COPY --from=0 /operator /usr/local/bin/fusionapp-operator