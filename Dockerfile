# Use the latest Go image
FROM golang:latest

# Set environment variables
ENV GO111MODULE=on
ENV GOPROXY=https://goproxy.cn,direct

# Set the working directory in the container
WORKDIR /app

# Install git and essential build tools
RUN apt-get update && apt-get install -y \
    git \
    wget \
    && rm -rf /var/lib/apt/lists/*

# Fix Git ownership issue by adding /app as a safe directory
RUN git config --global --add safe.directory /app

# Copy go.mod and go.sum
COPY go.mod go.sum ./

# Download dependencies
RUN go mod tidy && go mod download

# Install the correct version of Bee
RUN go install github.com/beego/bee/v2@latest

# Copy the entire project
COPY . .

# Expose the port
EXPOSE 8080

# Start the application using Bee
CMD ["bee", "run", "-buildvcs=false"]


