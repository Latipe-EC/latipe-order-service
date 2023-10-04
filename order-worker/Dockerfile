FROM ubuntu:latest
FROM golang:1.20
LABEL service-owner="Latipe"


ENV GO111MODULE=on
  CGO_ENABLED=0

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY *.go ./


RUN CGO_ENABLED=0 GOOS=linux go build -o apiserver ./cmd/main.go


CMD ["./apiserver"]