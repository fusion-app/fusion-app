FROM golang:1.15.3-alpine3.12

ENV GOPATH /go

COPY . /go/src/github.com/fusion-app/fusion-app

RUN go build -o /operator /go/src/github.com/fusion-app/fusion-app/cmd/manager

FROM alpine:3.12

COPY --from=0 /operator /usr/local/bin/fusionapp-operator