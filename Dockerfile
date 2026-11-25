FROM golang:1.23-alpine

RUN apk add --no-cache ca-certificates tzdata git

WORKDIR /app