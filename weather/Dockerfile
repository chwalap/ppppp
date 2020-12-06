FROM golang:latest

ENV GOPATH=/go

COPY src /go/src
WORKDIR /go/src
RUN go build -o weather.exe weather

CMD ["/go/src/weather.exe"]
