package main

import (
	"context"
	"github.com/loveuer/nf"
	"log"
	"time"
)

var (
	app  = nf.New()
	quit = make(chan bool)
)

func main() {

	app.Get("/name", handleGet)

	go func() {
		err := app.Run(":80")
		log.Print("run with err=", err)
		quit <- true
	}()

	<-quit
}

func handleGet(c *nf.Ctx) error {
	type Req struct {
		Name string   `query:"name"`
		Addr []string `query:"addr"`
	}

	var (
		err error
		req = Req{}
	)

	if err = c.QueryParser(&req); err != nil {
		return nf.NewNFError(400, err.Error())
	}

	if req.Name == "quit" {

		go func() {
			time.Sleep(2 * time.Second)
			log.Print("app quit = ", app.Shutdown(context.TODO()))
		}()
	}

	return c.JSON(nf.Map{"req_map": req})
}
