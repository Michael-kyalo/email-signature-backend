# Use the official Golang image as a builder
FROM golang:1.23 as builder

# Set the working directory
WORKDIR /app

# Copy go.mod and go.sum
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy the rest of the application source code
COPY . .

# Build the application binary
RUN go build -o main .

# Final stage: Updated base image with newer glibc
FROM ubuntu:22.04

# Set up working directory
WORKDIR /app

# Copy the built binary from the builder stage
COPY --from=builder /app/main .

# Copy any other necessary files (e.g., migrations or assets)
COPY ./db ./db

# Expose the application port
EXPOSE 3000

# Run the application
CMD ["./main"]
