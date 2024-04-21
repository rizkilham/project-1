# Use the official Golang image based on Alpine Linux as builder
FROM golang:alpine AS builder

# Set the current working directory inside the container
WORKDIR /app

# Copy the Go module files and download dependencies
COPY go.mod .

RUN go mod download

# Copy the rest of the application source code
COPY . .

# Build the Go application
RUN go build -o webserver .

# Start a new stage from Alpine Linux
FROM alpine:latest

# Set the working directory inside the container
WORKDIR /app

# Install necessary packages for Redis and MongoDB clients
RUN apk update && apk add --no-cache redis mongodb-tools

# Copy the compiled Go application from the builder stage
COPY --from=builder /app/webserver .

# Expose port 8080 to the outside world
EXPOSE 8080

# Command to run the executable
CMD ["./webserver"]
