# Use the official Golang image for development
FROM golang:1.24.4

# Set the Current Working Directory inside the container
WORKDIR /app

# Install poppler-utils for PDF processing
RUN apt-get update && apt-get install -y poppler-utils

# Install Air for hot-reloading. We do this before copying app files to leverage Docker's layer caching.
RUN go install github.com/air-verse/air@latest

# Copy go mod and sum files
COPY go.mod go.sum ./

# Download all dependencies. Dependencies will be cached if the go.mod and go.sum files are not changed
RUN go mod download

# Copy the source code into the container
# This makes the image self-contained. For development, the volume mount in
# docker-compose.yml will override this directory with the host's source code.
COPY . .

# Expose port 8080 to the outside world
EXPOSE 8080

# Command to run the application using Air.
# Using the full path to the binary avoids any potential PATH issues.
CMD ["air", "-c", ".air.toml"]
