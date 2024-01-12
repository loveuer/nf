package main

import (
	"github.com/loveuer/nf"
	"log"
)

func main() {
	app := nf.New()

	app.Get("/hello/:name", func(c *nf.Ctx) error {
		name := c.Param("name")
		return c.JSON(nf.Map{"status": 200, "data": "hello, " + name})
	})

	log.Fatal(app.Run("0.0.0.0:80"))
}
