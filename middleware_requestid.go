package ursa

import (
	"github.com/google/uuid"
)

// RequestIDConfig defines the config for RequestID middleware
type RequestIDConfig struct {
	// Header is the header key where the request ID will be stored
	// Default: "X-Request-ID"
	Header string

	// Generator is a function that generates the request ID
	// Default: UUID v7 generator
	Generator func() string

	// ContextKey is the key used to store the request ID in the context
	// If empty, the request ID will not be stored in the context
	// Default: "RequestID"
	ContextKey string
}

// DefaultRequestIDConfig is the default RequestID middleware config
var DefaultRequestIDConfig = RequestIDConfig{
	Header: "X-Request-ID",
	Generator: func() string {
		return uuid.Must(uuid.NewV7()).String()
	},
	ContextKey: "RequestID",
}

// NewRequestID returns a RequestID middleware with default config
func NewRequestID() HandlerFunc {
	return NewRequestIDWithConfig(DefaultRequestIDConfig)
}

// NewRequestIDWithConfig returns a RequestID middleware with custom config
func NewRequestIDWithConfig(config RequestIDConfig) HandlerFunc {
	// Set defaults
	if config.Header == "" {
		config.Header = DefaultRequestIDConfig.Header
	}
	if config.Generator == nil {
		config.Generator = DefaultRequestIDConfig.Generator
	}
	if config.ContextKey == "" {
		config.ContextKey = DefaultRequestIDConfig.ContextKey
	}

	return func(c *Ctx) error {
		// Check if request ID already exists in the request header
		requestID := c.Get(config.Header)
		
		// If not, generate a new one
		if requestID == "" {
			requestID = config.Generator()
		}

		// Set the request ID in the response header
		c.Set(config.Header, requestID)

		// Store in context if ContextKey is provided
		if config.ContextKey != "" {
			c.Locals(config.ContextKey, requestID)
		}

		return c.Next()
	}
}
