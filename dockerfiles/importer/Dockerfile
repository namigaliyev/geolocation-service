FROM golang:1.21.5 AS builder
WORKDIR /app
RUN apt-get update && apt-get install -y curl
COPY ./ ./

# Build the importer binary
RUN go build -o importer ./cmd/importer