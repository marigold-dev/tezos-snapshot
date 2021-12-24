FROM golang:1.17-alpine
WORKDIR /app
COPY go.mod ./
COPY go.sum ./
RUN go mod download
COPY services/photographer/*.go ./
RUN go build -o /main

FROM tezos/tezos:master
RUN sudo apk add curl lz4 xz jq
COPY --from=0 /main ./
ENTRYPOINT ["./main"]
CMD [""]
