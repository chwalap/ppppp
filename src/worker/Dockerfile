FROM golang:latest AS builder
COPY src/github.com /go/src/github.com
COPY src/golang.org /go/src/golang.org
COPY src/shared /go/src/shared
COPY src/worker /go/src/worker
RUN go build -o /go/bin/worker worker

FROM golang:latest
COPY --from=builder /go/bin/worker /opt/ppppp/bin/worker
COPY --from=builder /go/src/shared/city.json /db/cities.json
WORKDIR /opt/ppppp/
CMD ["/opt/ppppp/bin/worker"]
