FROM golang:1.22.2-alpine3.19

WORKDIR /app

# airをインストール
RUN go install github.com/cosmtrek/air@latest

COPY go.mod .
COPY go.sum .

RUN go mod download

COPY . .

CMD ["air", "-c", ".air.toml"]