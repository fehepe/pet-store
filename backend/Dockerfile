# Build stage
FROM golang:1.24-alpine AS builder

# Install build dependencies
RUN apk add --no-cache git

# Set working directory
WORKDIR /app

# Copy go mod files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main ./cmd/server

# Final stage
FROM alpine:latest

# Install runtime dependencies
RUN apk --no-cache add ca-certificates

# Create non-root user
RUN addgroup -g 1000 -S petstore && \
    adduser -u 1000 -S petstore -G petstore

# Set working directory
WORKDIR /app

# Copy binary from builder
COPY --from=builder /app/main .

# Copy migrations
COPY --from=builder /app/internal/database/migrations ./internal/database/migrations

# Create uploads directory and set permissions
RUN mkdir -p ./uploads && \
    chown -R petstore:petstore ./uploads && \
    chown -R petstore:petstore ./internal

# Change to non-root user
USER petstore

# Expose port
EXPOSE 8080

# Run the application
CMD ["./main"]