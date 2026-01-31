package ursa

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

// Helper function to check if a string contains a substring
func contains(s, substr string) bool {
	return strings.Contains(s, substr)
}

// TestHTTPMethods tests all HTTP methods (GET, POST, PUT, DELETE, PATCH, HEAD, OPTIONS)
func TestHTTPMethods(t *testing.T) {
	app := New()

	// Test GET method
	app.Get("/test", func(c *Ctx) error {
		return c.SendString("GET")
	})

	// Test POST method
	app.Post("/test", func(c *Ctx) error {
		return c.SendString("POST")
	})

	// Test PUT method
	app.Put("/test", func(c *Ctx) error {
		return c.SendString("PUT")
	})

	// Test DELETE method
	app.Delete("/test", func(c *Ctx) error {
		return c.SendString("DELETE")
	})

	// Test PATCH method
	app.Patch("/test", func(c *Ctx) error {
		return c.SendString("PATCH")
	})

	// Test HEAD method
	app.Head("/test", func(c *Ctx) error {
		return c.SendStatus(200)
	})

	// Test OPTIONS method
	app.Options("/test", func(c *Ctx) error {
		return c.SendString("OPTIONS")
	})

	// Test GET request
	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	w := httptest.NewRecorder()
	app.ServeHTTP(w, req)
	if w.Code != 200 {
		t.Errorf("Expected status 200, got %d", w.Code)
	}
	if w.Body.String() != "GET" {
		t.Errorf("Expected body 'GET', got '%s'", w.Body.String())
	}

	// Test POST request
	req = httptest.NewRequest(http.MethodPost, "/test", nil)
	w = httptest.NewRecorder()
	app.ServeHTTP(w, req)
	if w.Code != 200 {
		t.Errorf("Expected status 200, got %d", w.Code)
	}
	if w.Body.String() != "POST" {
		t.Errorf("Expected body 'POST', got '%s'", w.Body.String())
	}

	// Test PUT request
	req = httptest.NewRequest(http.MethodPut, "/test", nil)
	w = httptest.NewRecorder()
	app.ServeHTTP(w, req)
	if w.Code != 200 {
		t.Errorf("Expected status 200, got %d", w.Code)
	}
	if w.Body.String() != "PUT" {
		t.Errorf("Expected body 'PUT', got '%s'", w.Body.String())
	}

	// Test DELETE request
	req = httptest.NewRequest(http.MethodDelete, "/test", nil)
	w = httptest.NewRecorder()
	app.ServeHTTP(w, req)
	if w.Code != 200 {
		t.Errorf("Expected status 200, got %d", w.Code)
	}
	if w.Body.String() != "DELETE" {
		t.Errorf("Expected body 'DELETE', got '%s'", w.Body.String())
	}

	// Test PATCH request
	req = httptest.NewRequest(http.MethodPatch, "/test", nil)
	w = httptest.NewRecorder()
	app.ServeHTTP(w, req)
	if w.Code != 200 {
		t.Errorf("Expected status 200, got %d", w.Code)
	}
	if w.Body.String() != "PATCH" {
		t.Errorf("Expected body 'PATCH', got '%s'", w.Body.String())
	}

	// Test HEAD request
	req = httptest.NewRequest(http.MethodHead, "/test", nil)
	w = httptest.NewRecorder()
	app.ServeHTTP(w, req)
	if w.Code != 200 {
		t.Errorf("Expected status 200, got %d", w.Code)
	}

	// Test OPTIONS request
	req = httptest.NewRequest(http.MethodOptions, "/test", nil)
	w = httptest.NewRecorder()
	app.ServeHTTP(w, req)
	if w.Code != 200 {
		t.Errorf("Expected status 200, got %d", w.Code)
	}
	if w.Body.String() != "OPTIONS" {
		t.Errorf("Expected body 'OPTIONS', got '%s'", w.Body.String())
	}
}

// TestRouteParameters tests path parameters like :param
func TestRouteParameters(t *testing.T) {
	app := New()

	app.Get("/users/:id", func(c *Ctx) error {
		id := c.Param("id")
		return c.JSON(Map{"id": id})
	})

	app.Get("/users/:id/posts/:postId", func(c *Ctx) error {
		id := c.Param("id")
		postId := c.Param("postId")
		return c.JSON(Map{"userId": id, "postId": postId})
	})

	// Test single parameter
	req := httptest.NewRequest(http.MethodGet, "/users/123", nil)
	w := httptest.NewRecorder()
	app.ServeHTTP(w, req)
	if w.Code != 200 {
		t.Errorf("Expected status 200, got %d", w.Code)
	}
	if !contains(w.Body.String(), "123") {
		t.Errorf("Expected body to contain '123', got '%s'", w.Body.String())
	}

	// Test multiple parameters
	req = httptest.NewRequest(http.MethodGet, "/users/456/posts/789", nil)
	w = httptest.NewRecorder()
	app.ServeHTTP(w, req)
	if w.Code != 200 {
		t.Errorf("Expected status 200, got %d", w.Code)
	}
	if !contains(w.Body.String(), "456") {
		t.Errorf("Expected body to contain '456', got '%s'", w.Body.String())
	}
	if !contains(w.Body.String(), "789") {
		t.Errorf("Expected body to contain '789', got '%s'", w.Body.String())
	}
}

