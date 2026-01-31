package ursa

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

// TestServerTimeouts tests server timeout configuration
func TestServerTimeouts(t *testing.T) {
	app := New(Config{
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 5 * time.Second,
		IdleTimeout:  60 * time.Second,
	})

	app.Get("/test", func(c *Ctx) error {
		return c.SendString("test")
	})

	// Verify timeouts are set in server config
	if app.config.ReadTimeout != 5*time.Second {
		t.Errorf("Expected ReadTimeout 5s, got %v", app.config.ReadTimeout)
	}
	if app.config.WriteTimeout != 5*time.Second {
		t.Errorf("Expected WriteTimeout 5s, got %v", app.config.WriteTimeout)
	}
	if app.config.IdleTimeout != 60*time.Second {
		t.Errorf("Expected IdleTimeout 60s, got %v", app.config.IdleTimeout)
	}
}

// TestCORSMiddleware tests CORS middleware
func TestCORSMiddleware(t *testing.T) {
	app := New()
	app.Use(NewCORS())

	app.Get("/test", func(c *Ctx) error {
		return c.SendString("test")
	})

	// Test preflight request
	req := httptest.NewRequest(http.MethodOptions, "/test", nil)
	req.Header.Set("Origin", "https://example.com")
	req.Header.Set("Access-Control-Request-Method", "GET")
	w := httptest.NewRecorder()
	app.ServeHTTP(w, req)

	if w.Code != http.StatusNoContent {
		t.Errorf("Expected status 204, got %d", w.Code)
	}

	if w.Header().Get("Access-Control-Allow-Origin") == "" {
		t.Error("Expected Access-Control-Allow-Origin header to be set")
	}

	// Test actual request
	req = httptest.NewRequest(http.MethodGet, "/test", nil)
	req.Header.Set("Origin", "https://example.com")
	w = httptest.NewRecorder()
	app.ServeHTTP(w, req)

	if w.Code != 200 {
		t.Errorf("Expected status 200, got %d", w.Code)
	}

	if w.Header().Get("Access-Control-Allow-Origin") == "" {
		t.Error("Expected Access-Control-Allow-Origin header to be set")
	}
}

// TestCORSMiddlewareWithConfig tests CORS middleware with custom config
func TestCORSMiddlewareWithConfig(t *testing.T) {
	app := New()
	app.Use(NewCORSWithConfig(CORSConfig{
		AllowOrigins:     []string{"https://example.com", "https://test.com"},
		AllowMethods:     []string{"GET", "POST"},
		AllowHeaders:     []string{"Content-Type"},
		AllowCredentials: true,
		MaxAge:           3600,
	}))

	app.Get("/test", func(c *Ctx) error {
		return c.SendString("test")
	})

	// Test with allowed origin
	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	req.Header.Set("Origin", "https://example.com")
	w := httptest.NewRecorder()
	app.ServeHTTP(w, req)

	if w.Header().Get("Access-Control-Allow-Origin") != "https://example.com" {
		t.Errorf("Expected origin https://example.com, got %s", w.Header().Get("Access-Control-Allow-Origin"))
	}

	if w.Header().Get("Access-Control-Allow-Credentials") != "true" {
		t.Error("Expected credentials to be allowed")
	}
}

// TestSecureMiddleware tests Secure middleware
func TestSecureMiddleware(t *testing.T) {
	app := New()
	app.Use(NewSecure())

	app.Get("/test", func(c *Ctx) error {
		return c.SendString("test")
	})

	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	w := httptest.NewRecorder()
	app.ServeHTTP(w, req)

	if w.Code != 200 {
		t.Errorf("Expected status 200, got %d", w.Code)
	}

	// Check security headers
	if w.Header().Get("X-XSS-Protection") != "1; mode=block" {
		t.Errorf("Expected X-XSS-Protection header, got %s", w.Header().Get("X-XSS-Protection"))
	}

	if w.Header().Get("X-Content-Type-Options") != "nosniff" {
		t.Errorf("Expected X-Content-Type-Options header, got %s", w.Header().Get("X-Content-Type-Options"))
	}

	if w.Header().Get("X-Frame-Options") != "SAMEORIGIN" {
		t.Errorf("Expected X-Frame-Options header, got %s", w.Header().Get("X-Frame-Options"))
	}
}

// TestSecureMiddlewareWithConfig tests Secure middleware with custom config
func TestSecureMiddlewareWithConfig(t *testing.T) {
	app := New()
	app.Use(NewSecureWithConfig(SecureConfig{
		XSSProtection:         "1; mode=block",
		ContentTypeNosniff:    "nosniff",
		XFrameOptions:         "DENY",
		HSTSMaxAge:            31536000,
		HSTSIncludeSubdomains: true,
		ContentSecurityPolicy: "default-src 'self'",
		ReferrerPolicy:        "no-referrer",
	}))

	app.Get("/test", func(c *Ctx) error {
		return c.SendString("test")
	})

	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	w := httptest.NewRecorder()
	app.ServeHTTP(w, req)

	if w.Header().Get("X-Frame-Options") != "DENY" {
		t.Errorf("Expected X-Frame-Options DENY, got %s", w.Header().Get("X-Frame-Options"))
	}

	if w.Header().Get("Content-Security-Policy") != "default-src 'self'" {
		t.Errorf("Expected CSP header, got %s", w.Header().Get("Content-Security-Policy"))
	}

	if w.Header().Get("Referrer-Policy") != "no-referrer" {
		t.Errorf("Expected Referrer-Policy header, got %s", w.Header().Get("Referrer-Policy"))
	}
}

// TestRequestIDMiddleware tests RequestID middleware
func TestRequestIDMiddleware(t *testing.T) {
	app := New()
	app.Use(NewRequestID())

	app.Get("/test", func(c *Ctx) error {
		requestID := c.Locals("RequestID")
		if requestID == nil {
			t.Error("Expected RequestID to be set in Locals")
		}
		return c.SendString("test")
	})

	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	w := httptest.NewRecorder()
	app.ServeHTTP(w, req)

	if w.Code != 200 {
		t.Errorf("Expected status 200, got %d", w.Code)
	}

	requestID := w.Header().Get("X-Request-ID")
	if requestID == "" {
		t.Error("Expected X-Request-ID header to be set")
	}
}

// TestRequestIDMiddlewareWithExisting tests RequestID middleware with existing ID
func TestRequestIDMiddlewareWithExisting(t *testing.T) {
	app := New()
	app.Use(NewRequestID())

	app.Get("/test", func(c *Ctx) error {
		return c.SendString("test")
	})

	existingID := "existing-request-id"
	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	req.Header.Set("X-Request-ID", existingID)
	w := httptest.NewRecorder()
	app.ServeHTTP(w, req)

	if w.Header().Get("X-Request-ID") != existingID {
		t.Errorf("Expected RequestID %s, got %s", existingID, w.Header().Get("X-Request-ID"))
	}
}
