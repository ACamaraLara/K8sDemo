# Use an official Golang image to build the application
FROM golang:1.22.0 as builder

# Define a build-time argument for the service name (default to "api-gateway")
ARG SERVICE_NAME

WORKDIR /app

# Copy the microservice directory using the service name variable
COPY ./microservices/${SERVICE_NAME}/ ./microservices/${SERVICE_NAME}/
COPY ./shared/ ./shared/

WORKDIR /app/microservices/${SERVICE_NAME}

# Build the Go application
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 GO111MODULE=on go build -a -o ${SERVICE_NAME} ./cmd/main.go

# Use a lightweight image for the final container
FROM alpine:latest

ARG SERVICE_NAME
ENV SERVICE_NAME=${SERVICE_NAME}

# Copy the compiled binary from the builder stage using the service name variable
COPY --from=builder /app/microservices/${SERVICE_NAME}/${SERVICE_NAME} .

# Expose port 80
EXPOSE 80

# Run the binary using the service name variable
CMD ["sh", "-c", "./$SERVICE_NAME"]