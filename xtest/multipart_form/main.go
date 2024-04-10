package main

import (
	"github.com/loveuer/nf"
	"github.com/loveuer/nf/nft/resp"
	"log"
)

func main() {
	app := nf.New(nf.Config{BodyLimit: 10 * 1024 * 1024})

	app.Post("/upload", func(c *nf.Ctx) error {
		fs, err := c.MultipartForm()
		if err != nil {
			return resp.Resp400(c, err.Error())
		}

		fm := make(map[string][]string)
		for key := range fs.File {
			if _, exist := fm[key]; !exist {
				fm[key] = make([]string, 0)
			}

			for f := range fs.File[key] {
				fm[key] = append(fm[key], fs.File[key][f].Filename)
			}
		}

		return resp.Resp200(c, nf.Map{"value": fs.Value, "files": fm})
	})

	log.Fatal(app.Run(":13322"))
}
