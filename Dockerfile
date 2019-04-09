# build stage
FROM golang:1.12-alpine AS builder
RUN apk add --update ca-certificates tzdata && update-ca-certificates

# final stage
FROM scratch

COPY --from=builder /usr/share/zoneinfo /usr/share/zoneinfo
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY scooter-spotter /bin/

EXPOSE 80
ENTRYPOINT ["/bin/scooter-spotter"]
