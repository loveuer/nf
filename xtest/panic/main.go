package main

import (
	"github.com/loveuer/nf"
	"log"
)

func main() {
	app := nf.New(nf.Config{
		DisableRecover: false,
	})

	app.Get("/hello/:name", func(c *nf.Ctx) error {
		name := c.Param("name")

		if name == "nf" {
			panic("name is nf")
		}

		return c.JSON("nice")
	})

	log.Fatal(app.Run("0.0.0.0:80"))
}
