package main

import (
	"github.com/loveuer/nf"
	"log"
)

func main() {
	app := nf.New(nf.Config{DisableLogger: false})

	app.Get("/hello", func(c *nf.Ctx) error {
		return c.SendString("world")
	})

	app.Use(ml())

	log.Fatal(app.Run(":80"))
}

func ml() nf.HandlerFunc {
	return func(c *nf.Ctx) error {
		index := []byte(`<h1>my not found</h1>`)
		c.Set("Content-Type", "text/html")
		c.Status(403)
		_, err := c.Write(index)
		return err
	}
}
