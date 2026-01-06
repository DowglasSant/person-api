#!/bin/bash

echo "=========================================="
echo "Load Test - List Persons Endpoint"
echo "=========================================="
echo ""

# Test with different concurrency levels and pagination
echo "Test 1: 500 requests, 10 concurrent - Page 1, default size"
ab -n 500 -c 10 "http://localhost:8080/api/v1/persons?page=1&page_size=10"

echo ""
echo "=========================================="
echo "Test 2: 1000 requests, 25 concurrent - Page 1, size 20"
ab -n 1000 -c 25 "http://localhost:8080/api/v1/persons?page=1&page_size=20"

echo ""
echo "=========================================="
echo "Test 3: 2000 requests, 50 concurrent - Sorted by name"
ab -n 2000 -c 50 "http://localhost:8080/api/v1/persons?sort=name&order=asc"

echo ""
echo "=========================================="
echo "Load test completed!"
echo "=========================================="
