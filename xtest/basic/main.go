package main

import (
	"errors"
	"github.com/loveuer/nf"
	"github.com/loveuer/nf/nft/resp"
	"log"
	"net"
	"time"
)

func main() {
	app := nf.New(nf.Config{EnableNotImplementHandler: true})

	app.Get("/hello/:name", func(c *nf.Ctx) error {
		name := c.Param("name")
		return c.JSON(nf.Map{"status": 200, "data": "hello, " + name})
	})
	app.Get("/not_impl")
	app.Patch("/world", func(c *nf.Ctx) error {
		time.Sleep(5 * time.Second)
		c.Status(404)
		return c.JSON(nf.Map{"method": c.Method, "status": c.StatusCode})
	})
	app.Get("/error", func(c *nf.Ctx) error {
		return resp.RespError(c, resp.NewError(404, "not found", errors.New("NNNot Found"), nil))
	})

	ln, _ := net.Listen("tcp", ":80")
	log.Fatal(app.RunListener(ln))
}
