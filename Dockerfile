FROM golang:1.17-alpine as build
RUN apk --no-cache add ca-certificates tzdata && update-ca-certificates

FROM scratch

COPY --from=build /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=build /usr/share/zoneinfo /usr/share/zoneinfo

COPY ./alertmanager-telegram /
ENTRYPOINT ["/alertmanager-telegram"]
CMD [ "daemon" ]
