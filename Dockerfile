# Build stage
FROM golang:1.25-alpine AS builder

# Set build directory
WORKDIR /app

# Enable Go modules and verify dependencies
COPY go.mod go.sum ./
RUN go mod download && go mod verify

# Copy source code and build optimized binary
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build \
    -ldflags="-s -w" \
    -o main ./cmd/main.go

# Final stage
FROM alpine:3.19

# Security: Add non-root user
RUN addgroup -S appgroup && adduser -S appuser -G appgroup

# Add certificates for HTTPS requests
RUN apk --no-cache add ca-certificates

WORKDIR /app

# Copy the binary from builder to a location not masked by volumes
COPY --from=builder /app/main /usr/local/bin/app-main
COPY --from=builder /app/.env.example .env.example

# Set ownership to non-root user
RUN mkdir -p logs build && \
    chown -R appuser:appgroup /app

# Use non-root user
USER appuser

EXPOSE 8080

CMD ["app-main"]
