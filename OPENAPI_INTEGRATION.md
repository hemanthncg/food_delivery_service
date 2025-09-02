# OpenAPI Integration Guide

## Overview

This document describes the OpenAPI 3.0 integration for the Food Delivery Service API. 
The integration provides automatic request/response validation, comprehensive API 
documentation, and interactive testing capabilities.

## Features

- **OpenAPI 3.0 Specification**: Complete API documentation in YAML format
- **Request Validation**: Automatic validation of incoming requests against schemas
- **Response Validation**: Validation of API responses for consistency
- **Interactive Documentation**: Swagger UI for testing and exploration
- **Schema Validation**: JSON Schema validation for all data structures
- **Error Handling**: Standardized error responses with proper HTTP status codes

## API Endpoints

### Core Service Endpoints

| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/health` | Health check endpoint |
| POST | `/restaurant/enroll` | Register new restaurant |
| GET | `/restaurants` | List all restaurants |
| GET | `/menu` | Get restaurant menu |
| POST | `/driver/enroll` | Register new driver |
| GET | `/drivers` | List all drivers |
| GET/POST | `/driver/location` | Get/update driver location |
| POST | `/order` | Place new order |
| POST | `/order/delivered` | Mark order as delivered |
| GET | `/order/get` | Get order details |
| GET | `/orders` | List all orders |

### OpenAPI Documentation Endpoints

| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/openapi.json` | OpenAPI specification in JSON format |
| GET | `/docs` | Interactive Swagger UI |

## Getting Started

### 1. Prerequisites

Ensure you have the following dependencies:
- Go 1.25.0 or higher
- kin-openapi library (automatically installed)

### 2. Running the Service

```bash
# Start the service
go run main.go

# Or build and run
go build -o food-delivery-service
./food-delivery-service
```

The service will start on port 8080 with OpenAPI integration enabled.

### 3. Accessing Documentation

- **OpenAPI Spec**: http://localhost:8080/openapi.json
- **Swagger UI**: http://localhost:8080/docs

## OpenAPI Specification

### Schema Definitions

#### Location
```yaml
Location:
  type: object
  properties:
    lat:
      type: number
      format: double
      description: Latitude coordinate
    lng:
      type: number
      format: double
      description: Longitude coordinate
  required: [lat, lng]
```

#### Restaurant
```yaml
Restaurant:
  type: object
  properties:
    id:
      type: string
      description: Unique restaurant identifier
    name:
      type: string
      description: Restaurant name
    location:
      $ref: '#/components/schemas/Location'
  required: [id, name, location]
```

#### Order
```yaml
Order:
  type: object
  properties:
    id:
      type: string
      description: Unique order identifier
    restaurant_id:
      type: string
      description: Restaurant identifier
    items:
      type: array
      items:
        type: string
      description: Ordered food items
    user_location:
      $ref: '#/components/schemas/Location'
    driver_id:
      type: string
      description: Assigned driver identifier
    status:
      type: string
      enum: [created, preparing, out_for_delivery, delivered]
      description: Order status
  required: [id, restaurant_id, items, user_location]
```

## Validation Features

### Request Validation

The OpenAPI middleware automatically validates:
- **Path Parameters**: URL path variables
- **Query Parameters**: URL query string parameters
- **Request Headers**: HTTP headers including Content-Type
- **Request Body**: JSON payload validation against schemas
- **HTTP Methods**: Allowed HTTP methods per endpoint

### Response Validation

Response validation ensures:
- **Status Codes**: Valid HTTP status codes
- **Response Headers**: Proper header format
- **Response Body**: JSON schema compliance
- **Content Types**: Correct MIME types

### Error Handling

Standardized error responses:
```json
{
  "error": "Invalid request data",
  "code": "VALIDATION_ERROR"
}
```

## Testing with Swagger UI

### 1. Access Swagger UI
Navigate to http://localhost:8080/docs

### 2. Test Endpoints
- Expand any endpoint to see detailed information
- Click "Try it out" to test the API
- Fill in required parameters
- Execute the request
- View response details and validation

### 3. Example Requests

#### Enroll Restaurant
```bash
curl -X POST "http://localhost:8080/restaurant/enroll" \
  -H "Content-Type: application/json" \
  -d '{
    "id": "rest3",
    "name": "Sushi Bar",
    "location": {
      "lat": 19.0760,
      "lng": 72.8777
    },
    "menu": ["California Roll", "Salmon Nigiri"]
  }'
```

