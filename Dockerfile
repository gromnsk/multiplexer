FROM golang:1.14 AS builder

COPY . /multiplexer
RUN cd /multiplexer && make build OUTPUT=/go/bin/multiplexer

FROM alpine:latest AS certs

RUN apk --update upgrade && apk add ca-certificates && update-ca-certificates

FROM scratch

COPY --from=builder /go/bin/multiplexer /multiplexer
COPY --from=certs /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

ENV MULTIPLEXER_HTTP_PORT=8080

EXPOSE $MULTIPLEXER_HTTP_PORT

ENTRYPOINT ["/multiplexer"]
