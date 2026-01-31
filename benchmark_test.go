package ursa

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

// BenchmarkRouteLookup benchmarks route lookup performance
func BenchmarkRouteLookup(b *testing.B) {
	app := New()

	// Register many routes with unique names
	for i := 0; i < 100; i++ {
		routeName := fmt.Sprintf("/route%d", i)
		app.Get(routeName, func(c *Ctx) error {
			return c.SendString("ok")
		})
	}

	// Register parameter routes
	for i := 0; i < 50; i++ {
		routeName := fmt.Sprintf("/param%d/:id", i)
		app.Get(routeName, func(c *Ctx) error {
			return c.SendString("ok")
		})
	}

	// Benchmark static route lookup
	req := httptest.NewRequest(http.MethodGet, "/route50", nil)
	w := httptest.NewRecorder()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		app.ServeHTTP(w, req)
	}
}

// BenchmarkParameterExtraction benchmarks parameter extraction performance
func BenchmarkParameterExtraction(b *testing.B) {
	app := New()

	app.Get("/users/:id/posts/:postId/comments/:commentId", func(c *Ctx) error {
		_ = c.Param("id")
		_ = c.Param("postId")
		_ = c.Param("commentId")
		return c.SendString("ok")
	})

	req := httptest.NewRequest(http.MethodGet, "/users/123/posts/456/comments/789", nil)
	w := httptest.NewRecorder()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		app.ServeHTTP(w, req)
	}
}

// BenchmarkMiddlewareChain benchmarks middleware chain performance
func BenchmarkMiddlewareChain(b *testing.B) {
	app := New()

	// Add multiple middleware
	for i := 0; i < 10; i++ {
		app.Use(func(c *Ctx) error {
			return c.Next()
		})
	}

	app.Get("/test", func(c *Ctx) error {
		return c.SendString("ok")
	})

	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	w := httptest.NewRecorder()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		app.ServeHTTP(w, req)
	}
}

// BenchmarkJSONResponse benchmarks JSON response performance
func BenchmarkJSONResponse(b *testing.B) {
	app := New()

	app.Get("/json", func(c *Ctx) error {
		return c.JSON(Map{
			"message": "Hello, World!",
			"status":  "ok",
			"data": Map{
				"id":   123,
				"name": "Test User",
				"tags": []string{"golang", "web", "api"},
			},
		})
	})

	req := httptest.NewRequest(http.MethodGet, "/json", nil)
	w := httptest.NewRecorder()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		app.ServeHTTP(w, req)
	}
}

// BenchmarkStringResponse benchmarks string response performance
func BenchmarkStringResponse(b *testing.B) {
	app := New()

	app.Get("/string", func(c *Ctx) error {
		return c.SendString("Hello, World! This is a test response string.")
	})

	req := httptest.NewRequest(http.MethodGet, "/string", nil)
	w := httptest.NewRecorder()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		app.ServeHTTP(w, req)
	}
}

// BenchmarkContextLocals benchmarks context locals performance
func BenchmarkContextLocals(b *testing.B) {
	app := New()

	app.Use(func(c *Ctx) error {
		c.Locals("user", "john")
		c.Locals("role", "admin")
		c.Locals("session", "abc123")
		return c.Next()
	})

	app.Get("/test", func(c *Ctx) error {
		_ = c.Locals("user")
		_ = c.Locals("role")
		_ = c.Locals("session")
		return c.SendString("ok")
	})

	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	w := httptest.NewRecorder()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		app.ServeHTTP(w, req)
	}
}

// BenchmarkWildcardRoutes benchmarks wildcard route performance
func BenchmarkWildcardRoutes(b *testing.B) {
	app := New()

	app.Get("/files/*filepath", func(c *Ctx) error {
		_ = c.Param("filepath")
		return c.SendString("ok")
	})

	req := httptest.NewRequest(http.MethodGet, "/files/path/to/very/long/file/name.txt", nil)
	w := httptest.NewRecorder()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		app.ServeHTTP(w, req)
	}
}

