# Development Guide

## Setting Up Development Environment

### Prerequisites
- Go 1.19 or higher
- Git

### Getting Started
1. Clone the repository:
```bash
git clone https://github.com/thewizardshell/froggit.git
cd froggit
```

2. Install dependencies:
```bash
go mod download
```

3. Run the application:
```bash
go run main.go
```

## Project Structure
- `main.go`: Application entry point
- `internal/`: Core packages
- `docs/`: Documentation

## Building
```bash
go build -o froggit
```

## Testing
```bash
go test ./...
```

## Contributing
See [Contributing Guidelines](contributing.md) for details on our code of conduct and process for submitting pull requests.