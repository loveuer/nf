package main

import (
	"github.com/loveuer/nf"
	"log"
)

func main() {
	app := nf.New(nf.Config{EnableNotImplementHandler: true})

	api := app.Group("/api")
	api.Get("/1", func(c *nf.Ctx) error {
		return c.SendString("nice")
	})

	log.Fatal(app.Run(":80"))
}
