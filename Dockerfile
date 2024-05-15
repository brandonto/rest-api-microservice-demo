FROM golang:1.21 AS builder

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -o main .

FROM alpine:latest

WORKDIR /app
COPY --from=builder /app/main .
EXPOSE 8080
CMD ["./main"]