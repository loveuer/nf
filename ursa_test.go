package ursa

import (
	"context"
	"net"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

// TestNotFoundHandling tests custom 404 handling
func TestNotFoundHandling(t *testing.T) {
	app := New()

	// Test default 404 handler
	req := httptest.NewRequest(http.MethodGet, "/nonexistent", nil)
	w := httptest.NewRecorder()
	app.ServeHTTP(w, req)
	if w.Code != 404 {
		t.Errorf("Expected status 404, got %d", w.Code)
	}
	if !contains(strings.ToLower(w.Body.String()), "not found") {
		t.Errorf("Expected body to contain 'not found', got '%s'", w.Body.String())
	}

	// Test custom 404 handler
	app = New(Config{
		NotFoundHandler: func(c *Ctx) error {
			return c.Status(404).JSON(Map{"error": "Custom 404", "message": "Resource not found"})
		},
	})

	app.Get("/test", func(c *Ctx) error {
		return c.SendString("test")
	})

	req = httptest.NewRequest(http.MethodGet, "/nonexistent", nil)
	w = httptest.NewRecorder()
	app.ServeHTTP(w, req)
	if w.Code != 404 {
		t.Errorf("Expected status 404, got %d", w.Code)
	}
	if !contains(w.Body.String(), "Custom 404") {
		t.Errorf("Expected body to contain 'Custom 404', got '%s'", w.Body.String())
	}
}

// TestMethodNotAllowed tests 405 Method Not Allowed handling
func TestMethodNotAllowed(t *testing.T) {
	app := New()

	app.Get("/test", func(c *Ctx) error {
		return c.SendString("GET test")
	})

	// Test POST to GET-only route
	req := httptest.NewRequest(http.MethodPost, "/test", nil)
	w := httptest.NewRecorder()
	app.ServeHTTP(w, req)
	if w.Code != 405 {
		t.Errorf("Expected status 405, got %d", w.Code)
	}
	if w.Body.String() != "405 Method Not Allowed" {
		t.Errorf("Expected body '405 Method Not Allowed', got '%s'", w.Body.String())
	}

	// Test custom 405 handler
	app = New(Config{
		MethodNotAllowedHandler: func(c *Ctx) error {
			return c.Status(405).JSON(Map{"error": "Method not allowed", "allowed": []string{"GET"}})
		},
	})

	app.Get("/test", func(c *Ctx) error {
		return c.SendString("GET test")
	})

	req = httptest.NewRequest(http.MethodPost, "/test", nil)
	w = httptest.NewRecorder()
	app.ServeHTTP(w, req)
	if w.Code != 405 {
		t.Errorf("Expected status 405, got %d", w.Code)
	}
	if !contains(w.Body.String(), "Method not allowed") {
		t.Errorf("Expected body to contain 'Method not allowed', got '%s'", w.Body.String())
	}
}

// TestPanicRecovery tests panic recovery functionality
func TestPanicRecovery(t *testing.T) {
	app := New()

	// Enable recovery (it's enabled by default)
	app.Use(NewRecover(false))

	app.Get("/panic", func(c *Ctx) error {
		panic("test panic")
	})

	app.Get("/normal", func(c *Ctx) error {
		return c.SendString("normal response")
	})

	// Test panic recovery
	req := httptest.NewRequest(http.MethodGet, "/panic", nil)
	w := httptest.NewRecorder()
	app.ServeHTTP(w, req)
	if w.Code != 500 {
		t.Errorf("Expected status 500, got %d", w.Code)
	}
	if w.Body.String() != "Internal Server Error" {
		t.Errorf("Expected body 'Internal Server Error', got '%s'", w.Body.String())
	}

	// Test that app still works after panic
	req = httptest.NewRequest(http.MethodGet, "/normal", nil)
	w = httptest.NewRecorder()
	app.ServeHTTP(w, req)
	if w.Code != 200 {
		t.Errorf("Expected status 200, got %d", w.Code)
	}
	if w.Body.String() != "normal response" {
		t.Errorf("Expected body 'normal response', got '%s'", w.Body.String())
	}
}

// TestAppConfiguration tests various app configurations
func TestAppConfiguration(t *testing.T) {
	// Test with custom configuration
	config := Config{
		DisableBanner:       true,
		DisableLogger:       true,
		DisableRecover:      true,
		BodyLimit:           1024,
		DisableHttpErrorLog: true,
	}

	app := New(config)

	app.Get("/test", func(c *Ctx) error {
		return c.SendString("test")
	})

	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	w := httptest.NewRecorder()
	app.ServeHTTP(w, req)
	if w.Code != 200 {
		t.Errorf("Expected status 200, got %d", w.Code)
	}
}

// TestGracefulShutdown tests graceful shutdown functionality
func TestGracefulShutdown(t *testing.T) {
	app := New()

	app.Get("/test", func(c *Ctx) error {
		return c.SendString("test")
	})

	// Create a test listener manually
	ln, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("Failed to create listener: %v", err)
	}

	// Start server in goroutine
	go func() {
		if err := app.RunListener(ln); err != nil && err != http.ErrServerClosed {
			t.Errorf("Server error: %v", err)
		}
	}()

	// Wait a bit for server to start
	time.Sleep(100 * time.Millisecond)

	// Create context with timeout for shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Test graceful shutdown
	err = app.Shutdown(ctx)
	if err != nil {
		t.Errorf("Expected no error during shutdown, got %v", err)
	}
}

