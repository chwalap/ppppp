FROM alpine:latest AS db
WORKDIR /opt/ppppp/
COPY --from=builder /go/bin/db /opt/ppppp/bin/db
CMD /opt/ppppp/bin/db
