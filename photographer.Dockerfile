FROM golang:1.20-alpine
WORKDIR /app
COPY go.mod ./
COPY go.sum ./
RUN go mod download
COPY ./cmd ./cmd
COPY ./pkg ./pkg
RUN cd cmd/photographer && go build -o /main

FROM tezos/tezos:v15.1
RUN sudo apk add curl lz4 xz jq
COPY --from=0 /main ./
ENTRYPOINT ["./main"]
CMD [""]
