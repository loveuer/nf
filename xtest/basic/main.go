package main

import (
	"github.com/loveuer/nf"
	"log"
)

func main() {
	app := nf.New(nf.Config{})

	app.Get("/ok", func(c *nf.Ctx) error {
		return c.SendStatus(200)
	})

	api := app.Group("/api")
	api.Get("/1", func(c *nf.Ctx) error {
		return c.SendString("nice")
	})

	log.Fatal(app.Run(":80"))
}
