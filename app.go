package nf

import (
	"context"
	"crypto/tls"
	"errors"
	"fmt"
	"net"
	"net/http"
	"strings"
)

type App struct {
	*RouterGroup
	config *Config
	router *router
	groups []*RouterGroup
	server *http.Server
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

func (a *App) run(ln net.Listener) error {
	if !a.config.DisableBanner {
		fmt.Println(banner + "nf serve at: " + a.server.Addr + "\n")
	}

	return a.server.Serve(ln)
}

func (a *App) Run(address string) error {
	ln, err := net.Listen("tcp", address)
	if err != nil {
		return err
	}

	a.server = &http.Server{}

	return a.run(ln)
}

func (a *App) RunTLS(address string, tlsConfig *tls.Config) error {
	ln, err := tls.Listen("tcp", address, tlsConfig)
	if err != nil {
		return err
	}

	a.server = &http.Server{}

	return a.run(ln)
}

func (a *App) RunListener(ln net.Listener) error {
	a.server = &http.Server{}

	return a.run(ln)
}

func (a *App) Shutdown(ctx context.Context) error {
	return a.server.Shutdown(ctx)
}
