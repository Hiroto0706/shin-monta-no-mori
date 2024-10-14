FROM golang:1.22.2-alpine3.19 AS builder
WORKDIR /app
COPY go.mod go.sum ./
# airのインストール
RUN go install github.com/cosmtrek/air@v1.51.0
RUN go mod download
COPY . .
RUN go build -o /build/main ./cmd/main.go

FROM alpine:3.19
WORKDIR /app
COPY --from=builder /build/main /usr/local/bin/main
COPY ./app.env /app/app.env
COPY ./internal/db/migration /app/internal/db/migration
COPY ./credential.json /app/credential.json

# 環境変数を読み込み
ENV CREDENTIAL_FILE_PATH=/app/credential.json