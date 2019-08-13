FROM golang:1.12-stretch as builder
ADD . /src
RUN cd /src && CGO_ENABLED=0 make all

FROM alpine:latest as certs
RUN apk --update add ca-certificates

FROM scratch as prod
COPY --from=certs /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt

COPY --from=builder /src/bin/sensor-pipe .

CMD ["./sensor-pipe"]