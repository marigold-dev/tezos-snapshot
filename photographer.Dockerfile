ARG VERSION=latest
FROM golang:1.21 as builder

WORKDIR /app
COPY go.mod ./
COPY go.sum ./
RUN go mod download
COPY ./cmd ./cmd
COPY ./pkg ./pkg
# To use the libc functions for net and os/user, and still get a static binary (for containers)
# https://github.com/remotemobprogramming/mob/issues/393
RUN cd cmd/photographer && go build -ldflags "-linkmode 'external' -extldflags '-static'" -o /main

FROM tezos/tezos:${VERSION} as tezos

FROM debian:buster-slim
COPY --from=builder /main ./
COPY --from=tezos /usr/lib/ /usr/lib/
COPY --from=tezos /lib/ /lib/
COPY --from=tezos /usr/local/bin/octez-client /usr/local/bin/octez-client
COPY --from=tezos /usr/local/bin/octez-node /usr/local/bin/octez-node
COPY --from=tezos etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt

RUN addgroup -S tezos && adduser -S tezos -G tezos
USER tezos
ENV USER=tezos

ENTRYPOINT ["./main"]
CMD [""]
