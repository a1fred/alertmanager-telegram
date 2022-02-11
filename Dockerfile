FROM golang:1.17-alpine as build
RUN apk --no-cache add ca-certificates tzdata && update-ca-certificates

FROM scratch

COPY --from=build /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=build /usr/share/zoneinfo /usr/share/zoneinfo
COPY ./alertmanager-telegram /

HEALTHCHECK --interval=5s --timeout=30s --start-period=5s --retries=3 CMD [ "/alertmanager-telegram", "health"]

ENTRYPOINT ["/alertmanager-telegram"]
CMD [ "daemon" ]
