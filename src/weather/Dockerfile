FROM alpine:latest
COPY --from=builder /go/bin/weather /opt/ppppp/bin/weather
WORKDIR /opt/ppppp/
CMD /opt/ppppp/bin/weather
