package main

import (
	"github.com/loveuer/nf"
	"log"
)

func main() {
	app := nf.New()

	app.Get("/hello", func(c *nf.Ctx) error {
		type Req struct {
			Name  string   `query:"name"`
			Age   int      `query:"age"`
			Likes []string `query:"likes"`
		}

		var (
			err error
			req = new(Req)
		)

		if err = c.QueryParser(req); err != nil {
			return nf.NewNFError(400, err.Error())
		}

		return c.JSON(nf.Map{"status": 200, "data": req})
	})

	log.Fatal(app.Run("0.0.0.0:80"))
}
