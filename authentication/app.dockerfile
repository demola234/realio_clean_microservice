# Stage 1: Build the Go service
FROM golang:1.22-alpine AS build

# Install necessary build tools
RUN apk --no-cache add gcc g++ make ca-certificates

# Set the working directory
WORKDIR /go/src/job_portal

# Copy go.mod, go.sum, and vendor directory
COPY go.mod go.sum ./
COPY vendor vendor

# Copy the service code
COPY account account

# Build the service
RUN GO111MODULE=on go build -mod vendor -o ./bin/app ./authentication/cmd/authentication

# Stage 2: Final image with minimal dependencies
FROM alpine:3.18

# Install required certificates for TLS connections
RUN apk --no-cache add ca-certificates

# Set the working directory
WORKDIR /usr/bin

# Copy the compiled binary from the build stage
COPY --from=build /go/src/job_portal/bin/app .

# Expose the port the service listens on
EXPOSE 8080

# Run the service
CMD ["./app"]
