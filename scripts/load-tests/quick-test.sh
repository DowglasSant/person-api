#!/bin/bash

echo "=========================================="
echo "Quick Load Test"
echo "=========================================="
echo ""

# Check if server is running
if ! curl -s http://localhost:8080/health > /dev/null; then
    echo "ERROR: API is not running on http://localhost:8080"
    echo ""
    echo "To start the API, run in another terminal:"
    echo "  go run cmd/api/main.go"
    echo ""
    echo "Or use this one-liner to start and background it:"
    echo "  go run cmd/api/main.go > /tmp/api.log 2>&1 &"
    echo ""
    exit 1
fi

echo "✓ API is running"
echo ""

# Quick test on health endpoint
echo "Running quick test on /health endpoint..."
echo "Test: 1000 requests, 20 concurrent connections"
echo ""

ab -n 1000 -c 20 -q http://localhost:8080/health

echo ""
echo "=========================================="
echo "Quick test completed!"
echo ""
echo "To run comprehensive tests, use:"
echo "  bash scripts/load-tests/run-all-tests.sh"
echo ""
echo "To run individual tests:"
echo "  bash scripts/load-tests/test-health.sh"
echo "  bash scripts/load-tests/test-list-persons.sh"
echo "  bash scripts/load-tests/test-find-by-cpf.sh"
echo "=========================================="
