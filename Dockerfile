# Start from the official golang image
FROM golang:1.17-alpine AS build

# Set the current working directory inside the container
WORKDIR /app

# Copy go mod and sum files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy the source code into the container
COPY . .

# Build the Go app
RUN go build -o groupF_Assignment_4 .

# Final image
FROM alpine:latest

WORKDIR /app

# Copy the binary from the build stage
COPY --from=build /app/groupF_Assignment_4 .

# Expose port 8080
EXPOSE 8080

# Command to run the executable
CMD ["./groupF_Assignment_4"]
