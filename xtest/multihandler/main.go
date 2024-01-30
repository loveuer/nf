package main

import (
	"github.com/loveuer/nf"
	"log"
)

func main() {
	app := nf.New()

	app.Get("/nice", h1, h2)

	log.Fatal(app.Run(":3333"))
}

func h1(c *nf.Ctx) error {
	you := c.Query("to")
	if you == "you" {
		return c.JSON(nf.Map{"status": 201, "msg": "nice to meet you"})
	}

	return c.Next()
}

func h2(c *nf.Ctx) error {
	return c.JSON(nf.Map{"status": 200, "msg": "hello world"})
}
