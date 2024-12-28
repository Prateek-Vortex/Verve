# Step 1: Use Golang base image to build the application
FROM golang:1.23.4-alpine AS build

# Set the Current Working Directory inside the container
WORKDIR /app

# Copy the Go modules and sum files
COPY go.mod go.sum ./

# Download all dependencies. Dependencies are cached if the go.mod and go.sum files are not changed
RUN go mod tidy

# Copy the source code into the container
COPY . .

# Build the Go app
RUN go build -o main ./cmd/api

# Step 2: Create a smaller image to run the app
FROM alpine:latest  

# Install dependencies for running the Go binary
RUN apk --no-cache add ca-certificates

# Set the Current Working Directory inside the container
WORKDIR /root/

# Copy the binary from the build container
COPY --from=build /app/main .

# Expose the port the app runs on
EXPOSE 8080

# Command to run the application
CMD ["./main"]
