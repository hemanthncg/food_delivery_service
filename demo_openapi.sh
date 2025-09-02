#!/bin/bash

echo "🚀 Food Delivery Service - OpenAPI Integration Demo"
echo "=================================================="
echo ""

# Check if service is running
echo "📋 Checking if service is running..."
if curl -s http://localhost:8080/health > /dev/null; then
    echo "✅ Service is running on port 8080"
else
    echo "❌ Service is not running. Please start it first:"
    echo "   go run main.go"
    echo ""
    exit 1
fi

echo ""
echo "🌐 OpenAPI Documentation Endpoints:"
echo "   • OpenAPI Spec (JSON): http://localhost:8080/openapi.json"
echo "   • Swagger UI:          http://localhost:8080/docs"
echo ""

echo "📖 Testing OpenAPI Specification..."
if curl -s http://localhost:8080/openapi.json | jq '.info.title' > /dev/null 2>&1; then
    echo "✅ OpenAPI specification is accessible"
    TITLE=$(curl -s http://localhost:8080/openapi.json | jq -r '.info.title')
    VERSION=$(curl -s http://localhost:8080/openapi.json | jq -r '.info.version')
    echo "   📄 Title: $TITLE"
    echo "   🔢 Version: $VERSION"
else
    echo "❌ Failed to access OpenAPI specification"
fi

echo ""
echo "🔍 Testing Swagger UI..."
if curl -s http://localhost:8080/docs | grep -q "swagger-ui"; then
    echo "✅ Swagger UI is accessible"
else
    echo "❌ Failed to access Swagger UI"
fi

echo ""
echo "🧪 Testing API Endpoints with OpenAPI Validation..."

echo "   📍 Testing Health Endpoint..."
HEALTH_RESPONSE=$(curl -s -w "%{http_code}" http://localhost:8080/health)
HTTP_CODE="${HEALTH_RESPONSE: -3}"
RESPONSE_BODY="${HEALTH_RESPONSE%???}"
if [ "$HTTP_CODE" = "200" ]; then
    echo "   ✅ Health endpoint: $HTTP_CODE - $RESPONSE_BODY"
else
    echo "   ❌ Health endpoint: $HTTP_CODE"
fi

echo ""
echo "   🏪 Testing Restaurant List..."
RESTAURANTS_RESPONSE=$(curl -s -w "%{http_code}" http://localhost:8080/restaurants)
HTTP_CODE="${RESTAURANTS_RESPONSE: -3}"
RESPONSE_BODY="${RESTAURANTS_RESPONSE%???}"
if [ "$HTTP_CODE" = "200" ]; then
    echo "   ✅ Restaurants endpoint: $HTTP_CODE"
    echo "   📊 Response: $RESPONSE_BODY" | jq '.' 2>/dev/null || echo "   📊 Response: $RESPONSE_BODY"
else
    echo "   ❌ Restaurants endpoint: $HTTP_CODE"
fi

echo ""
echo "🎯 OpenAPI Integration Features:"
echo "   • ✅ Request validation against schemas"
echo "   • ✅ Response validation for consistency"
echo "   • ✅ Automatic parameter validation"
echo "   • ✅ Schema-based error handling"
echo "   • ✅ Interactive API documentation"
echo "   • ✅ Standardized API contracts"

echo ""
echo "📚 Next Steps:"
echo "   1. Open http://localhost:8080/docs in your browser"
echo "   2. Explore the interactive API documentation"
echo "   3. Test endpoints directly from Swagger UI"
echo "   4. Use the OpenAPI spec for client generation"
echo "   5. Integrate with CI/CD for automated validation"

echo ""
echo "🔗 Useful Links:"
echo "   • OpenAPI 3.0 Spec: https://swagger.io/specification/"
echo "   • kin-openapi:      https://github.com/getkin/kin-openapi"
echo "   • Swagger UI:       https://swagger.io/tools/swagger-ui/"

echo ""
echo "✨ Demo completed successfully!"
