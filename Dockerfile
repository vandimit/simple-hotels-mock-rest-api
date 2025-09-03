# Builder stage
FROM golang:1.19-bullseye AS builder

# Set working directory
WORKDIR /app

# Copy the entire source code
COPY . .

# Build the application using vendored dependencies
RUN CGO_ENABLED=0 GOOS=linux go build -mod=vendor -o hotels-api .

# Runtime stage
FROM debian:bullseye-slim

# Set working directory
WORKDIR /app

# Install ca-certificates for HTTPS
RUN apt-get update && \
    apt-get install -y --no-install-recommends ca-certificates && \
    apt-get clean && \
    rm -rf /var/lib/apt/lists/*

# Copy the binary from the builder stage
COPY --from=builder /app/hotels-api .

# Copy required data files
COPY --from=builder /app/mock-data /app/mock-data
COPY --from=builder /app/public /app/public

# Expose the port the app runs on
EXPOSE 8080

# Command to run the executable
CMD ["./hotels-api"]