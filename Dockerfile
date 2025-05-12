# Build stage
FROM golang:1.24-alpine AS builder

# Set working directory
WORKDIR /app

# Copy go mod and sum files
COPY go.mod go.sum ./

# Download all dependencies
RUN go mod download

# Copy the source code
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o bin/app ./cmd/api

# Runtime stage
FROM alpine:latest

# Add necessary packages
RUN apk --no-cache add ca-certificates tzdata

# Set working directory
WORKDIR /app

# Copy binary from build stage
COPY --from=builder /app/bin/app .

# Copy static assets and configuration files
COPY --from=builder /app/db/migration /app/db/migration
COPY --from=builder /app/.env* .

# Set environment variables
ENV PORT=8080
ENV APP_ENV=production

# Expose port
EXPOSE 8080

# Run the application
CMD ["./app"] 