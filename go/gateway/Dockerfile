# syntax=docker/dockerfile:1

# Build stage
FROM golang:1.23-alpine AS builder
WORKDIR /app

# Install git (for go mod download)
RUN apk add --no-cache git

# Copy the entire Go workspace
COPY . .

# Set Go environment variables
ENV CGO_ENABLED=0 \
    GO111MODULE=on

# Build the gateway binary
WORKDIR /app/go/gateway
RUN go build -o gateway server.go

# Final stage
FROM alpine:latest
WORKDIR /app

# Copy the built binary from the builder
COPY --from=builder /app/go/gateway/gateway .

# Expose the port your gateway listens on (change if needed)
EXPOSE 8080

# Copy the config file
COPY --from=builder /app/go/gateway/config/gateway_config.json ./go/gateway/config/gateway_config.json

# Run the gateway server
CMD ["./gateway"]
