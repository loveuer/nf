package ursa

// SecureConfig defines the config for Secure middleware
type SecureConfig struct {
	// XSSProtection provides protection against cross-site scripting attack (XSS)
	// by setting the `X-XSS-Protection` header.
	// Default: "1; mode=block"
	XSSProtection string

	// ContentTypeNosniff provides protection against overriding Content-Type
	// header by setting the `X-Content-Type-Options` header.
	// Default: "nosniff"
	ContentTypeNosniff string

	// XFrameOptions can be used to indicate whether or not a browser should
	// be allowed to render a page in a <frame>, <iframe> or <object>.
	// Default: "SAMEORIGIN"
	XFrameOptions string

	// HSTSMaxAge sets the `Strict-Transport-Security` header to indicate how
	// long (in seconds) browsers should remember that this site is only to
	// be accessed using HTTPS.
	// Default: 0 (disabled)
	HSTSMaxAge int

	// HSTSIncludeSubdomains adds the includeSubDomains directive to the
	// `Strict-Transport-Security` header.
	// Default: false
	HSTSIncludeSubdomains bool

	// ContentSecurityPolicy sets the `Content-Security-Policy` header providing
	// security against cross-site scripting (XSS), clickjacking and other code
	// injection attacks.
	// Default: ""
	ContentSecurityPolicy string

	// ReferrerPolicy sets the `Referrer-Policy` header providing security against
	// leaking referrer information.
	// Default: ""
	ReferrerPolicy string

	// PermissionsPolicy sets the `Permissions-Policy` header providing control
	// over which features and APIs can be used in the browser.
	// Default: ""
	PermissionsPolicy string
}

// DefaultSecureConfig is the default Secure middleware config
var DefaultSecureConfig = SecureConfig{
	XSSProtection:      "1; mode=block",
	ContentTypeNosniff: "nosniff",
	XFrameOptions:      "SAMEORIGIN",
	HSTSMaxAge:         0,
	ReferrerPolicy:     "",
}

// NewSecure returns a Secure middleware with default config
func NewSecure() HandlerFunc {
	return NewSecureWithConfig(DefaultSecureConfig)
}

// NewSecureWithConfig returns a Secure middleware with custom config
func NewSecureWithConfig(config SecureConfig) HandlerFunc {
	// Set defaults
	if config.XSSProtection == "" {
		config.XSSProtection = DefaultSecureConfig.XSSProtection
	}
	if config.ContentTypeNosniff == "" {
		config.ContentTypeNosniff = DefaultSecureConfig.ContentTypeNosniff
	}
	if config.XFrameOptions == "" {
		config.XFrameOptions = DefaultSecureConfig.XFrameOptions
	}

	return func(c *Ctx) error {
		// X-XSS-Protection
		if config.XSSProtection != "" {
			c.Set("X-XSS-Protection", config.XSSProtection)
		}

		// X-Content-Type-Options
		if config.ContentTypeNosniff != "" {
			c.Set("X-Content-Type-Options", config.ContentTypeNosniff)
		}

		// X-Frame-Options
		if config.XFrameOptions != "" {
			c.Set("X-Frame-Options", config.XFrameOptions)
		}

		// Strict-Transport-Security
		if config.HSTSMaxAge > 0 {
			hsts := "max-age=" + string(rune(config.HSTSMaxAge))
			if config.HSTSIncludeSubdomains {
				hsts += "; includeSubDomains"
			}
			c.Set("Strict-Transport-Security", hsts)
		}

		// Content-Security-Policy
		if config.ContentSecurityPolicy != "" {
			c.Set("Content-Security-Policy", config.ContentSecurityPolicy)
		}

		// Referrer-Policy
		if config.ReferrerPolicy != "" {
			c.Set("Referrer-Policy", config.ReferrerPolicy)
		}

		// Permissions-Policy
		if config.PermissionsPolicy != "" {
			c.Set("Permissions-Policy", config.PermissionsPolicy)
		}

		return c.Next()
	}
}