// TestWildcardRoutes tests wildcard routes like *wildcard
func TestWildcardRoutes(t *testing.T) {
	app := New()

	app.Get("/files/*filepath", func(c *Ctx) error {
		filepath := c.Param("filepath")
		return c.JSON(Map{"filepath": filepath})
	})

	app.Get("/api/*path", func(c *Ctx) error {
		return c.SendString("API endpoint")
	})

	// Test wildcard parameter
	req := httptest.NewRequest(http.MethodGet, "/files/path/to/file.txt", nil)
	w := httptest.NewRecorder()
	app.ServeHTTP(w, req)
	if w.Code != 200 {
		t.Errorf("Expected status 200, got %d", w.Code)
	}
	if !contains(w.Body.String(), "path/to/file.txt") {
		t.Errorf("Expected body to contain 'path/to/file.txt', got '%s'", w.Body.String())
	}

	// Test wildcard with empty path
	req = httptest.NewRequest(http.MethodGet, "/files/", nil)
	w = httptest.NewRecorder()
	app.ServeHTTP(w, req)
	if w.Code != 200 {
		t.Errorf("Expected status 200, got %d", w.Code)
	}

	// Test API wildcard
	req = httptest.NewRequest(http.MethodGet, "/api/v1/users", nil)
	w = httptest.NewRecorder()
	app.ServeHTTP(w, req)
	if w.Code != 200 {
		t.Errorf("Expected status 200, got %d", w.Code)
	}
	if w.Body.String() != "API endpoint" {
		t.Errorf("Expected body 'API endpoint', got '%s'", w.Body.String())
	}
}

// TestRouteGroups tests route group functionality
func TestRouteGroups(t *testing.T) {
	app := New()

	// Create a route group
	v1 := app.Group("/api/v1")
	v1.Get("/users", func(c *Ctx) error {
		return c.SendString("v1 users")
	})
	v1.Post("/users", func(c *Ctx) error {
		return c.SendString("v1 create user")
	})

	// Create another route group
	v2 := app.Group("/api/v2")
	v2.Get("/users", func(c *Ctx) error {
		return c.SendString("v2 users")
	})

	// Test v1 group
	req := httptest.NewRequest(http.MethodGet, "/api/v1/users", nil)
	w := httptest.NewRecorder()
	app.ServeHTTP(w, req)
	if w.Code != 200 {
		t.Errorf("Expected status 200, got %d", w.Code)
	}
	if w.Body.String() != "v1 users" {
		t.Errorf("Expected body 'v1 users', got '%s'", w.Body.String())
	}

	// Test v1 POST
	req = httptest.NewRequest(http.MethodPost, "/api/v1/users", nil)
	w = httptest.NewRecorder()
	app.ServeHTTP(w, req)
	if w.Code != 200 {
		t.Errorf("Expected status 200, got %d", w.Code)
	}
	if w.Body.String() != "v1 create user" {
		t.Errorf("Expected body 'v1 create user', got '%s'", w.Body.String())
	}

	// Test v2 group
	req = httptest.NewRequest(http.MethodGet, "/api/v2/users", nil)
	w = httptest.NewRecorder()
	app.ServeHTTP(w, req)
	if w.Code != 200 {
		t.Errorf("Expected status 200, got %d", w.Code)
	}
	if w.Body.String() != "v2 users" {
		t.Errorf("Expected body 'v2 users', got '%s'", w.Body.String())
	}
}

// TestNestedGroups tests nested route groups
func TestNestedGroups(t *testing.T) {
	app := New()

	// Create nested groups
	api := app.Group("/api")
	v1 := api.Group("/v1")
	users := v1.Group("/users")

	users.Get("/", func(c *Ctx) error {
		return c.SendString("v1 users list")
	})

	users.Get("/:id", func(c *Ctx) error {
		id := c.Param("id")
		return c.JSON(Map{"id": id, "version": "v1"})
	})

	// Test nested group route
	req := httptest.NewRequest(http.MethodGet, "/api/v1/users", nil)
	w := httptest.NewRecorder()
	app.ServeHTTP(w, req)
	if w.Code != 200 {
		t.Errorf("Expected status 200, got %d", w.Code)
	}
	if w.Body.String() != "v1 users list" {
		t.Errorf("Expected body 'v1 users list', got '%s'", w.Body.String())
	}

	// Test nested group with parameter
	req = httptest.NewRequest(http.MethodGet, "/api/v1/users/123", nil)
	w = httptest.NewRecorder()
	app.ServeHTTP(w, req)
	if w.Code != 200 {
		t.Errorf("Expected status 200, got %d", w.Code)
	}
	if !contains(w.Body.String(), "123") {
		t.Errorf("Expected body to contain '123', got '%s'", w.Body.String())
	}
	if !contains(w.Body.String(), "v1") {
		t.Errorf("Expected body to contain 'v1', got '%s'", w.Body.String())
	}
}

