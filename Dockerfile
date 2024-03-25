# Start with the official Go image as the base image
FROM golang:1.17-alpine AS builder

# Set the working directory inside the container
WORKDIR /app

# Copy the Go modules and build files
COPY go.mod go.sum ./
RUN go mod download

# Copy the rest of the application source code
COPY . .

# Build the Go application
RUN go build -o myapp .

# Start a new stage from scratch
FROM alpine:latest

# Set the working directory inside the container
WORKDIR /app

# Copy the built executable from the previous stage
COPY --from=builder /app/myapp .

# Expose the port your application listens on
EXPOSE 8080

# Command to run the application
CMD ["./myapp"]
