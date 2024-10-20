# Use official Golang image as the base
FROM golang:1.19-alpine as builder

WORKDIR /app

# Copy go.mod and go.sum to the container
COPY go.mod go.sum ./
RUN go mod download

# Copy the source code to the container
COPY . .

# Build the Go app
RUN go build -o /authentication ./cmd/authentication/main.go

# Use a minimal image for the final executable
FROM alpine:latest

WORKDIR /

# Copy the compiled Go binary from the builder stage
COPY --from=builder /authentication .

# Install golang-migrate CLI tool
RUN apk add --no-cache curl \
    && curl -L https://github.com/golang-migrate/migrate/releases/download/v4.15.2/migrate.linux-amd64.tar.gz | tar xvz \
    && mv migrate /usr/local/bin/

# Copy migrations
COPY ./db/migrations /migrations

# Run the migrations and start the app
CMD migrate -path /migrations -database "${DB_URL}" up && ./authentication
