# Build stage
FROM golang:1.21-alpine AS builder

# Install git and ca-certificates
RUN apk add --no-cache git ca-certificates tzdata curl

# Create non-root user
RUN adduser -D -s /bin/sh -u 1001 appuser

WORKDIR /app

# Copy go mod files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download && go mod verify

# Copy source code
COPY . .

# Build arguments for version info
ARG VERSION=dev
ARG BUILD_TIME
ARG GIT_COMMIT

# Build the application
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build \
    -ldflags="-w -s -X main.Version=${VERSION} -X main.BuildTime=${BUILD_TIME} -X main.GitCommit=${GIT_COMMIT}" \
    -o server cmd/server/main.go

# Final stage
FROM alpine:3.18

# Install ca-certificates and curl for health checks
RUN apk --no-cache add ca-certificates curl tzdata && \
    update-ca-certificates

# Create non-root user
RUN adduser -D -s /bin/sh -u 1001 appuser

WORKDIR /app

# Copy the binary and config files
COPY --from=builder /app/server .
COPY --from=builder /app/.env* ./
COPY --from=builder /etc/passwd /etc/passwd
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

# Change ownership
RUN chown -R appuser:appuser /app

# Switch to non-root user
USER appuser

# Expose port
EXPOSE 8080

# Health check
HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
    CMD curl -f http://localhost:8080/health || exit 1

# Run the application
CMD ["./server"] 