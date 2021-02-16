FROM golang:latest AS builder
COPY src/github.com /go/src/github.com
COPY src/golang.org /go/src/golang.org
COPY src/shared /go/src/shared
COPY src/webserver /go/src/webserver
RUN go get webserver
RUN go build -o /go/bin/webserver webserver

FROM golang:latest
COPY --from=builder /go/bin/webserver /opt/ppppp/bin/webserver
COPY --from=builder /go/src/shared/city.json /db/cities.json
COPY --from=builder /go/src/webserver/templates/ /opt/ppppp/html/
COPY --from=builder /go/src/webserver/ext/ /opt/ppppp/ext/
WORKDIR /opt/ppppp/
CMD ["/opt/ppppp/bin/webserver"]
