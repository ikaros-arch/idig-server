# Stage 1: Build the Go application
FROM golang:1.23 as builder

# Set the working directory inside the container
WORKDIR /app

# Copy the Go module files and download dependencies
# COPY go.mod go.sum ./
# RUN go mod download

# Copy the source code
COPY ./src .

# Build the Go application
RUN go build -o idig-server

# Stage 2: Run the application
FROM debian:bookworm-slim

# Set the working directory inside the container
WORKDIR /app

# Copy the built binary from the builder stage
COPY --from=builder /app/idig-server .

# Copy the entrypoint script
COPY ./entrypoint.sh .

# Make the entrypoint script executable
RUN chmod +x entrypoint.sh

# Expose the port the application runs on
EXPOSE 9000

# Set the entrypoint script as the default command
ENTRYPOINT ["./entrypoint.sh"]