// BenchmarkRouteRegistration benchmarks route registration performance
func BenchmarkRouteRegistration(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		app := New()
		for j := 0; j < 100; j++ {
			app.Get("/route/"+string(rune(j)), func(c *Ctx) error {
				return c.SendString("ok")
			})
		}
	}
}

// BenchmarkGroupRoutes benchmarks route group performance
func BenchmarkGroupRoutes(b *testing.B) {
	app := New()

	api := app.Group("/api/v1")
	users := api.Group("/users")
	posts := users.Group("/posts")

	posts.Get("/test", func(c *Ctx) error {
		return c.SendString("ok")
	})

	req := httptest.NewRequest(http.MethodGet, "/api/v1/users/posts/test", nil)
	w := httptest.NewRecorder()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		app.ServeHTTP(w, req)
	}
}

// BenchmarkNotFound benchmarks 404 handling performance
func BenchmarkNotFound(b *testing.B) {
	app := New()

	app.Get("/test", func(c *Ctx) error {
		return c.SendString("ok")
	})

	req := httptest.NewRequest(http.MethodGet, "/nonexistent", nil)
	w := httptest.NewRecorder()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		app.ServeHTTP(w, req)
	}
}

// BenchmarkMethodNotAllowed benchmarks 405 handling performance
func BenchmarkMethodNotAllowed(b *testing.B) {
	app := New()

	app.Get("/test", func(c *Ctx) error {
		return c.SendString("ok")
	})

	req := httptest.NewRequest(http.MethodPost, "/test", nil)
	w := httptest.NewRecorder()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		app.ServeHTTP(w, req)
	}
}

// BenchmarkComplexRouting benchmarks complex routing scenarios
func BenchmarkComplexRouting(b *testing.B) {
	app := New()

	// Register various route types
	app.Get("/", func(c *Ctx) error { return c.SendString("home") })
	app.Get("/health", func(c *Ctx) error { return c.SendString("ok") })
	app.Get("/users", func(c *Ctx) error { return c.SendString("users") })
	app.Post("/users", func(c *Ctx) error { return c.SendStatus(201) })
	app.Get("/users/:id", func(c *Ctx) error { return c.SendString("user") })
	app.Put("/users/:id", func(c *Ctx) error { return c.SendString("updated") })
	app.Delete("/users/:id", func(c *Ctx) error { return c.SendStatus(204) })
	app.Get("/api/v1/*", func(c *Ctx) error { return c.SendString("api") })

	// Test different routes in sequence
	routes := []string{
		"/",
		"/health",
		"/users",
		"/users/123",
		"/api/v1/test",
	}

	reqs := make([]*http.Request, len(routes))
	for i, route := range routes {
		method := http.MethodGet
		if route == "/users" {
			method = http.MethodPost
		} else if route == "/users/123" {
			method = http.MethodPut
		}
		reqs[i] = httptest.NewRequest(method, route, nil)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		reqIndex := i % len(reqs)
		w := httptest.NewRecorder()
		app.ServeHTTP(w, reqs[reqIndex])
	}
}

// BenchmarkConcurrentRequests benchmarks concurrent request handling
func BenchmarkConcurrentRequests(b *testing.B) {
	app := New()

	app.Get("/test", func(c *Ctx) error {
		return c.SendString("ok")
	})

	req := httptest.NewRequest(http.MethodGet, "/test", nil)

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			w := httptest.NewRecorder()
			app.ServeHTTP(w, req)
		}
	})
}

// BenchmarkLargeResponse benchmarks large response performance
func BenchmarkLargeResponse(b *testing.B) {
	app := New()

	// Create a large JSON response
	largeData := make([]Map, 1000)
	for i := range largeData {
		largeData[i] = Map{
			"id":    i,
			"name":  "Item " + string(rune(i)),
			"value": i * 2,
			"tags":  []string{"tag1", "tag2", "tag3"},
		}
	}

	app.Get("/large", func(c *Ctx) error {
		return c.JSON(Map{"data": largeData})
	})

	req := httptest.NewRequest(http.MethodGet, "/large", nil)
	w := httptest.NewRecorder()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		app.ServeHTTP(w, req)
	}
}
