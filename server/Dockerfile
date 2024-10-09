# Stage 1: Build the Go binary using the official Golang image
FROM golang:1.23 as builder

# Set the working directory inside the container
WORKDIR /app

# Copy the entire project from the parent directory (since context is ..)
COPY ./ .

# Move into the server directory to build the server binary
WORKDIR /app/cacophony/server

# Download dependencies
RUN go mod download

# Build the server binary with CGO disabled for static linking
RUN CGO_ENABLED=0 GOOS=linux go build -o s

# Stage 2: Create a lightweight container with only the binary using Alpine Linux
FROM alpine:3.18

# Copy only the server binary from the builder stage
COPY --from=builder /app/cacophony/server/s /s

# Expose port 8080 for the server
EXPOSE 8080

# Ensure the binary has executable permissions
RUN chmod +x /s

CMD ["/bin/sh", "-c", "echo 'Waiting for MySQL to be ready...'; sleep 10; /s"]

