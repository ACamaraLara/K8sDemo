# Use an official Golang image to build the application
FROM golang:1.22.0 as builder

WORKDIR /app

# Copy the api-gateway directory
COPY ./microservices/api-gateway/ ./microservices/api-gateway/
COPY ./shared/ ./shared/

WORKDIR /app/microservices/api-gateway

# Build the Go application
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 GO111MODULE=on go build -a -o api-gateway ./cmd/main.go

# Use a lightweight image for the final container
FROM alpine:latest

# Copy the compiled binary from the builder stage
COPY --from=builder /app/microservices/api-gateway/api-gateway .

# Expose port 80
EXPOSE 80

# Run the binary
CMD ["./api-gateway"]
