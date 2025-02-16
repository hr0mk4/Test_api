FROM golang:1.23.1-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod tidy

COPY . .

RUN go build -o test_api ./cmd/app

FROM alpine:latest

WORKDIR /app

COPY --from=builder /app/test_api .

EXPOSE 8080

CMD ["./test_api"]
