FROM alpine:latest
COPY --from=builder /go/bin/webserver /opt/ppppp/bin/webserver
COPY --from=builder /go/src/shared/city.json /db/cities.json
COPY --from=builder /go/src/webserver/templates/ /opt/ppppp/html/
COPY --from=builder /go/src/webserver/ext/ /opt/ppppp/ext/
WORKDIR /opt/ppppp/
CMD ["/opt/ppppp/bin/webserver"]
