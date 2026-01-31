package ursa

import (
	"bytes"
	"io"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

// TestContextLocals tests context locals storage and retrieval
func TestContextLocals(t *testing.T) {
	app := New()

	app.Use(func(c *Ctx) error {
		c.Locals("user", "john")
		c.Locals("role", "admin")
		return c.Next()
	})

	app.Get("/test", func(c *Ctx) error {
		user := c.Locals("user")
		role := c.Locals("role")
		missing := c.Locals("missing")

		if user != "john" {
			t.Errorf("Expected user 'john', got '%v'", user)
		}
		if role != "admin" {
			t.Errorf("Expected role 'admin', got '%v'", role)
		}
		if missing != nil {
			t.Errorf("Expected missing to be nil, got '%v'", missing)
		}

		return c.JSON(Map{"user": user, "role": role})
	})

	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	w := httptest.NewRecorder()
	app.ServeHTTP(w, req)

	if w.Code != 200 {
		t.Errorf("Expected status 200, got %d", w.Code)
	}
	if !contains(w.Body.String(), "john") || !contains(w.Body.String(), "admin") {
		t.Errorf("Expected body to contain user and role, got '%s'", w.Body.String())
	}
}

// TestParameterExtraction tests various parameter extraction methods
func TestParameterExtraction(t *testing.T) {
	app := New()

	// Test path parameters
	app.Get("/users/:id/posts/:postId", func(c *Ctx) error {
		id := c.Param("id")
		postId := c.Param("postId")
		missing := c.Param("missing")

		if id != "123" {
			t.Errorf("Expected id '123', got '%s'", id)
		}
		if postId != "456" {
			t.Errorf("Expected postId '456', got '%s'", postId)
		}
		if missing != "" {
			t.Errorf("Expected missing to be empty, got '%s'", missing)
		}

		return c.JSON(Map{"userId": id, "postId": postId})
	})

	// Test query parameters
	app.Get("/search", func(c *Ctx) error {
		q := c.Query("q")
		limit := c.Query("limit")
		missing := c.Query("missing")

		if q != "golang" {
			t.Errorf("Expected query 'golang', got '%s'", q)
		}
		if limit != "10" {
			t.Errorf("Expected limit '10', got '%s'", limit)
		}
		if missing != "" {
			t.Errorf("Expected missing query to be empty, got '%s'", missing)
		}

		return c.JSON(Map{"query": q, "limit": limit})
	})

	// Test path parameters
	req := httptest.NewRequest(http.MethodGet, "/users/123/posts/456", nil)
	w := httptest.NewRecorder()
	app.ServeHTTP(w, req)
	if w.Code != 200 {
		t.Errorf("Expected status 200, got %d", w.Code)
	}

	// Test query parameters
	req = httptest.NewRequest(http.MethodGet, "/search?q=golang&limit=10", nil)
	w = httptest.NewRecorder()
	app.ServeHTTP(w, req)
	if w.Code != 200 {
		t.Errorf("Expected status 200, got %d", w.Code)
	}
}

// TestRequestParsing tests JSON and form body parsing
func TestRequestParsing(t *testing.T) {
	app := New()

	// Test JSON body parsing
	app.Post("/json", func(c *Ctx) error {
		type User struct {
			Name string `json:"name"`
			Age  int    `json:"age"`
		}

		var user User
		if err := c.BodyParser(&user); err != nil {
			t.Errorf("Failed to parse JSON: %v", err)
		}

		if user.Name != "John" || user.Age != 30 {
			t.Errorf("Expected user{Name:John, Age:30}, got %+v", user)
		}

		return c.JSON(Map{"parsed": user})
	})

	// Test form data parsing
	app.Post("/form", func(c *Ctx) error {
		name := c.FormValue("name")
		age := c.FormValue("age")

		if name != "Jane" || age != "25" {
			t.Errorf("Expected form values name=Jane, age=25, got name=%s, age=%s", name, age)
		}

		return c.JSON(Map{"name": name, "age": age})
	})

	// Test JSON request
	jsonBody := `{"name": "John", "age": 30}`
	req := httptest.NewRequest(http.MethodPost, "/json", strings.NewReader(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	app.ServeHTTP(w, req)
	if w.Code != 200 {
		t.Errorf("Expected status 200, got %d", w.Code)
	}

	// Test form request
	formBody := strings.NewReader("name=Jane&age=25")
	req = httptest.NewRequest(http.MethodPost, "/form", formBody)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	w = httptest.NewRecorder()
	app.ServeHTTP(w, req)
	if w.Code != 200 {
		t.Errorf("Expected status 200, got %d", w.Code)
	}
}

// TestResponseWriting tests various response writing methods
func TestResponseWriting(t *testing.T) {
	app := New()

	// Test SendString
	app.Get("/string", func(c *Ctx) error {
		return c.SendString("Hello, World!")
	})

	// Test SendStatus
	app.Get("/status", func(c *Ctx) error {
		return c.SendStatus(204)
	})

	// Test JSON response
	app.Get("/json", func(c *Ctx) error {
		return c.JSON(Map{"message": "Hello", "status": "ok"})
	})

	// Test string response
	req := httptest.NewRequest(http.MethodGet, "/string", nil)
	w := httptest.NewRecorder()
	app.ServeHTTP(w, req)
	if w.Code != 200 || w.Body.String() != "Hello, World!" {
		t.Errorf("String response failed: code=%d, body=%s", w.Code, w.Body.String())
	}

	// Test status response
	req = httptest.NewRequest(http.MethodGet, "/status", nil)
	w = httptest.NewRecorder()
	app.ServeHTTP(w, req)
	if w.Code != 204 {
		t.Errorf("Expected status 204, got %d", w.Code)
	}

	// Test JSON response
	req = httptest.NewRequest(http.MethodGet, "/json", nil)
	w = httptest.NewRecorder()
	app.ServeHTTP(w, req)
	if w.Code != 200 || !contains(w.Body.String(), "Hello") {
		t.Errorf("JSON response failed: code=%d, body=%s", w.Code, w.Body.String())
	}
}

// TestHeaders tests header manipulation
func TestHeaders(t *testing.T) {
	app := New()

	app.Get("/headers", func(c *Ctx) error {
		c.SetHeader("X-Custom", "custom-value")
		c.SetHeader("Content-Type", "text/plain")

		userAgent := c.Request.Header.Get("User-Agent")
		if userAgent == "" {
			t.Error("Expected User-Agent header")
		}

		return c.SendString("headers test")
	})

	req := httptest.NewRequest(http.MethodGet, "/headers", nil)
	req.Header.Set("User-Agent", "test-agent")
	w := httptest.NewRecorder()
	app.ServeHTTP(w, req)

	if w.Code != 200 {
		t.Errorf("Expected status 200, got %d", w.Code)
	}
	if w.Header().Get("X-Custom") != "custom-value" {
		t.Errorf("Expected X-Custom header to be 'custom-value', got '%s'", w.Header().Get("X-Custom"))
	}
}

// TestFileUpload tests file upload handling
func TestFileUpload(t *testing.T) {
	app := New()

	app.Post("/upload", func(c *Ctx) error {
		file, header, err := c.Request.FormFile("file")
		if err != nil {
			t.Errorf("Failed to get file: %v", err)
			return c.Status(400).SendString("No file")
		}
		defer file.Close()

		if header.Filename != "test.txt" {
			t.Errorf("Expected filename 'test.txt', got '%s'", header.Filename)
		}

		content, _ := io.ReadAll(file)
		if string(content) != "Hello, File!" {
			t.Errorf("Expected file content 'Hello, File!', got '%s'", string(content))
		}

		return c.JSON(Map{"filename": header.Filename, "size": header.Size})
	})

	// Create multipart form
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	part, _ := writer.CreateFormFile("file", "test.txt")
	part.Write([]byte("Hello, File!"))
	writer.Close()

	req := httptest.NewRequest(http.MethodPost, "/upload", body)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	w := httptest.NewRecorder()
	app.ServeHTTP(w, req)

	if w.Code != 200 {
		t.Errorf("Expected status 200, got %d", w.Code)
	}
	if !contains(w.Body.String(), "test.txt") {
		t.Errorf("Expected body to contain filename, got '%s'", w.Body.String())
	}
}

// TestQueryParser tests query parameter parsing into structs
func TestQueryParser(t *testing.T) {
	app := New()

	type SearchRequest struct {
		Query string   `query:"q"`
		Page  int      `query:"page"`
		Limit int      `query:"limit"`
		Tags  []string `query:"tags"`
		Sort  string   `query:"sort"`
	}

	app.Get("/search", func(c *Ctx) error {
		var req SearchRequest
		if err := c.QueryParser(&req); err != nil {
			t.Errorf("Failed to parse query: %v", err)
		}

		if req.Query != "golang" {
			t.Errorf("Expected query 'golang', got '%s'", req.Query)
		}
		if req.Page != 1 {
			t.Errorf("Expected page 1, got %d", req.Page)
		}
		if req.Limit != 10 {
			t.Errorf("Expected limit 10, got %d", req.Limit)
		}
		if len(req.Tags) != 2 || req.Tags[0] != "web" || req.Tags[1] != "api" {
			t.Errorf("Expected tags [web, api], got %v", req.Tags)
		}

		return c.JSON(Map{"parsed": req})
	})

	req := httptest.NewRequest(http.MethodGet, "/search?q=golang&page=1&limit=10&tags=web&tags=api&sort=desc", nil)
	w := httptest.NewRecorder()
	app.ServeHTTP(w, req)

	if w.Code != 200 {
		t.Errorf("Expected status 200, got %d", w.Code)
	}
}

// TestContextMethods tests various context helper methods
func TestContextMethods(t *testing.T) {
	app := New()

	app.Get("/methods", func(c *Ctx) error {
		// Test Method
		if c.Method() != "GET" {
			t.Errorf("Expected method GET, got %s", c.Method())
		}

		// Test Path
		if c.Path() != "/methods" {
			t.Errorf("Expected path /methods, got %s", c.Path())
		}

		// Test IP (basic check)
		ip := c.IP()
		if ip == "" {
			t.Error("Expected IP to be set")
		}

		// Test Status
		c.Status(202)
		if c.StatusCode != 202 {
			t.Errorf("Expected status 202, got %d", c.StatusCode)
		}

		return c.SendString("methods test")
	})

	req := httptest.NewRequest(http.MethodGet, "/methods", nil)
	w := httptest.NewRecorder()
	app.ServeHTTP(w, req)

	if w.Code != 202 {
		t.Errorf("Expected status 202, got %d", w.Code)
	}
}

// TestMultipleBodyParsing tests parsing body multiple times
func TestMultipleBodyParsing(t *testing.T) {
	app := New()

	app.Post("/multi", func(c *Ctx) error {
		// First parse as struct
		type User struct {
			Name string `json:"name"`
		}
		var user User
		if err := c.BodyParser(&user); err != nil {
			t.Errorf("Failed to parse body as struct: %v", err)
		}

		// Second parse as map
		var data map[string]interface{}
		if err := c.BodyParser(&data); err != nil {
			t.Errorf("Failed to parse body as map: %v", err)
		}

		if user.Name != "Alice" {
			t.Errorf("Expected name 'Alice', got '%s'", user.Name)
		}
		if data["name"] != "Alice" {
			t.Errorf("Expected map name 'Alice', got '%v'", data["name"])
		}

		return c.JSON(Map{"struct": user, "map": data})
	})

	jsonBody := `{"name": "Alice"}`
	req := httptest.NewRequest(http.MethodPost, "/multi", strings.NewReader(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	app.ServeHTTP(w, req)

	if w.Code != 200 {
		t.Errorf("Expected status 200, got %d", w.Code)
	}
}

// TestContextReset tests context reset functionality
func TestContextReset(t *testing.T) {
	app := New()

	var contextID string

	app.Use(func(c *Ctx) error {
		contextID = c.Request.Context().Value(TraceKey).(string)
		return c.Next()
	})

	app.Get("/test1", func(c *Ctx) error {
		currentID := c.Request.Context().Value(TraceKey).(string)
		if currentID != contextID {
			t.Errorf("Expected context ID to be consistent")
		}
		return c.SendString("test1")
	})

	app.Get("/test2", func(c *Ctx) error {
		currentID := c.Request.Context().Value(TraceKey).(string)
		if currentID != contextID {
			t.Errorf("Expected context ID to be consistent")
		}
		return c.SendString("test2")
	})

	// First request
	req1 := httptest.NewRequest(http.MethodGet, "/test1", nil)
	w1 := httptest.NewRecorder()
	app.ServeHTTP(w1, req1)

	// Second request
	req2 := httptest.NewRequest(http.MethodGet, "/test2", nil)
	w2 := httptest.NewRecorder()
	app.ServeHTTP(w2, req2)

	if w1.Code != 200 || w2.Code != 200 {
		t.Errorf("Expected both requests to succeed")
	}
}
