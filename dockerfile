FROM golang:latest

# Set the working directory
WORKDIR /usr/src/app

# Pre-copy/cache go.mod for pre-downloading dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy the rest of the application code
COPY . .

# Build the application
RUN go build -v -o /usr/local/bin/app ./cmd/app

# Set the port environment variable
ENV PORT=8080

# Expose the port
EXPOSE 8080

# Run the application
CMD ["/usr/local/bin/app"]