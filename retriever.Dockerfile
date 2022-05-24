FROM golang:1.18-alpine
WORKDIR /app
COPY go.mod ./
COPY go.sum ./
RUN go mod download
COPY ./cmd ./cmd
COPY ./pkg ./pkg
RUN cd cmd/retriever && go build -o /main
EXPOSE 8080
ENTRYPOINT ["./main"]
CMD [""]
