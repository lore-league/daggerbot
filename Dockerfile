# Use Go 1.23 bookworm as base image
FROM golang:alpine AS base

# Development stage
# =============================================================================
# Create a development stage based on the "base" image
FROM base AS development

# Change the working directory to /app
WORKDIR /app

# Install the air CLI for auto-reloading
RUN go install github.com/air-verse/air@latest

# Copy the go.mod and go.sum files to the /app directory
COPY go.mod go.sum ./

# Install dependencies
RUN go mod download

# Start air for live reloading
CMD ["air"]

# Builder stage
# =============================================================================
# Create a builder stage based on the "base" image
FROM base AS builder

# Move to working directory /build
WORKDIR /build

# Copy the go.mod and go.sum files to the /build directory
COPY go.mod go.sum ./

# Install dependencies
RUN go mod download

# Copy the entire source code into the container
COPY . .

# Build the application
# Turn off CGO to ensure static binaries
RUN CGO_ENABLED=0 go build -o daggerbot

# Production stage
# =============================================================================
# Create a production stage to run the application binary
FROM scratch AS production

# Move to working directory /prod
WORKDIR /prod

# Create a non-root user and group
USER 1001

# Copy binary from builder stage
COPY --from=builder /build/daggerbot ./

# Add a healthcheck to ensure the container is running
HEALTHCHECK --interval=30s --timeout=5s --start-period=5s --retries=3 CMD [ "/prod/daggerbot", "--healthcheck" ] || exit 1

# Start the application
CMD ["/prod/daggerbot"]
