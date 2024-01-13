package main

import (
	"context"
	"github.com/loveuer/nf"
	"log"
	"time"
)

func main() {
	app := nf.New()
	quit := make(chan bool)

	app.Get("/name", handleGet)

	go func() {
		err := app.Run(":7383")
		log.Print("run with err=", err)
	}()

	go func() {
		time.Sleep(5 * time.Second)
		err := app.Shutdown(context.TODO())
		log.Print("quit with err=", err)
		quit <- true
	}()

	<-quit

	log.Print("quited")
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

	return c.JSON(nf.Map{"req_map": req})
}
