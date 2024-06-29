FROM golang:1.22.2-alpine3.19 AS builder

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY ./server ./server
RUN go build -o /build/main ./server/cmd/main.go

FROM alpine:3.19
WORKDIR /app
COPY --from=builder /build/main /usr/local/bin/main
COPY ./server/app.env /app/app.env
COPY ./server/internal/db/migration /app/internal/db/migration
COPY ./credential.json /credential.json

EXPOSE 8080
CMD [ "/usr/local/bin/main" ]