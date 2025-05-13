# Use the official Golang image as the base image
FROM golang:1.22-alpine AS builder

# Set the working directory inside the container
WORKDIR /app

# Install timezone data
RUN apk add --no-cache tzdata

# Copy the Go module files
COPY go.mod .
COPY go.sum .

# Download dependencies
RUN go mod download

# Copy the source code
COPY . .
COPY config.yaml .
COPY credentials.json .
COPY credentials-prod.json .
COPY token.json .
COPY token-prod.json .

# Build the Go application
RUN go build -o mailcast cmd/main.go

# Use a minimal Alpine image for the final stage
FROM alpine:latest

# Install timezone data in the final runtime container
RUN apk add --no-cache tzdata

# Set timezone inside Docker (optional)
# ENV TZ=Asia/Jakarta
ENV TZ=UTC

# Set the working directory
WORKDIR /app

# Copy the binary from the builder stage
COPY --from=builder /app/mailcast .
COPY --from=builder /app/config.yaml .
COPY --from=builder /app/credentials.json .
COPY --from=builder /app/credentials-prod.json .
COPY --from=builder /app/token.json .
COPY --from=builder /app/token-prod.json .

# Run the application
CMD ["./mailcast"]