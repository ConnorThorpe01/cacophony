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

# Build the server binary
RUN go build -o server

# Stage 2: Create a lightweight container with only the binary using Alpine Linux
FROM alpine:3.18

# Install any necessary dependencies (e.g., ca-certificates for SSL)
#RUN apk --no-cache add ca-certificates

# Copy only the server binary from the builder stage
COPY --from=builder /app/cacophony/server /server

# Expose port 8080 for the server
EXPOSE 8080

RUN chmod +x ./server/server

# Run the server binary
CMD ["server"]
