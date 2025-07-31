.PHONY: build test run clean deps lint

# Build the application
build:
	go build -o bin/myllm main.go

# Run tests
test:
	go test -v ./...

# Run tests with coverage
test-coverage:
	go test -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out

# Run the application
run:
	go run main.go

# Install dependencies
deps:
	go mod tidy
	go mod download

# Clean build artifacts
clean:
	rm -rf bin/
	rm -f coverage.out

# Lint the code
lint:
	golangci-lint run

# Run the example script (requires server to be running)
example:
	./example.sh

# Build and run in one command
dev: build
	./bin/myllm

# Docker build
docker-build:
	docker build -t myllm .

# Docker run
docker-run:
	docker run -p 8080:8080 -e OPENAI_API_KEY=$$OPENAI_API_KEY myllm 