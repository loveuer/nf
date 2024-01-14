package main

import (
	"github.com/loveuer/nf"
	"log"
	"net"
	"time"
)

func main() {
	app := nf.New()

	app.Get("/hello/:name", func(c *nf.Ctx) error {
		name := c.Param("name")
		return c.JSON(nf.Map{"status": 200, "data": "hello, " + name})
	})
	app.Patch("/world", func(c *nf.Ctx) error {
		time.Sleep(5 * time.Second)
		c.Status(404)
		return c.JSON(nf.Map{"method": c.Method, "status": c.StatusCode})
	})

	ln, _ := net.Listen("tcp", ":80")
	log.Fatal(app.RunListener(ln))
}
