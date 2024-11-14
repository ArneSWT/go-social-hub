# Use the official Golang image as the base image
FROM golang:1.20

# Install GCC
RUN apt-get update && apt-get install -y gcc

# Set the working directory inside the container
WORKDIR /app

# Copy the Go module files
COPY go.mod go.sum ./

# Download the Go module dependencies
RUN go mod download

# Copy the rest of the application code
COPY . .

# Build the Go application
RUN CGO_ENABLED=1 GOOS=linux go build -o main .

# Expose the port the application runs on
EXPOSE 8080

# Run the Go application
CMD ["./main"]