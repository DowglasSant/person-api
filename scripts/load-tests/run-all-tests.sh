#!/bin/bash

echo "=========================================="
echo "Running ALL Load Tests"
echo "=========================================="
echo ""

# Check if server is running
if ! curl -s http://localhost:8080/health > /dev/null; then
    echo "ERROR: API is not running on http://localhost:8080"
    echo "Please start the API first: go run cmd/api/main.go"
    exit 1
fi

echo "API is running. Starting load tests..."
echo ""

# Run each test
bash scripts/load-tests/test-health.sh
sleep 2

bash scripts/load-tests/test-list-persons.sh
sleep 2

bash scripts/load-tests/test-find-by-cpf.sh
sleep 2

echo ""
echo "=========================================="
echo "All load tests completed!"
echo "=========================================="
echo ""
echo "Summary:"
echo "- Tested /health endpoint with up to 100 concurrent connections"
echo "- Tested /api/v1/persons (list) with up to 50 concurrent connections"
echo "- Tested /api/v1/persons/cpf/:cpf with up to 50 concurrent connections"
echo ""
echo "Review the results above for performance metrics."
echo "=========================================="
