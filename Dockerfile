FROM golang:1.13-alpine3.10

ENV CGO_ENABLED=0

RUN apk add --update \
    git \
  && rm -rf /var/cache/apk/*

WORKDIR /src

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go test -v ./...
