package proxy

import (
	"fmt"
	"github.com/loveuer/nf"
	"net/http/httputil"
	"net/url"
)

// NewProxy takes target host and creates a reverse proxy
func NewProxy(targetHost string) nf.HandlerFunc {

	return func(c *nf.Ctx) error {
		parse, err := url.Parse(targetHost)
		if err != nil {
			return err
		}

		proxy := httputil.NewSingleHostReverseProxy(parse)
		proxy.ServeHTTP(c.Writer, c.Request)
		fmt.Println(parse)

		return nil
	}
}
