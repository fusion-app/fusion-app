FROM golang:1.12.10-alpine3.10

ENV GOPATH /go

COPY . /go/src/github.com/fusion-app/fusion-app

RUN go build -o /operator /go/src/github.com/fusion-app/fusion-app/cmd/manager && \
    go build -o /backend /go/src/github.com/fusion-app/fusion-app/dashboard/backend

FROM alpine3.10

RUN apk update \
    && apk add --no-cache \
        ca-certificates \
        && update-ca-certificates 2>/dev/null || true

COPY --from=0 /operator /usr/local/bin/fusionapp-operator
COPY --from=0 /backend /usr/local/bin/fusionapp-dashboard
