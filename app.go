package nf

import (
	"errors"
	"fmt"
	"net/http"
	"strings"
)

type App struct {
	*RouterGroup
	config *Config
	router *router
	groups []*RouterGroup
}

func (a *App) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	c := newContext(a, writer, request)

	for _, group := range a.groups {
		if strings.HasPrefix(request.URL.Path, group.prefix) {
			c.handlers = append(c.handlers, group.middlewares...)
		}
	}

	if err := a.router.handle(c); err != nil {
		var ne = &Err{}

		if errors.As(err, ne) {
			writer.WriteHeader(ne.Status)
		} else {
			writer.WriteHeader(500)
		}

		_, _ = writer.Write([]byte(err.Error()))
	}
}

func (a *App) Run(address string) error {
	if !a.config.DisableBanner {
		fmt.Println(banner + "nf serve at: " + address + "\n")
	}
	return http.ListenAndServe(address, a)
}
