#!/bin/bash

echo "=========================================="
echo "Load Test - Health Endpoint"
echo "=========================================="
echo ""

# Test with different concurrency levels
echo "Test 1: 1000 requests, 10 concurrent connections"
ab -n 1000 -c 10 http://localhost:8080/health

echo ""
echo "=========================================="
echo "Test 2: 5000 requests, 50 concurrent connections"
ab -n 5000 -c 50 http://localhost:8080/health

echo ""
echo "=========================================="
echo "Test 3: 10000 requests, 100 concurrent connections"
ab -n 10000 -c 100 http://localhost:8080/health

echo ""
echo "=========================================="
echo "Load test completed!"
echo "=========================================="
