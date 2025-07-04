# Default task
default:
    just --list

# Build the application
build:
    go build ./cmd/app/main.go

# Run the application for development IE not using a docker container for the backend
dev:
    sudo service postgresql stop
    docker compose up redis_db -d
    docker compose up dev_db -d
    air

# Run the application for production IE using a docker container (Not configured yet)
prod:
    sudo service postgresql stop
    docker compose up redis_db -d
    docker compose up prod_db -d
    docker compose up service_base

# Run tests with coverage
test:
    sudo service postgresql stop
    go test -cover ./tests/integration_tests/ ./tests/unit_tests/

# Run swagger docs generation
swag:
    rm -rf ./docs 
    swag init -g cmd/app/main.go --dir ./