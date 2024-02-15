FROM golang:1.21 AS builder

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download
RUN go mod tidy && \
    go mod verify && \
    go mod download

COPY . .


RUN go build -o ./bin/api ./cmd/api

FROM golang:1.21

WORKDIR /app

COPY --from=builder /app/bin/api .

CMD ["./api"]