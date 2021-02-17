FROM alpine:latest
COPY --from=builder /go/bin/worker /opt/ppppp/bin/worker
COPY --from=builder /go/src/shared/city.json /db/cities.json
WORKDIR /opt/ppppp/
CMD ["/opt/ppppp/bin/worker"]
