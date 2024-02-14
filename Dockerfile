# Use an official Go image as the base image
FROM golang:latest

# Set the working directory inside the container
WORKDIR /app

# Copy the current directory contents into the container's working directory
COPY . .

# Download and install Go dependencies
RUN go mod download

# Build the Go application
RUN go build -o payment_gateway .

# Download wait for it tool.
ADD https://raw.githubusercontent.com/vishnubob/wait-for-it/master/wait-for-it.sh /wait-for-it
RUN chmod +x /wait-for-it