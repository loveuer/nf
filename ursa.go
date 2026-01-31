package ursa

import (
	"sync"
	"time"
)

const (
	banner   = " _   _                  \n| | | |_ __ ___  __ _  \n| | | | '__/ __|/ _` | \n| |_| | |  \\__ \\ (_| | \n \\___/|_|  |___/\\__,_| \n "
	_404     = "<!doctype html><html lang=\"en\"><head><meta charset=\"UTF-8\"><meta name=\"viewport\" content=\"width=device-width,user-scalable=no,initial-scale=1,maximum-scale=1,minimum-scale=1\"><meta http-equiv=\"X-UA-Compatible\" content=\"ie=edge\"><title>Not Found</title><style>body{background:#333;margin:0;color:#ccc;display:flex;align-items:center;max-height:100vh;height:100vh;justify-content:center}textarea{min-height:5rem;min-width:20rem;text-align:center;border:none;background:0 0;color:#ccc;resize:none;user-input:none;user-select:none;cursor:default;-webkit-user-select:none;-webkit-touch-callout:none;-moz-user-select:none;-ms-user-select:none;outline:0}</style></head><body><textarea id=\"banner\" readonly=\"readonly\"></textarea><script type=\"text/javascript\">let htmlCodes = [\n    ' _   _                  ',\n    '| | | |_ __ ___  __ _  ',\n    '| | | | \\'__/ __|/ _\\` | ',\n    '| |_| | |  \\\\__ \\\\ (_| | ',\n    ' \\\\___/|_|  |___/\\\\__,_| '\n].join('\\n');\ndocument.querySelector('#banner').value = htmlCodes</script></body></html>"
	_405     = `405 Method Not Allowed`
	_500     = `500 Internal Server Error`
	TraceKey = "X-Trace-Id"
	version  = "ursa-1.0.0"
)

type Map map[string]interface{}

type Config struct {
	DisableMessagePrint bool `json:"-"`
	// Default: 4 * 1024 * 1024
	BodyLimit int64 `json:"-"`

	// Server timeout configurations
	ReadTimeout  time.Duration `json:"-"` // Default: 10s, maximum duration for reading the entire request
	WriteTimeout time.Duration `json:"-"` // Default: 10s, maximum duration before timing out writes of the response
	IdleTimeout  time.Duration `json:"-"` // Default: 120s, maximum amount of time to wait for the next request

	// if report http.ErrServerClosed as run err
	ErrServeClose bool `json:"-"`

	DisableBanner       bool `json:"-"`
	DisableLogger       bool `json:"-"`
	DisableRecover      bool `json:"-"`
	DisableHttpErrorLog bool `json:"-"`

	// EnableNotImplementHandler bool        `json:"-"`
	NotFoundHandler         HandlerFunc  `json:"-"`
	MethodNotAllowedHandler HandlerFunc  `json:"-"`
	BeforeServeFn           func(a *App) `json:"-"`
}

var defaultConfig = &Config{
	BodyLimit:    4 * 1024 * 1024,
	ReadTimeout:  10 * time.Second,
	WriteTimeout: 10 * time.Second,
	IdleTimeout:  120 * time.Second,
	NotFoundHandler: func(c *Ctx) error {
		c.Set("Content-Type", MIMETextHTML)
		_, err := c.Status(404).Write([]byte(_404))
		return err
	},
	MethodNotAllowedHandler: func(c *Ctx) error {
		c.Set("Content-Type", MIMETextPlain)
		_, err := c.Status(405).Write([]byte(_405))
		return err
	},
}

func New(config ...Config) *App {
	app := &App{
		RouterGroup: RouterGroup{
			Handlers: nil,
			basePath: "/",
			root:     true,
		},

		pool: &sync.Pool{},

		redirectTrailingSlash:  true,
		redirectFixedPath:      false,
		handleMethodNotAllowed: true,
		useRawPath:             false,
		unescapePathValues:     true,
		removeExtraSlash:       false,
	}

	app.config = defaultConfig

	if len(config) > 0 {
		cfg := config[0]

		if cfg.DisableMessagePrint {
			app.config.DisableMessagePrint = cfg.DisableMessagePrint
		}

		if cfg.DisableBanner {
			app.config.DisableBanner = cfg.DisableBanner
		}

		if cfg.DisableLogger {
			app.config.DisableLogger = cfg.DisableLogger
		}

		if cfg.DisableRecover {
			app.config.DisableRecover = cfg.DisableRecover
		}

		if cfg.DisableHttpErrorLog {
			app.config.DisableHttpErrorLog = cfg.DisableHttpErrorLog
		}

		if cfg.ErrServeClose {
			app.config.ErrServeClose = cfg.ErrServeClose
		}

		if cfg.BodyLimit > 0 {
			app.config.BodyLimit = cfg.BodyLimit
		}

		if cfg.ReadTimeout > 0 {
			app.config.ReadTimeout = cfg.ReadTimeout
		}

		if cfg.WriteTimeout > 0 {
			app.config.WriteTimeout = cfg.WriteTimeout
		}

		if cfg.IdleTimeout > 0 {
			app.config.IdleTimeout = cfg.IdleTimeout
		}

		if cfg.NotFoundHandler != nil {
			app.config.NotFoundHandler = cfg.NotFoundHandler
		}

		if cfg.MethodNotAllowedHandler != nil {
			app.config.MethodNotAllowedHandler = cfg.MethodNotAllowedHandler
		}

		if cfg.BeforeServeFn != nil {
			app.config.BeforeServeFn = cfg.BeforeServeFn
		}
	}

	app.RouterGroup.app = app

	app.Use(func(c *Ctx) error {
		c.SetHeader("server", version)
		return c.Next()
	})

	if !app.config.DisableLogger {
		app.Use(NewLogger())
	}

	if !app.config.DisableRecover {
		app.Use(NewRecover(true))
	}

	app.pool.New = func() any {
		return app.allocateContext()
	}

	return app
}