// TestRoutePriority tests route matching priority
func TestRoutePriority(t *testing.T) {
	app := New()

	// Static route should have higher priority than parameter route
	app.Get("/users/special", func(c *Ctx) error {
		return c.SendString("special user")
	})

	app.Get("/users/:id", func(c *Ctx) error {
		id := c.Param("id")
		return c.JSON(Map{"id": id})
	})

	// Test static route (should match first)
	req := httptest.NewRequest(http.MethodGet, "/users/special", nil)
	w := httptest.NewRecorder()
	app.ServeHTTP(w, req)
	if w.Code != 200 {
		t.Errorf("Expected status 200, got %d", w.Code)
	}
	if w.Body.String() != "special user" {
		t.Errorf("Expected body 'special user', got '%s'", w.Body.String())
	}

	// Test parameter route
	req = httptest.NewRequest(http.MethodGet, "/users/123", nil)
	w = httptest.NewRecorder()
	app.ServeHTTP(w, req)
	if w.Code != 200 {
		t.Errorf("Expected status 200, got %d", w.Code)
	}
	if !contains(w.Body.String(), "123") {
		t.Errorf("Expected body to contain '123', got '%s'", w.Body.String())
	}
}

// TestRouteConflicts tests route conflict handling
func TestRouteConflicts(t *testing.T) {
	app := New()

	// This should work fine
	app.Get("/test", func(c *Ctx) error {
		return c.SendString("test")
	})

	// Test that duplicate routes panic
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("Expected panic for duplicate route")
		}
	}()
	app.Get("/test", func(c *Ctx) error {
		return c.SendString("duplicate")
	})
}

// TestTrailingSlash tests trailing slash handling
func TestTrailingSlash(t *testing.T) {
	app := New()

	app.Get("/api", func(c *Ctx) error {
		return c.SendString("api without slash")
	})

	app.Get("/api-with-slash-suffix", func(c *Ctx) error {
		return c.SendString("api with suffix")
	})

	// Test without trailing slash
	req := httptest.NewRequest(http.MethodGet, "/api", nil)
	w := httptest.NewRecorder()
	app.ServeHTTP(w, req)
	if w.Code != 200 {
		t.Errorf("Expected status 200, got %d", w.Code)
	}
	if w.Body.String() != "api without slash" {
		t.Errorf("Expected body 'api without slash', got '%s'", w.Body.String())
	}

	// Test with different suffix
	req = httptest.NewRequest(http.MethodGet, "/api-with-slash-suffix", nil)
	w = httptest.NewRecorder()
	app.ServeHTTP(w, req)
	if w.Code != 200 {
		t.Errorf("Expected status 200, got %d", w.Code)
	}
	if w.Body.String() != "api with suffix" {
		t.Errorf("Expected body 'api with suffix', got '%s'", w.Body.String())
	}
}

// TestAnyMethod tests the Any method functionality
func TestAnyMethod(t *testing.T) {
	app := New()

	app.Any("/any", func(c *Ctx) error {
		return c.JSON(Map{"method": c.Method()})
	})

	methods := []string{http.MethodGet, http.MethodPost, http.MethodPut, http.MethodDelete, http.MethodPatch}

	for _, method := range methods {
		req := httptest.NewRequest(method, "/any", nil)
		w := httptest.NewRecorder()
		app.ServeHTTP(w, req)
		if w.Code != 200 {
			t.Errorf("Expected status 200, got %d", w.Code)
		}
		if !contains(w.Body.String(), method) {
			t.Errorf("Expected body to contain '%s', got '%s'", method, w.Body.String())
		}
	}
}

// TestMatchMethod tests the Match method functionality
func TestMatchMethod(t *testing.T) {
	app := New()

	app.Match([]string{http.MethodGet, http.MethodPost}, "/match", func(c *Ctx) error {
		return c.JSON(Map{"method": c.Method()})
	})

	// Test GET
	req := httptest.NewRequest(http.MethodGet, "/match", nil)
	w := httptest.NewRecorder()
	app.ServeHTTP(w, req)
	if w.Code != 200 {
		t.Errorf("Expected status 200, got %d", w.Code)
	}
	if !contains(w.Body.String(), "GET") {
		t.Errorf("Expected body to contain 'GET', got '%s'", w.Body.String())
	}

	// Test POST
	req = httptest.NewRequest(http.MethodPost, "/match", nil)
	w = httptest.NewRecorder()
	app.ServeHTTP(w, req)
	if w.Code != 200 {
		t.Errorf("Expected status 200, got %d", w.Code)
	}
	if !contains(w.Body.String(), "POST") {
		t.Errorf("Expected body to contain 'POST', got '%s'", w.Body.String())
	}

	// Test PUT (should not match)
	req = httptest.NewRequest(http.MethodPut, "/match", nil)
	w = httptest.NewRecorder()
	app.ServeHTTP(w, req)
	if w.Code != 405 {
		t.Errorf("Expected status 405, got %d", w.Code)
	}
}
