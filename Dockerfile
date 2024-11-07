# Use the official Go image to build the application
FROM golang:1.20-alpine AS builder

# Set the working directory
WORKDIR /app

# Copy go.mod and go.sum and download dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy the application code
COPY . .

# Build the application
RUN go build -o main .

# Use a minimal Docker image to run the compiled Go binary
FROM alpine:latest

# Set working directory and copy the compiled binary
WORKDIR /root/
COPY --from=builder /app/main .

# Copy any other required files (like static assets)
COPY --from=builder /app/uploads ./uploads

# Expose the application port
EXPOSE 8080

# Command to run the application
CMD ["./main"]
