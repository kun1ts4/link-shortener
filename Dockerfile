FROM golang:1.26-alpine3.23 AS builder

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN go build -o link-shortener ./cmd/link-shortener/main.go

FROM alpine:3.23
WORKDIR /root/
COPY --from=builder /app/link-shortener .
COPY --from=builder /app/config ./config
CMD ["./link-shortener"]