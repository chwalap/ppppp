FROM alpine:latest AS worker
WORKDIR /opt/ppppp/
COPY src/shared/city.json /db/cities.json
COPY --from=builder /go/bin/worker /opt/ppppp/bin/worker
CMD /opt/ppppp/bin/worker
