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
WORKDIR /app/go/user
RUN go build -o user_service main.go

# Final stage
FROM alpine:latest
WORKDIR /app

# Copy the built binary from the builder
COPY --from=builder /app/go/user/user_service .

# Expose the port your gateway listens on (change if needed)
EXPOSE 50051

# Environment variables
ENV CONFIG_PATH=/app/go/user/config/

# Copy the config file
COPY --from=builder /app/go/user/config/user_service_config.json ./go/user/config/user_service_config.json
# Copy sample data file
COPY --from=builder /app/go/user/data/user_db/users.json ./go/user/data/user_db/users.json

# Run the gateway server
CMD ["./user_service"]
