package ursa

import (
	"net/http"
	"strconv"
	"strings"
)

// CORSConfig defines the config for CORS middleware
type CORSConfig struct {
	// AllowOrigins defines a list of origins that may access the resource.
	// Default: []string{"*"}
	AllowOrigins []string

	// AllowMethods defines a list of methods allowed when accessing the resource.
	// Default: []string{"GET", "POST", "HEAD", "PUT", "DELETE", "PATCH"}
	AllowMethods []string

	// AllowHeaders defines a list of request headers that can be used when making the actual request.
	// Default: []string{}
	AllowHeaders []string

	// AllowCredentials indicates whether or not the response to the request can be exposed when the credentials flag is true.
	// Default: false
	AllowCredentials bool

	// ExposeHeaders defines a whitelist headers that clients are allowed to access.
	// Default: []string{}
	ExposeHeaders []string

	// MaxAge indicates how long (in seconds) the results of a preflight request can be cached.
	// Default: 0 (no cache)
	MaxAge int
}

// DefaultCORSConfig is the default CORS middleware config
var DefaultCORSConfig = CORSConfig{
	AllowOrigins: []string{"*"},
	AllowMethods: []string{
		http.MethodGet,
		http.MethodPost,
		http.MethodHead,
		http.MethodPut,
		http.MethodDelete,
		http.MethodPatch,
	},
	AllowHeaders:     []string{},
	AllowCredentials: false,
	ExposeHeaders:    []string{},
	MaxAge:           0,
}

// NewCORS returns a CORS middleware with default config
func NewCORS() HandlerFunc {
	return NewCORSWithConfig(DefaultCORSConfig)
}

// NewCORSWithConfig returns a CORS middleware with custom config
func NewCORSWithConfig(config CORSConfig) HandlerFunc {
	// Set defaults
	if len(config.AllowOrigins) == 0 {
		config.AllowOrigins = DefaultCORSConfig.AllowOrigins
	}
	if len(config.AllowMethods) == 0 {
		config.AllowMethods = DefaultCORSConfig.AllowMethods
	}

	allowMethods := strings.Join(config.AllowMethods, ",")
	allowHeaders := strings.Join(config.AllowHeaders, ",")
	exposeHeaders := strings.Join(config.ExposeHeaders, ",")
	maxAge := strconv.Itoa(config.MaxAge)

	return func(c *Ctx) error {
		origin := c.Get("Origin")
		
		// Check if origin is allowed
		allowOrigin := ""
		if len(config.AllowOrigins) == 1 && config.AllowOrigins[0] == "*" {
			allowOrigin = "*"
		} else {
			for _, o := range config.AllowOrigins {
				if o == origin || o == "*" {
					allowOrigin = origin
					break
				}
			}
		}

		// If origin is not allowed, continue without CORS headers
		if allowOrigin == "" && origin != "" {
			allowOrigin = config.AllowOrigins[0]
		}

		// Set CORS headers
		if allowOrigin != "" {
			c.Set("Access-Control-Allow-Origin", allowOrigin)
		}

		if config.AllowCredentials {
			c.Set("Access-Control-Allow-Credentials", "true")
		}

		// Handle preflight requests
		if c.Request.Method == http.MethodOptions {
			c.Set("Access-Control-Allow-Methods", allowMethods)
			
			if allowHeaders != "" {
				c.Set("Access-Control-Allow-Headers", allowHeaders)
			} else {
				// If no allow headers are specified, echo the request headers
				if h := c.Get("Access-Control-Request-Headers"); h != "" {
					c.Set("Access-Control-Allow-Headers", h)
				}
			}

			if exposeHeaders != "" {
				c.Set("Access-Control-Expose-Headers", exposeHeaders)
			}

			if config.MaxAge > 0 {
				c.Set("Access-Control-Max-Age", maxAge)
			}

			return c.SendStatus(http.StatusNoContent)
		}

		// Set expose headers for actual requests
		if exposeHeaders != "" {
			c.Set("Access-Control-Expose-Headers", exposeHeaders)
		}

		return c.Next()
	}
}
