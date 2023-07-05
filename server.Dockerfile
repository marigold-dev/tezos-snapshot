FROM golang:1.20 as builder
WORKDIR /app
COPY go.mod ./
COPY go.sum ./
RUN go mod download
COPY ./cmd ./cmd
COPY ./pkg ./pkg
# To use the libc functions for net and os/user, and still get a static binary (for containers)
# https://github.com/remotemobprogramming/mob/issues/393
RUN cd cmd/server && go build -ldflags "-linkmode 'external' -extldflags '-static'" -o /main

FROM debian:buster-slim
COPY --from=builder /main ./
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt
EXPOSE 8080
ENTRYPOINT ["/main"]
CMD [""]
