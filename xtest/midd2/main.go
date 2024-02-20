package main

import (
	"errors"
	"github.com/loveuer/nf"
	"github.com/loveuer/nf/nft/resp"
	"log"
)

func main() {
	app := nf.New()

	api := app.Group("/api")

	api.Get("/hello",
		auth(),
		func(c *nf.Ctx) error {
			return resp.Resp403(c, errors.New("in hello"))
		},
	)

	log.Fatal(app.Run(":80"))
}

func auth() nf.HandlerFunc {
	return func(c *nf.Ctx) error {
		token := c.Query("token")
		if token != "zyp" {
			return resp.Resp401(c, errors.New("no auth"))
		}

		return c.Next()
	}
}
