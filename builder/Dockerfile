# build all services inside golang container
FROM golang:alpine AS builder
RUN apk add --no-cache --virtual .build-deps bash gcc musl-dev openssl git
COPY src /go/src
RUN go get -d -v worker webserver weather db && \
    go build -ldflags="-extldflags=-static" -tags osusergo,netgo,sqlite_omit_load_extension -o /go/bin/worker worker && \
    go build -ldflags="-extldflags=-static" -tags osusergo,netgo,sqlite_omit_load_extension -o /go/bin/webserver webserver && \
    go build -ldflags="-extldflags=-static" -tags osusergo,netgo,sqlite_omit_load_extension -o /go/bin/weather weather && \
    go build -ldflags="-extldflags=-static" -tags osusergo,netgo,sqlite_omit_load_extension -o /go/bin/db db
