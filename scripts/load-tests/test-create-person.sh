#!/bin/bash

echo "=========================================="
echo "Load Test - Create Person Endpoint"
echo "=========================================="
echo ""

# Create a JSON file with sample data
cat > /tmp/person-payload.json <<EOF
{
  "name": "Load Test User",
  "cpf": "11144477735",
  "birth_date": "1990-01-01T00:00:00Z",
  "phone": "81 91234-5678",
  "email": "loadtest@example.com"
}
EOF

echo "Test 1: 100 requests, 5 concurrent connections"
ab -n 100 -c 5 -p /tmp/person-payload.json -T application/json http://localhost:8080/api/v1/persons

echo ""
echo "=========================================="
echo "Test 2: 200 requests, 10 concurrent connections"
ab -n 200 -c 10 -p /tmp/person-payload.json -T application/json http://localhost:8080/api/v1/persons

echo ""
echo "=========================================="
echo "Note: This test will create duplicate CPF entries and may fail due to database constraints."
echo "For realistic testing, you may need to populate the database with test data first."
echo "=========================================="

# Cleanup
rm -f /tmp/person-payload.json

echo ""
echo "Load test completed!"
