FROM golang:latest AS builder
COPY src/github.com /go/src/github.com
COPY src/golang.org /go/src/golang.org
COPY src/shared /go/src/shared
COPY src/weather /go/src/weather
RUN go get weather
RUN go build -o /go/bin/weather weather

FROM golang:latest
COPY --from=builder /go/bin/weather /opt/ppppp/bin/weather
WORKDIR /opt/ppppp/
CMD /opt/ppppp/bin/weather
