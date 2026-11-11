# Stage 1: Builder
FROM golang:1.27.3-bookworm AS builder

WORKDIR /app

# Copy dependency manifests
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Build Linux executable
ENV CGO_ENABLED=0 GOOS=linux
RUN go build -o server ./cmd/server

# Stage 2: Runtime
FROM alpine:3.24 AS runtime

WORKDIR /app

# Copy executable from builder
COPY --from=builder /app/server .

# Configure runtime user
USER nobody

# Start application
CMD ["./server"]