package apis

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestInitOpenAPI(t *testing.T) {
	// Test successful initialization
	err := InitOpenAPI()
	if err != nil {
		t.Errorf("Expected no error, got: %v", err)
	}

	// Verify OpenAPI doc is loaded
	if openAPIDoc == nil {
		t.Error("Expected OpenAPI document to be loaded")
	}

	// Verify router is created
	if router == nil {
		t.Error("Expected router to be created")
	}
}

func TestServeOpenAPIDoc(t *testing.T) {
	// Initialize OpenAPI first
	err := InitOpenAPI()
	if err != nil {
		t.Fatalf("Failed to initialize OpenAPI: %v", err)
	}

	// Create test request
	req := httptest.NewRequest("GET", "/openapi.json", nil)
	w := httptest.NewRecorder()

	// Call handler
	ServeOpenAPIDoc(w, req)

	// Check response
	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}

	contentType := w.Header().Get("Content-Type")
	if contentType != "application/json" {
		t.Errorf("Expected Content-Type application/json, got %s", contentType)
	}

	// Verify response body is not empty
	if len(w.Body.Bytes()) == 0 {
		t.Error("Expected non-empty response body")
	}
}

func TestServeOpenAPISwagger(t *testing.T) {
	// Create test request
	req := httptest.NewRequest("GET", "/docs", nil)
	w := httptest.NewRecorder()

	// Call handler
	ServeOpenAPISwagger(w, req)

	// Check response
	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}

	contentType := w.Header().Get("Content-Type")
	if contentType != "text/html" {
		t.Errorf("Expected Content-Type text/html, got %s", contentType)
	}

	// Verify response body contains Swagger UI
	body := w.Body.String()
	if len(body) == 0 {
		t.Error("Expected non-empty response body")
	}

	if !contains(body, "swagger-ui") {
		t.Error("Expected response to contain Swagger UI")
	}
}

func TestOpenAPIMiddleware(t *testing.T) {
	// Initialize OpenAPI first
	err := InitOpenAPI()
	if err != nil {
		t.Fatalf("Failed to initialize OpenAPI: %v", err)
	}

	// Create a simple test handler
	testHandler := func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"message": "test"}`))
	}

	// Create middleware
	middleware := OpenAPIMiddleware(testHandler)

	// Test with valid request
	req := httptest.NewRequest("GET", "/health", nil)
	w := httptest.NewRecorder()

	middleware.ServeHTTP(w, req)

	// Should pass through to handler
	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}
}

// Helper function to check if string contains substring
func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr ||
		(len(s) > len(substr) &&
			(s[:len(substr)] == substr ||
				s[len(s)-len(substr):] == substr ||
				containsSubstring(s, substr))))
}

func containsSubstring(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
