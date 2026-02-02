# Build stage
FROM golang:1.23-alpine AS builder

WORKDIR /app

# Install git for any private dependencies
RUN apk add --no-cache git

# Copy go mod and sum files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy the source code
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -o farmily-app main.go

# Final stage
FROM alpine:latest

WORKDIR /root/

# Install ca-certificates for secure database connections
RUN apk --no-cache add ca-certificates

# Copy the binary from the builder stage
COPY --from=builder /app/farmily-app .

# Copy static assets and templates (preserving directory structure)
COPY --from=builder /app/static ./static
COPY --from=builder /app/app/templates ./app/templates

# Expose port 8000
EXPOSE 8000

# ENV variables can be overridden at runtime
ENV PORT=8000
ENV LOCAL_DB=false

# Command to run
CMD ["./farmily-app"]
