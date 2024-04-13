FROM golang:1.22.2-alpine3.19 AS builder

WORKDIR /app
COPY go.mod go.sum ./
COPY . .
RUN go build -o /build/main ./cmd/main.go

FROM alpine:3.19
WORKDIR /app
COPY --from=builder /build/main /usr/local/bin/main
# COPY app.env .
# COPY credential.json .
# COPY db/migration ./db/migration