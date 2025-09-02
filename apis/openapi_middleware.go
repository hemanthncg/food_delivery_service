package apis

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/getkin/kin-openapi/openapi3"
)

var (
	openAPIDoc *openapi3.T
)

// InitOpenAPI initializes OpenAPI specification
func InitOpenAPI() error {
	// Try multiple possible paths for the OpenAPI spec
	specPaths := []string{
		"openapi.yaml",
		"../openapi.yaml",
		"../../openapi.yaml",
	}

	var doc *openapi3.T
	var err error
	loader := openapi3.NewLoader()

	for _, specPath := range specPaths {
		if _, statErr := os.Stat(specPath); statErr == nil {
			doc, err = loader.LoadFromFile(specPath)
			if err == nil {
				break
			}
		}
	}

	if doc == nil {
		return fmt.Errorf("OpenAPI spec not found in any of: %v", specPaths)
	}

	// Validate the specification
	if err := doc.Validate(loader.Context); err != nil {
		return fmt.Errorf("invalid OpenAPI spec: %v", err)
	}

	openAPIDoc = doc
	fmt.Println("OpenAPI specification loaded successfully")

	return nil
}

// GetOpenAPIDoc returns the OpenAPI document
func GetOpenAPIDoc() *openapi3.T {
	return openAPIDoc
}

// ServeOpenAPIDoc serves the OpenAPI specification
func ServeOpenAPIDoc(w http.ResponseWriter, r *http.Request) {
	if openAPIDoc == nil {
		http.Error(w, "OpenAPI specification not loaded", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(openAPIDoc)
}

// ServeOpenAPISwagger serves Swagger UI for the OpenAPI spec
func ServeOpenAPISwagger(w http.ResponseWriter, r *http.Request) {
	swaggerHTML := `<!DOCTYPE html>
<html>
<head>
    <title>Food Delivery Service API - Swagger UI</title>
    <link rel="stylesheet" type="text/css" href="https://unpkg.com/swagger-ui-dist@5.9.0/swagger-ui.css" />
    <style>
        html { box-sizing: border-box; overflow: -moz-scrollbars-vertical; overflow-y: scroll; }
        *, *:before, *:after { box-sizing: inherit; }
        body { margin:0; background: #fafafa; }
    </style>
</head>
<body>
    <div id="swagger-ui"></div>
    <script src="https://unpkg.com/swagger-ui-dist@5.9.0/swagger-ui-bundle.js"></script>
    <script src="https://unpkg.com/swagger-ui-dist@5.9.0/swagger-ui-standalone-preset.js"></script>
    <script>
        window.onload = function() {
            const ui = SwaggerUIBundle({
                url: '/openapi.json',
                dom_id: '#swagger-ui',
                deepLinking: true,
                presets: [
                    SwaggerUIBundle.presets.apis,
                    SwaggerUIStandalonePreset
                ],
                plugins: [
                    SwaggerUIBundle.plugins.DownloadUrl
                ],
                layout: "StandaloneLayout"
            });
        };
    </script>
</body>
</html>`

	w.Header().Set("Content-Type", "text/html")
	w.Write([]byte(swaggerHTML))
}
