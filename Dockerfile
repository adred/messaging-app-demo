FROM golang:1.20-alpine AS builder
WORKDIR /app

# Copy go module files and download dependencies.
COPY go.mod go.sum ./
RUN go mod download

# Copy the entire source code.
COPY . .

# Build the application.
RUN CGO_ENABLED=0 GOOS=linux go build -o messaging-service ./cmd/messaging-service

FROM alpine:latest
RUN apk --no-cache add ca-certificates

# Set the working directory to /app.
WORKDIR /app

# Copy the built binary from the builder stage.
COPY --from=builder /app/messaging-service .

# Copy the static and docs directories to the final image.
COPY --from=builder /app/static ./static
COPY --from=builder /app/docs ./docs

EXPOSE 3000
CMD ["./messaging-service"]
