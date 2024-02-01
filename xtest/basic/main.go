package main

import (
	"errors"
	"github.com/loveuer/nf"
	"github.com/loveuer/nf/nft/resp"
	"log"
	"time"
)

func main() {
	app := nf.New(nf.Config{EnableNotImplementHandler: true})

	app.Get("/hello/:name", func(c *nf.Ctx) error {
		name := c.Param("name")
		return c.JSON(nf.Map{"status": 200, "data": "hello, " + name})
	})
	app.Get("/not_impl")
	app.Patch("/world", func(c *nf.Ctx) error {
		time.Sleep(5 * time.Second)
		c.Status(404)
		return c.JSON(nf.Map{"method": c.Method, "status": c.StatusCode})
	})
	app.Get("/error", func(c *nf.Ctx) error {
		return resp.RespError(c, resp.NewError(404, "not found", errors.New("NNNot Found"), nil))
	})
	app.Post("/data", func(c *nf.Ctx) error {
		type Req struct {
			Name string `json:"name"`
		}

		var (
			err error
			req = new(Req)
			rm  = make(map[string]any)
		)

		if err = c.BodyParser(req); err != nil {
			return c.JSON(nf.Map{"status": 400, "msg": err.Error()})
		}

		if err = c.BodyParser(&rm); err != nil {
			return c.JSON(nf.Map{"status": 400, "msg": err.Error()})
		}

		return c.JSON(nf.Map{"status": 200, "data": req, "map": rm})
	})

	log.Fatal(app.Run(":80"))
}
