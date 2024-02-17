FROM golang:1.21-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download
RUN go mod tidy && \
    go mod verify

COPY . .


RUN CGO_ENABLED=0 GOOS=linux go build -o api ./cmd/api

FROM alpine:edge

COPY --from=builder /app/api api
COPY --from=builder /app/.env .env

EXPOSE 5132

CMD ["./api"]

