# Stage 1: Build the Go application
FROM golang:1.23-alpine AS builder

# Install required packages
RUN apk update && apk add --no-cache git

# Set environment variables
ENV GO111MODULE=on
WORKDIR /app

# Copy go.mod and go.sum for dependency installation
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy the rest of the source code
COPY . .

# Build the Go applications
RUN go build -o tic-tac-toe-app main.go
RUN go build -o kafka-consumer kafka_consumer.go

# Stage 2: Create a minimal image for runtime
FROM alpine:latest

# Install necessary certificates
RUN apk --no-cache add ca-certificates

WORKDIR /root/

# Copy the built binaries
COPY --from=builder /app/tic-tac-toe-app .
COPY --from=builder /app/kafka-consumer .

# Copy static files for the web UI
COPY --from=builder /app/static ./static


# Expose the web app port
EXPOSE 8080

# Script to run both the main app and the Kafka consumer in the same container
COPY run.sh /root/run.sh
RUN chmod +x /root/run.sh

# Start both the main app and Kafka consumer
CMD ["./run.sh"]
