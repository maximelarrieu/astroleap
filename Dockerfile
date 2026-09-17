# Build stage 1: Compile the Go server and WASM
FROM golang:1.27-alpine AS builder

WORKDIR /app

# Install build dependencies
RUN apk add --no-cache git ca-certificates

# Install Go module dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# 1. Build WASM binary into web/ directory
RUN GOOS=js GOARCH=wasm go build -o web/game.wasm ./cmd/game

# 2. Build HTTP Server binary
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-w -s" -o /server ./cmd/server

# Stage 2: Minimal runtime image
FROM alpine:3.20

RUN apk add --no-cache ca-certificates

WORKDIR /app
COPY --from=builder /server /server

# Cloud Run binds PORT dynamically (default 8080)
ENV PORT=8080
EXPOSE 8080

CMD ["/server"]
