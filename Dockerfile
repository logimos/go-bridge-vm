# Build stage
FROM golang:1.21-alpine AS builder

# Install git and ca-certificates
RUN apk add --no-cache git ca-certificates

# Set working directory
WORKDIR /app

# Copy go mod files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o myllm .

# Final stage
FROM alpine:latest

# Install ca-certificates for HTTPS requests and wget for health checks
RUN apk --no-cache add ca-certificates wget

# Create non-root user
RUN addgroup -g 1001 -S appgroup && \
    adduser -u 1001 -S appuser -G appgroup

# Create app directory
RUN mkdir -p /app && chown appuser:appgroup /app

# Set working directory
WORKDIR /app

# Copy binary from builder stage
COPY --from=builder /app/myllm .

# Copy configuration files
COPY --from=builder /app/configs ./configs
COPY --from=builder /app/USAGE_GUIDE.md .

# Change ownership to non-root user
RUN chown -R appuser:appgroup /app

# Switch to non-root user
USER appuser

# Expose port
EXPOSE 8080

# Set default environment variables
ENV AI_PROVIDER=enhanced_local
ENV INTENT_CONFIG_PATH=configs/personal_assistant.json
ENV AI_TEMPERATURE=0.1
ENV AI_MAX_TOKENS=1000

# Health check - check if the server responds
HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
  CMD wget --no-verbose --tries=1 --spider http://localhost:8080/api/v1/debug || exit 1

# Run the application
CMD ["./myllm"] 