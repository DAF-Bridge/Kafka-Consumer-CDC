# Start from the official Golang base image with Debian-based glibc support
FROM golang:1.24.0-bullseye AS builder

# Set the Current Working Directory inside the container
WORKDIR /app

# Install necessary dependencies
RUN apt-get update && apt-get install -y --no-install-recommends \
    ca-certificates gcc g++ libc-dev librdkafka-dev curl netcat && \
    rm -rf /var/lib/apt/lists/*

# Copy script
COPY create_connector.sh /usr/local/bin/create_connector.sh

# Make the script executable
RUN chmod +x /usr/local/bin/create_connector.sh

# Copy go mod and sum files
COPY go.mod go.sum ./

# Download all dependencies. Dependencies will be cached if go.mod and go.sum are not changed
RUN go mod download && go mod verify

# Copy the source from the current directory to the Working Directory inside the container
COPY . .

# Enable CGO for librdkafka linking
ENV CGO_ENABLED=1
ENV GO111MODULE=on

# Build the Go app
RUN go build -ldflags="-s -w" -v -o /usr/local/bin/app ./cmd/main.go

# Use minimal Debian-based image for production
FROM debian:bullseye-slim AS runner

# Install runtime dependencies
RUN apt-get update && apt-get install -y --no-install-recommends \ 
    ca-certificates netcat curl iputils-ping dnsutils && \
    rm -rf /var/lib/apt/lists/*

WORKDIR /app

# Copy the built binary from the builder
COPY --from=builder /usr/local/bin/app /usr/local/bin/app

# Copy SSL certificates for HTTPS requests
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

COPY --from=builder /usr/local/bin/create_connector.sh /usr/local/bin/create_connector.sh

# Set production environment
ENV ENVIRONMENT=production

EXPOSE 8000

# Command to run the executable
ENTRYPOINT ["sh", "-c", "/usr/local/bin/create_connector.sh && /usr/local/bin/app"]
