FROM golang:alpine AS builder

# install gcc for cgo
RUN apk add --no-cache --virtual .build-deps bash gcc musl-dev openssl git

# copy code and resources
COPY src /go/src

# download dependencies & build binaries
RUN go get -d -v worker webserver weather && \
    go build -ldflags="-extldflags=-static" -tags osusergo,netgo,sqlite_omit_load_extension -o /go/bin/worker worker && \
    go build -ldflags="-extldflags=-static" -tags osusergo,netgo,sqlite_omit_load_extension -o /go/bin/webserver webserver && \
    go build -ldflags="-extldflags=-static" -tags osusergo,netgo,sqlite_omit_load_extension -o /go/bin/weather weather
