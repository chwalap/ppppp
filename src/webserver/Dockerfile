FROM alpine:latest AS webserver
WORKDIR /opt/ppppp/
COPY src/shared/city.json /db/cities.json
COPY src/webserver/ext/ /opt/ppppp/ext/
COPY src/webserver/templates/ /opt/ppppp/html/
COPY --from=builder /go/bin/webserver /opt/ppppp/bin/webserver
CMD /opt/ppppp/bin/webserver