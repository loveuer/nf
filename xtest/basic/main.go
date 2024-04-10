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
	api.Use(func(c *nf.Ctx) error {
		c.SetParam("age", "18")
		return c.Next()
	})

	api.Get("/:name", func(c *nf.Ctx) error {
		name := c.Param("name")
		age := c.Param("age")
		return c.SendString(name + "@" + age)
	})

	log.Fatal(app.Run(":80"))
}
