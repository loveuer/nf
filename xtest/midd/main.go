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
	app.Get("/panic", func(c *nf.Ctx) error {
		panic("panic")
	})
	app.Use(ml())

	log.Fatal(app.Run(":7777"))
}

func ml() nf.HandlerFunc {
	return func(c *nf.Ctx) error {
		log.Printf("[ML] [%s] - [%s]", c.Method, c.Path())
		index := []byte(`<h1>my not found</h1>`)
		c.Set("Content-Type", "text/html")
		_, err := c.Write(index)
		return err
	}
}
