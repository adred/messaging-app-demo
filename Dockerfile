# Dockerfile
FROM golang:1.20-alpine AS builder
WORKDIR /app

# Copy go modules manifests and download dependencies.
COPY go.mod go.sum ./
RUN go mod download

# Copy the source code.
COPY . .

# Build the application.
RUN CGO_ENABLED=0 GOOS=linux go build -o messaging-service ./cmd/messaging-service

FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=builder /app/messaging-service .
EXPOSE 3000
CMD ["./messaging-service"]
