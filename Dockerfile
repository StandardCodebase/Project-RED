# Stage 1: Build the application
FROM golang:1.26-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o /red-engine ./cmd/red/main.go

# Stage 2: Construct the bare execution container
FROM alpine:latest
RUN apk --no-cache add ca-certificates

# FIX: Create the non-root user with a fixed UID/GID of 1000
RUN addgroup -g 1000 redgroup && \
    adduser -u 1000 -G redgroup -s /bin/sh -D reduser

WORKDIR /app
COPY --from=builder /red-engine ./red-engine

# Create the data directory explicitly before changing ownership
RUN mkdir -p /app/data

# Ensure the user owns the application directory
RUN chown -R reduser:redgroup /app

# Switch to the restricted user
USER reduser

EXPOSE 8080
VOLUME ["/app/data"]

CMD ["./red-engine", "-config", "/app/config.json"]
