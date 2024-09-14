# Build stage
FROM golang:alpine AS builder

WORKDIR /app
COPY . .

RUN go mod tidy
RUN go build -o main .

# Run stage
FROM alpine:latest

WORKDIR /app

COPY --from=builder /app/main .
COPY --from=builder /app/etc/rd-ticker.yaml ./etc/

EXPOSE 9090
CMD ["./main"]
