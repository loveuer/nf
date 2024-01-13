package main

import (
	"github.com/loveuer/nf"
	"log"
	"net"
)

func main() {
	app := nf.New()

	app.Get("/hello/:name", func(c *nf.Ctx) error {
		name := c.Param("name")
		return c.JSON(nf.Map{"status": 200, "data": "hello, " + name})
	})

	ln, _ := net.Listen("tcp", ":80")
	log.Fatal(app.RunListener(ln))
}
