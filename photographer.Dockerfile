ARG VERSION=latest
FROM golang:1.20 as builder

WORKDIR /app
COPY go.mod ./
COPY go.sum ./
RUN go mod download
COPY ./cmd ./cmd
COPY ./pkg ./pkg
RUN cd cmd/photographer && go build -o /main

FROM tezos/tezos:${VERSION} as tezos

FROM debian:buster-slim
COPY --from=builder /main ./
COPY --from=tezos /usr/lib/ /usr/lib/
COPY --from=tezos /lib/ /lib/
COPY --from=tezos /usr/local/bin/octez-client /usr/local/bin/octez-client
COPY --from=tezos /usr/local/bin/octez-node /usr/local/bin/octez-node
COPY --from=tezos etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt
ENTRYPOINT ["./main"]
CMD [""]
