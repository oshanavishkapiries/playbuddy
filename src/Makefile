.PHONY: build run clean test tidy

# Build the application
build:
	@echo "Building PlayBuddy Torrent Search..."
	go mod tidy
	go build -o bin/playbuddy cmd/main.go

# Run the application
run:
	@echo "Running PlayBuddy Torrent Search..."
	go mod tidy
	go run cmd/main.go

# Clean build artifacts
clean:
	@echo "Cleaning build artifacts..."
	rm -rf bin/
	rm -f *.exe

# Run tests
test:
	@echo "Running tests..."
	go test ./...

# Tidy dependencies
tidy:
	@echo "Tidying dependencies..."
	go mod tidy

# Install dependencies
deps:
	@echo "Installing dependencies..."
	go mod download

# Development setup
setup: deps tidy
	@echo "Setup complete!"

# Build for Windows
build-windows:
	@echo "Building for Windows..."
	GOOS=windows GOARCH=amd64 go build -o bin/playbuddy.exe cmd/main.go

# Build for Linux
build-linux:
	@echo "Building for Linux..."
	GOOS=linux GOARCH=amd64 go build -o bin/playbuddy cmd/main.go

# Build for macOS
build-macos:
	@echo "Building for macOS..."
	GOOS=darwin GOARCH=amd64 go build -o bin/playbuddy cmd/main.go 