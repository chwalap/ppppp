FROM alpine:latest AS weather
WORKDIR /opt/ppppp/
COPY --from=builder /go/bin/weather /opt/ppppp/bin/weather
CMD /opt/ppppp/bin/weather
