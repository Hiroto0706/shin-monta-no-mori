FROM golang:1.22.2-alpine3.19 AS builder

WORKDIR /app
COPY go.mod go.sum ./

# airのインストール
RUN go install github.com/air-verse/air@latest

RUN go mod download

COPY . .