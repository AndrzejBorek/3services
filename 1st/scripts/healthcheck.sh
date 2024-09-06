#!/bin/sh -xe

response=$(curl -s -o /dev/null -w "%{http_code}" -X GET localhost:8080/generate/1)

if [ "$response" -eq 200 ]; then
    echo "Healthcheck failed. Response: $response"
    exit 0
else
    echo "Healthcheck succeded. Response: $response"
    exit 1
fi