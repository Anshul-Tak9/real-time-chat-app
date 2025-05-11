# Build stage
FROM golang:1.23.9-alpine AS builder

WORKDIR /app

# Install git and Air
RUN apk add --no-cache git && \
    go install github.com/air-verse/air@latest

# Copy go mod and sum files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy the source code
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -o main .

# Final stage
FROM golang:1.23.9-alpine

WORKDIR /app

# Install git
RUN apk add --no-cache git

# Install Air
RUN go install github.com/air-verse/air@latest

# Copy the source code
COPY . .

# Copy the binary from builder
COPY --from=builder /app/main .

EXPOSE 4000

# Run the application with Air
CMD ["air"]
