FROM golang:1.21.5 AS builder
WORKDIR /app
RUN apt-get update && apt-get install -y curl
COPY ./ ./

# Build the API binary
RUN go build -o api ./cmd/api