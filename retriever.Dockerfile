FROM golang:1.17-alpine
WORKDIR /app
COPY go.mod ./
COPY go.sum ./
RUN go mod download
COPY ./cmd ./cmd
COPY ./pkg ./pkg
RUN cd cmd/retriever && go build -o /main

FROM tezos/tezos:release
RUN sudo apk add curl lz4 xz jq
COPY --from=0 /main ./
EXPOSE 8080
ENTRYPOINT ["./main"]
CMD [""]
