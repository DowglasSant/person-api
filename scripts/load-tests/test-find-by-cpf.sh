#!/bin/bash

echo "=========================================="
echo "Load Test - Find Person by CPF Endpoint"
echo "=========================================="
echo ""

# Test assuming CPF exists in database
CPF="11144477735"

echo "Test 1: 1000 requests, 10 concurrent connections"
ab -n 1000 -c 10 "http://localhost:8080/api/v1/persons/cpf/${CPF}"

echo ""
echo "=========================================="
echo "Test 2: 3000 requests, 30 concurrent connections"
ab -n 3000 -c 30 "http://localhost:8080/api/v1/persons/cpf/${CPF}"

echo ""
echo "=========================================="
echo "Test 3: 5000 requests, 50 concurrent connections"
ab -n 5000 -c 50 "http://localhost:8080/api/v1/persons/cpf/${CPF}"

echo ""
echo "=========================================="
echo "Note: Make sure the CPF ${CPF} exists in the database before running this test."
echo "=========================================="

echo ""
echo "Load test completed!"
