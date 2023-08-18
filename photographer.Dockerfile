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

FROM tezos/tezos:${VERSION}
COPY --from=builder /main ./

USER tezos
ENV USER=tezos

ENTRYPOINT ["./main"]
CMD [""]
