FROM golang:1.23.2 AS builder

WORKDIR /tasks

COPY go.mod go.sum ./

RUN go mod tidy

COPY . .

RUN mkdir -p /tasks/bin && go build -o /tasks/bin/main ./cmd

FROM debian:bookworm-slim

RUN apt-get update && apt-get install -y sqlite3

COPY --from=builder /tasks/bin/main /main

COPY config.json /config.json

CMD ["/main"]
