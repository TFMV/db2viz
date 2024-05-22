# Use the official Golang image as the base image
FROM golang:1.22

# Install necessary packages
RUN apt-get update && apt-get install -y \
    wget \
    unzip \
    libxml2 \
    libssl-dev \
    libpam0g \
    libaio1 \
    build-essential \
    postgresql-client

# Set the working directory inside the container
WORKDIR /app

# Copy the Go module files and download dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy the rest of the application code to the container
COPY . .

# Build the Go application
RUN go build -o db2viz ./cmd/main.go

# Command to run the application
CMD ["./db2viz"]
