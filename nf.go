package nf

const (
	banner = "  _  _     _     ___                 _ \n | \\| |___| |_  | __|__ _  _ _ _  __| |\n | .` / _ \\  _| | _/ _ \\ || | ' \\/ _` |\n |_|\\_\\___/\\__| |_|\\___/\\_,_|_||_\\__,_|\n "
)

type Map map[string]interface{}

type Config struct {
	// Default: 4 * 1024 * 1024
	BodyLimit int64 `json:"-"`

	// if report http.ErrServerClosed as run err
	ErrServeClose bool `json:"-"`

	DisableBanner  bool `json:"-"`
	DisableLogger  bool `json:"-"`
	DisableRecover bool `json:"-"`
}

var (
	defaultConfig = &Config{
		BodyLimit: 4 * 1024 * 1024,
	}
)

func New(config ...Config) *App {
	app := &App{
		router: newRouter(),
	}

	if len(config) > 0 {
		app.config = &config[0]
		if app.config.BodyLimit == 0 {
			app.config.BodyLimit = defaultConfig.BodyLimit
		}
	} else {
		app.config = defaultConfig
	}

	app.RouterGroup = &RouterGroup{app: app}
	app.groups = []*RouterGroup{app.RouterGroup}

	if !app.config.DisableRecover {
		app.Use(NewRecover(true))
	}

	return app
}
