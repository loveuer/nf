package main

import (
	"github.com/loveuer/nf"
	"log"
)

func main() {
	app := nf.New()

	app.Use(ml())
	app.Get("/hello", func(c *nf.Ctx) error {
		return c.SendString("world")
	})

	log.Fatal(app.Run(":7777"))
}

func ml() nf.HandlerFunc {
	return func(c *nf.Ctx) error {
		log.Printf("[ML] [%s] - [%s]", c.Method, c.Path())
		return c.Next()
	}
}
