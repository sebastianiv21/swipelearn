#!/bin/bash

# Test Runner for SwipeLearn API
# This script provides easy commands for running different types of tests

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Function to print colored output
print_status() {
    echo -e "${GREEN}$1${NC}"
}

print_warning() {
    echo -e "${YELLOW}$1${NC}"
}

print_error() {
    echo -e "${RED}$1${NC}"
}

# Function to run tests with appropriate command
run_tests() {
    local test_type=$1
    local description=$2
    
    print_status "🧪 Running $description..."
    
    case $test_type in
        "unit")
            echo "Running unit tests (services & repositories)..."
            go test -v ./internal/services ./internal/repositories
            ;;
        "services")
            echo "Running service tests only..."
            go test -v ./internal/services
            ;;
        "repositories")
            echo "Running repository tests only..."
            go test -v ./internal/repositories
            ;;
        "coverage")
            echo "Running tests with coverage..."
            go test -v -coverprofile=bin/coverage.out ./...
            go tool cover -html=bin/coverage.out -o bin/coverage.html
            print_status "📄 Coverage report generated: bin/coverage.html"
            ;;
        "integration")
            echo "Running integration tests (not implemented yet)..."
            print_warning "Integration tests are not yet implemented"
            ;;
        "all")
            echo "Running all available tests..."
            run_tests "unit" "Unit Tests"
            run_tests "integration" "Integration Tests"
            ;;
        *)
            print_error "Unknown test type: $1"
            echo "Available types: unit, services, repositories, coverage, integration, all"
            exit 1
            ;;
    esac
    
    if [ $? -eq 0 ]; then
        print_status "✅ Tests completed successfully!"
    else
        print_error "❌ Tests failed!"
        exit 1
    fi
}

# Parse command line arguments
case ${1:-unit} in
    "unit"|"")
        run_tests "unit" "Unit Tests"
        ;;
    "services")
        run_tests "services" "Service Tests"
        ;;
    "repositories")
        run_tests "repositories" "Repository Tests"
        ;;
    "coverage")
        run_tests "coverage" "Tests with Coverage"
        ;;
    "integration")
        run_tests "integration" "Integration Tests"
        ;;
    "all")
        run_tests "all" "All Tests"
        ;;
    "help"|"-h"|"--help")
        echo "SwipeLearn API Test Runner"
        echo ""
        echo "Usage: $0 [test_type]"
        echo ""
        echo "Available test types:"
        echo "  unit         Run unit tests (services & repositories) - default"
        echo "  services     Run service tests only"
        echo "  repositories Run repository tests only"
        echo "  coverage     Run tests with coverage report"
        echo "  integration  Run integration tests"
        echo "  all          Run all available tests"
        echo "  help         Show this help message"
        echo ""
        echo "Examples:"
        echo "  $0              # Run unit tests (default)"
        echo "  $0 unit         # Run unit tests"
        echo "  $0 coverage      # Run tests with coverage"
        echo "  $0 all           # Run all tests"
        ;;
    *)
        print_error "Unknown command: $1"
        echo "Use '$0 help' for usage information"
        exit 1
        ;;
esac