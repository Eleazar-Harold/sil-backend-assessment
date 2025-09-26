package docs

import (
	"encoding/json"
	"io"
	"net/http"
	"os"
	"path/filepath"

	"github.com/uptrace/bunrouter"
)

// DocsHandler handles documentation endpoints
type DocsHandler struct {
	docsPath string
}

// NewDocsHandler creates a new documentation handler
func NewDocsHandler(docsPath string) *DocsHandler {
	return &DocsHandler{
		docsPath: docsPath,
	}
}

// ServeSwaggerJSON serves the OpenAPI/Swagger JSON specification
func (h *DocsHandler) ServeSwaggerJSON(w http.ResponseWriter, req bunrouter.Request) error {
	swaggerPath := filepath.Join(h.docsPath, "swagger.json")
	
	file, err := os.Open(swaggerPath)
	if err != nil {
		http.Error(w, "Swagger specification not found", http.StatusNotFound)
		return err
	}
	defer file.Close()

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	
	_, err = io.Copy(w, file)
	return err
}

// ServeRedocHTML serves the ReDoc HTML documentation
func (h *DocsHandler) ServeRedocHTML(w http.ResponseWriter, req bunrouter.Request) error {
	redocPath := filepath.Join(h.docsPath, "redoc.html")
	
	file, err := os.Open(redocPath)
	if err != nil {
		http.Error(w, "ReDoc documentation not found", http.StatusNotFound)
		return err
	}
	defer file.Close()

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	
	_, err = io.Copy(w, file)
	return err
}

// ServeSwaggerUI serves the Swagger UI HTML documentation
func (h *DocsHandler) ServeSwaggerUI(w http.ResponseWriter, req bunrouter.Request) error {
	swaggerUIPath := filepath.Join(h.docsPath, "swagger-ui.html")
	
	file, err := os.Open(swaggerUIPath)
	if err != nil {
		http.Error(w, "Swagger UI documentation not found", http.StatusNotFound)
		return err
	}
	defer file.Close()

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	
	_, err = io.Copy(w, file)
	return err
}

// ServeGraphQLDocs serves the GraphQL documentation
func (h *DocsHandler) ServeGraphQLDocs(w http.ResponseWriter, req bunrouter.Request) error {
	graphqlDocsPath := filepath.Join(h.docsPath, "graphql-docs.html")
	
	file, err := os.Open(graphqlDocsPath)
	if err != nil {
		http.Error(w, "GraphQL documentation not found", http.StatusNotFound)
		return err
	}
	defer file.Close()

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	
	_, err = io.Copy(w, file)
	return err
}

// ServeDocsIndex serves the main documentation index
func (h *DocsHandler) ServeDocsIndex(w http.ResponseWriter, req bunrouter.Request) error {
	response := map[string]interface{}{
		"title": "SIL Backend Assessment API Documentation",
		"description": "Comprehensive REST and GraphQL API documentation",
		"version": "1.0.0",
		"documentation": map[string]interface{}{
			"rest": map[string]string{
				"swagger_json": "/swagger.json",
				"swagger_ui":   "/swagger-ui.html",
				"redoc":        "/redoc.html",
			},
			"graphql": map[string]string{
				"documentation": "/graphql-docs.html",
				"playground":    "/graphql/playground",
				"endpoint":      "/graphql",
			},
		},
		"endpoints": map[string]string{
			"rest_api":    "/api",
			"graphql_api": "/graphql",
			"health":      "/api/health",
		},
		"authentication": map[string]string{
			"jwt_auth":  "Bearer token for user operations",
			"oidc_auth": "Bearer token for customer operations",
		},
		"links": map[string]string{
			"project_readme":     "README.md",
			"api_documentation":  "API_DOCUMENTATION.md",
			"quick_reference":    "ENDPOINTS_QUICK_REFERENCE.md",
			"oidc_setup":        "OIDC_SETUP.md",
		},
	}

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	
	return json.NewEncoder(w).Encode(response)
}

// RegisterRoutes registers documentation routes
func (h *DocsHandler) RegisterRoutes(router *bunrouter.Router) {
	// Swagger/OpenAPI endpoints
	router.GET("/swagger.json", h.ServeSwaggerJSON)
	router.GET("/swagger-ui.html", h.ServeSwaggerUI)
	router.GET("/redoc.html", h.ServeRedocHTML)
	
	// GraphQL documentation
	router.GET("/graphql-docs.html", h.ServeGraphQLDocs)
	
	// Main documentation index (override the existing /docs endpoint)
	router.GET("/docs", h.ServeDocsIndex)
	
	// Additional convenience routes
	router.GET("/docs/rest", h.ServeSwaggerUI)
	router.GET("/docs/graphql", h.ServeGraphQLDocs)
	router.GET("/docs/redoc", h.ServeRedocHTML)
}