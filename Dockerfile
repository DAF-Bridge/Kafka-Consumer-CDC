FROM golang:1.24.0-bullseye AS builder

WORKDIR /app

RUN apt-get update && apt-get install -y --no-install-recommends \
    ca-certificates gcc g++ libc-dev librdkafka-dev curl netcat && \
    rm -rf /var/lib/apt/lists/*

COPY create_connector.sh /usr/local/bin/create_connector.sh

# Make the script executable
RUN chmod +x /usr/local/bin/create_connector.sh

# Copy go mod and sum files
COPY go.mod go.sum ./

RUN go mod download && go mod verify
COPY . .

# Enable CGO for librdkafka linking
ENV CGO_ENABLED=1
ENV GO111MODULE=on

RUN go build -ldflags="-s -w" -v -o /usr/local/bin/app ./cmd/main.go

# Use minimal Debian-based image for production
FROM debian:bullseye-slim AS runner

# Install runtime dependencies
RUN apt-get update && apt-get install -y --no-install-recommends \ 
    ca-certificates netcat curl iputils-ping dnsutils && \
    rm -rf /var/lib/apt/lists/*

WORKDIR /app

# Copy the built from the builder
COPY --from=builder /usr/local/bin/app /usr/local/bin/app
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /usr/local/bin/create_connector.sh /usr/local/bin/create_connector.sh

# Set production environment
ENV ENVIRONMENT=production

EXPOSE 8000

# Command to run the executable
ENTRYPOINT ["sh", "-c", "/usr/local/bin/create_connector.sh && /usr/local/bin/app"]
