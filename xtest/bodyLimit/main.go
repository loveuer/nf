package main

import (
	"github.com/loveuer/nf"
	"log"
)

func main() {
	app := nf.New(nf.Config{BodyLimit: 30})

	app.Post("/data", func(c *nf.Ctx) error {
		type Req struct {
			Name  string   `json:"name"`
			Age   int      `json:"age"`
			Likes []string `json:"likes"`
		}

		var (
			err error
			req = new(Req)
		)

		if err = c.BodyParser(req); err != nil {
			return c.JSON(nf.Map{"status": 400, "err": err.Error()})
		}

		return c.JSON(nf.Map{"status": 200, "data": req})
	})

	app.Post("/url", func(c *nf.Ctx) error {
		type Req struct {
			Name  string   `form:"name"`
			Age   int      `form:"age"`
			Likes []string `form:"likes"`
		}

		var (
			err error
			req = new(Req)
		)

		if err = c.BodyParser(req); err != nil {
			return c.JSON(nf.Map{"status": 400, "err": err.Error()})
		}

		return c.JSON(nf.Map{"status": 200, "data": req})
	})

	log.Fatal(app.Run("0.0.0.0:80"))
}
