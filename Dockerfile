# Build Stage
FROM golang:1.25.3-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

# Build binary
RUN CGO_ENABLED=0 GOOS=linux go build -o chat-service cmd/server/main.go

# Runtime Stage
FROM alpine:latest

WORKDIR /root/

COPY --from=builder /app/chat-service .
COPY .env .
COPY --from=builder /app/web ./web

# Expose port
EXPOSE 8080

# Environment variables should be passed at runtime
# CMD ["./chat-service"]
ENTRYPOINT ["./chat-service"]