#### Place Order
```bash
curl -X POST "http://localhost:8080/order" \
  -H "Content-Type: application/json" \
  -d '{
    "id": "order123",
    "restaurant_id": "rest1",
    "items": ["Margherita Pizza"],
    "user_location": {
      "lat": 19.0760,
      "lng": 72.8777
    }
  }'
```

## Configuration

### OpenAPI Specification File
The service loads the OpenAPI specification from `openapi.yaml` in the root 
directory. This file defines:
- API endpoints and operations
- Request/response schemas
- Validation rules
- Documentation metadata

### Middleware Configuration
The OpenAPI middleware is automatically enabled when the service starts. 
If initialization fails, the service continues without validation but logs 
a warning.

## Development Workflow

### 1. Adding New Endpoints
1. Add the endpoint to your handler
2. Update the OpenAPI specification (`openapi.yaml`)
3. Add the route to `apis/routes.go`
4. Test with Swagger UI

### 2. Modifying Schemas
1. Update the schema in `openapi.yaml`
2. Ensure your Go structs match the schema
3. Test validation with sample requests

### 3. Custom Validation
The middleware supports custom validation through:
- Schema extensions
- Custom format validators
- Additional validation rules

## Troubleshooting

### Common Issues

#### OpenAPI Spec Not Found
```
Error: OpenAPI spec not found: openapi.yaml
```
**Solution**: Ensure `openapi.yaml` exists in the project root directory.

#### Validation Errors
```
Request validation failed: [error details]
```
**Solution**: Check request format against OpenAPI schemas. Use Swagger UI 
to see expected request format.

#### Response Validation Warnings
```
Response validation failed: [error details]
```
**Solution**: Check response format in your handlers. Ensure they match 
the defined schemas.

### Debug Mode

Enable detailed logging by checking the console output for:
- OpenAPI initialization status
- Validation error details
- Middleware execution logs

## Best Practices

### 1. Schema Design
- Use descriptive property names
- Include examples for complex schemas
- Define required fields explicitly
- Use enums for status fields

### 2. Error Handling
- Provide meaningful error messages
- Use appropriate HTTP status codes
- Include error codes for programmatic handling
- Log validation failures for debugging

### 3. Documentation
- Keep OpenAPI spec up to date
- Include comprehensive descriptions
- Provide realistic examples
- Document all possible responses

### 4. Testing
- Test all endpoints through Swagger UI
- Validate error scenarios
- Test edge cases and boundary conditions
- Verify response schemas

## Integration with CI/CD

### Automated Validation
Include OpenAPI validation in your CI/CD pipeline:
```bash
# Validate OpenAPI spec
npx @redocly/cli lint openapi.yaml

# Generate client SDKs
npx @openapitools/openapi-generator-cli generate \
  -i openapi.yaml \
  -g go \
  -o ./generated-client
```

### Documentation Generation
Automatically generate documentation:
```bash
# Generate HTML documentation
npx @redocly/cli build-docs openapi.yaml -o docs/

# Generate Postman collection
npx @apimatic/openapi-to-postman openapi.yaml
```

## Security Considerations

### API Key Authentication
The OpenAPI spec includes API key authentication:
```yaml
securitySchemes:
  ApiKeyAuth:
    type: apiKey
    in: header
    name: X-API-Key
```

### Input Validation
- All inputs are validated against schemas
- Malicious payloads are rejected
- Request size limits are enforced
- Content-Type validation prevents injection

## Performance Impact

### Validation Overhead
- Request validation: ~1-5ms per request
- Response validation: ~1-3ms per response
- Schema compilation: One-time cost at startup
- Memory usage: Minimal increase

### Optimization Tips
- Use efficient JSON parsing
- Minimize schema complexity
- Cache compiled schemas
- Disable validation in production if needed

## Conclusion

The OpenAPI integration provides a robust foundation for API development, 
testing, and documentation. It ensures API consistency, improves developer 
experience, and enables automated testing and client generation.

For more information, refer to:
- [OpenAPI 3.0 Specification](https://swagger.io/specification/)
- [kin-openapi Documentation](https://github.com/getkin/kin-openapi)
- [Swagger UI Documentation](https://swagger.io/tools/swagger-ui/)