// TestMultipleHandlers tests multiple handlers on the same route
func TestMultipleHandlers(t *testing.T) {
	app := New()

	var order []string

	app.Get("/test",
		func(c *Ctx) error {
			order = append(order, "handler1")
			return c.Next()
		},
		func(c *Ctx) error {
			order = append(order, "handler2")
			return c.Next()
		},
		func(c *Ctx) error {
			order = append(order, "handler3")
			return c.SendString("final")
		},
	)

	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	w := httptest.NewRecorder()
	app.ServeHTTP(w, req)

	if w.Code != 200 {
		t.Errorf("Expected status 200, got %d", w.Code)
	}
	if w.Body.String() != "final" {
		t.Errorf("Expected body 'final', got '%s'", w.Body.String())
	}

	expectedOrder := []string{"handler1", "handler2", "handler3"}
	for i, expected := range expectedOrder {
		if i >= len(order) || order[i] != expected {
			t.Errorf("Expected handler #%d to be '%s', got '%s'", i, expected, order[i])
		}
	}
}

// TestTraceId tests trace ID generation and propagation
func TestTraceId(t *testing.T) {
	app := New()

	app.Get("/test", func(c *Ctx) error {
		traceId := c.Request.Context().Value(TraceKey)
		if traceId == nil {
			t.Error("Expected trace ID to be set")
		}
		return c.JSON(Map{"traceId": traceId})
	})

	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	w := httptest.NewRecorder()
	app.ServeHTTP(w, req)

	if w.Code != 200 {
		t.Errorf("Expected status 200, got %d", w.Code)
	}

	// Check that trace ID header is set
	traceHeader := w.Header().Get(TraceKey)
	if traceHeader == "" {
		t.Errorf("Expected trace ID header to be set")
	}
}

// TestComplexRouteHandling tests complex routing scenarios
func TestComplexRouteHandling(t *testing.T) {
	app := New()

	// Various route types
	app.Get("/", func(c *Ctx) error {
		return c.SendString("home")
	})

	app.Get("/health", func(c *Ctx) error {
		return c.JSON(Map{"status": "ok"})
	})

	app.Post("/users", func(c *Ctx) error {
		return c.SendStatus(201)
	})

	app.Get("/users/:id", func(c *Ctx) error {
		id := c.Param("id")
		return c.JSON(Map{"userId": id})
	})

	app.Put("/users/:id", func(c *Ctx) error {
		id := c.Param("id")
		return c.JSON(Map{"updated": id})
	})

	app.Delete("/users/:id", func(c *Ctx) error {
		return c.SendStatus(204)
	})

	app.Get("/api/v1/*path", func(c *Ctx) error {
		return c.SendString("API v1")
	})

	// Test root route
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	w := httptest.NewRecorder()
	app.ServeHTTP(w, req)
	if w.Code != 200 || w.Body.String() != "home" {
		t.Errorf("Root route failed: code=%d, body=%s", w.Code, w.Body.String())
	}

	// Test health check
	req = httptest.NewRequest(http.MethodGet, "/health", nil)
	w = httptest.NewRecorder()
	app.ServeHTTP(w, req)
	if w.Code != 200 || !contains(w.Body.String(), "ok") {
		t.Errorf("Health check failed: code=%d, body=%s", w.Code, w.Body.String())
	}

	// Test POST with created status
	req = httptest.NewRequest(http.MethodPost, "/users", nil)
	w = httptest.NewRecorder()
	app.ServeHTTP(w, req)
	if w.Code != 201 {
		t.Errorf("Expected status 201, got %d", w.Code)
	}

	// Test parameter route
	req = httptest.NewRequest(http.MethodGet, "/users/123", nil)
	w = httptest.NewRecorder()
	app.ServeHTTP(w, req)
	if w.Code != 200 || !contains(w.Body.String(), "123") {
		t.Errorf("Parameter route failed: code=%d, body=%s", w.Code, w.Body.String())
	}

	// Test DELETE with no content
	req = httptest.NewRequest(http.MethodDelete, "/users/123", nil)
	w = httptest.NewRecorder()
	app.ServeHTTP(w, req)
	if w.Code != 204 {
		t.Errorf("Expected status 204, got %d", w.Code)
	}

	// Test wildcard
	req = httptest.NewRequest(http.MethodGet, "/api/v1/users", nil)
	w = httptest.NewRecorder()
	app.ServeHTTP(w, req)
	if w.Code != 200 || w.Body.String() != "API v1" {
		t.Errorf("Wildcard route failed: code=%d, body=%s", w.Code, w.Body.String())
	}
}

// TestGetRoutes tests the GetRoutes method
func TestGetRoutes(t *testing.T) {
	app := New()

	app.Get("/users", func(c *Ctx) error {
		return c.SendString("users")
	})

	app.Post("/users", func(c *Ctx) error {
		return c.SendString("create user")
	})

	app.Get("/users/:id", func(c *Ctx) error {
		return c.SendString("user detail")
	})

	routes := app.GetRoutes()
	if len(routes) != 3 {
		t.Errorf("Expected 3 routes, got %d", len(routes))
	}

	// Check route methods and paths
	methodPathSet := make(map[string]bool)
	for _, route := range routes {
		key := route.Method + ":" + route.Path
		methodPathSet[key] = true
	}

	expectedRoutes := []string{
		"GET:/users",
		"POST:/users",
		"GET:/users/:id",
	}

	for _, expected := range expectedRoutes {
		if !methodPathSet[expected] {
			t.Errorf("Expected route '%s' not found", expected)
		}
	}
}
