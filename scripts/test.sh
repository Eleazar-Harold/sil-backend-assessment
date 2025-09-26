#!/bin/bash

# Test script for SIL Backend Assessment
# This script runs unit tests, integration tests, and e2e tests

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Configuration
TEST_DATABASE_URL="postgres://testuser:testpass@localhost:5432/test_db?sslmode=disable"
COVERAGE_THRESHOLD=80

echo -e "${BLUE}🧪 Starting test suite for SIL Backend Assessment${NC}"

# Function to run tests with coverage
run_tests_with_coverage() {
    local package=$1
    local test_name=$2
    
    echo -e "${BLUE}📊 Running ${test_name} with coverage...${NC}"
    
    go test -v -race -coverprofile=coverage.out -covermode=atomic ./${package}/...
    local exit_code=$?
    
    if [ $exit_code -eq 0 ]; then
        echo -e "${GREEN}✅ ${test_name} passed${NC}"
        
        # Calculate coverage
        local coverage=$(go tool cover -func=coverage.out | grep total | awk '{print $3}' | sed 's/%//')
        echo -e "${BLUE}📈 Coverage: ${coverage}%${NC}"
        
        # Check coverage threshold
        if (( $(echo "$coverage >= $COVERAGE_THRESHOLD" | bc -l) )); then
            echo -e "${GREEN}✅ Coverage meets threshold (${COVERAGE_THRESHOLD}%)${NC}"
        else
            echo -e "${RED}❌ Coverage below threshold (${COVERAGE_THRESHOLD}%)${NC}"
            exit 1
        fi
    else
        echo -e "${RED}❌ ${test_name} failed${NC}"
        exit $exit_code
    fi
}

# Function to run e2e tests
run_e2e_tests() {
    echo -e "${BLUE}🌐 Running E2E tests...${NC}"
    
    # Check if server is running
    if ! curl -s http://localhost:8080/api/health > /dev/null 2>&1; then
        echo -e "${YELLOW}⚠️  Server not running. Starting server for E2E tests...${NC}"
        go build -o server ./cmd/server
        ./server -config config.yaml > server.log 2>&1 &
        SERVER_PID=$!
        
        # Wait for server to start
        echo -e "${BLUE}⏳ Waiting for server to start...${NC}"
        for i in {1..30}; do
            if curl -s http://localhost:8080/api/health > /dev/null 2>&1; then
                echo -e "${GREEN}✅ Server started successfully${NC}"
                break
            fi
            sleep 1
        done
        
        if ! curl -s http://localhost:8080/api/health > /dev/null 2>&1; then
            echo -e "${RED}❌ Server failed to start${NC}"
            exit 1
        fi
    fi
    
    # Run E2E tests
    go test -v ./tests/e2e/...
    local exit_code=$?
    
    if [ $exit_code -eq 0 ]; then
        echo -e "${GREEN}✅ E2E tests passed${NC}"
    else
        echo -e "${RED}❌ E2E tests failed${NC}"
    fi
    
    # Cleanup server if we started it
    if [ ! -z "$SERVER_PID" ]; then
        echo -e "${BLUE}🧹 Cleaning up server...${NC}"
        kill $SERVER_PID || true
        rm -f server
    fi
    
    return $exit_code
}

# Function to run linting
run_linting() {
    echo -e "${BLUE}🔍 Running linting...${NC}"
    
    # Check if golangci-lint is installed
    if ! command -v golangci-lint > /dev/null 2>&1; then
        echo -e "${YELLOW}⚠️  golangci-lint not found. Installing...${NC}"
        curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(go env GOPATH)/bin v1.54.2
    fi
    
    golangci-lint run
    local exit_code=$?
    
    if [ $exit_code -eq 0 ]; then
        echo -e "${GREEN}✅ Linting passed${NC}"
    else
        echo -e "${RED}❌ Linting failed${NC}"
        exit $exit_code
    fi
}

# Function to run security scan
run_security_scan() {
    echo -e "${BLUE}🔒 Running security scan...${NC}"
    
    # Check if gosec is installed
    if ! command -v gosec > /dev/null 2>&1; then
        echo -e "${YELLOW}⚠️  gosec not found. Installing...${NC}"
        go install github.com/securecodewarrior/gosec/v2/cmd/gosec@latest
    fi
    
    gosec ./...
    local exit_code=$?
    
    if [ $exit_code -eq 0 ]; then
        echo -e "${GREEN}✅ Security scan passed${NC}"
    else
        echo -e "${RED}❌ Security scan found issues${NC}"
        exit $exit_code
    fi
}

# Main test execution
main() {
    # Parse command line arguments
    while [[ $# -gt 0 ]]; do
        case $1 in
            --unit-only)
                UNIT_ONLY=true
                shift
                ;;
            --e2e-only)
                E2E_ONLY=true
                shift
                ;;
            --no-coverage)
                NO_COVERAGE=true
                shift
                ;;
            --no-lint)
                NO_LINT=true
                shift
                ;;
            --no-security)
                NO_SECURITY=true
                shift
                ;;
            --help)
                echo "Usage: $0 [options]"
                echo "Options:"
                echo "  --unit-only     Run only unit tests"
                echo "  --e2e-only      Run only E2E tests"
                echo "  --no-coverage   Skip coverage check"
                echo "  --no-lint       Skip linting"
                echo "  --no-security   Skip security scan"
                echo "  --help          Show this help message"
                exit 0
                ;;
            *)
                echo -e "${RED}❌ Unknown option: $1${NC}"
                exit 1
                ;;
        esac
    done
    
    # Set environment variables for testing
    export TEST_DATABASE_URL="${TEST_DATABASE_URL}"
    export CGO_ENABLED=1
    
    # Run linting if not skipped
    if [ "$NO_LINT" != true ]; then
        run_linting
    fi
    
    # Run security scan if not skipped
    if [ "$NO_SECURITY" != true ]; then
        run_security_scan
    fi
    
    # Run unit tests if not e2e-only
    if [ "$E2E_ONLY" != true ]; then
        echo -e "${BLUE}🧪 Running unit tests...${NC}"
        
        # Test notification adapters
        run_tests_with_coverage "internal/adapters/notifications" "Notification adapter tests"
        
        # Test core services
        run_tests_with_coverage "internal/core/services" "Core service tests"
        
        # Test API handlers
        run_tests_with_coverage "internal/api/rest/handlers" "REST handler tests"
        
        # Test GraphQL resolvers
        run_tests_with_coverage "internal/api/graphql/resolvers" "GraphQL resolver tests"
        
        # Test repositories
        run_tests_with_coverage "internal/adapters/repositories" "Repository tests"
        
        # Test middleware
        run_tests_with_coverage "internal/adapters/middleware" "Middleware tests"
        
        # Test utilities
        run_tests_with_coverage "internal/testutils" "Test utility tests"
        
        echo -e "${GREEN}✅ All unit tests passed${NC}"
    fi
    
    # Run E2E tests if not unit-only
    if [ "$UNIT_ONLY" != true ]; then
        run_e2e_tests
    fi
    
    # Generate coverage report
    if [ "$NO_COVERAGE" != true ] && [ "$E2E_ONLY" != true ]; then
        echo -e "${BLUE}📊 Generating coverage report...${NC}"
        go tool cover -html=coverage.out -o coverage.html
        echo -e "${GREEN}✅ Coverage report generated: coverage.html${NC}"
    fi
    
    echo -e "${GREEN}🎉 All tests completed successfully!${NC}"
}

# Run main function
main "$@"
