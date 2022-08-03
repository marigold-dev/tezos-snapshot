FROM golang:1.18-alpine
WORKDIR /app
COPY go.mod ./
COPY go.sum ./
RUN go mod download
COPY ./cmd ./cmd
COPY ./pkg ./pkg
RUN cd cmd/photographer && go build -o /main

FROM tezos/tezos:v14.0
RUN sudo apk add curl lz4 xz jq
COPY --from=0 /main ./
ENTRYPOINT ["./main"]
CMD [""]
