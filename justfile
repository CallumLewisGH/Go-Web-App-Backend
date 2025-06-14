# Default task
default:
    just --list

# Build the application
build:
    go build ./cmd/app/main.go

# Run the application for development IE not using a docker container for the backend
dev:
    docker compose up dev_db -d
    go run ./cmd/app/main.go

# Run the application for production IE using a docker container
prod:
    docker compose up dev_db -d
    docker compose up service_base

# Run tests with coverage
test:
    go test ./tests/integration_tests

# Run swagger docs generation
swagger:
    rm -rf ./docs 
    swag init -g cmd/app/main.go --dir ./