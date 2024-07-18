package proxy

import (
	"fmt"
	"github.com/loveuer/nf"
	"net/http"
	"testing"
)

func TestNewProxy(t *testing.T) {

	app := nf.New()
	t1 := app.Group("/test")
	t1.Get("/hello", helloHandler)

	t2 := app.Group("/proxy")
	t2.Use(NewProxy("http://127.0.0.1:8081"))
	t2.Get("/hello", helloHandler)

	fmt.Println("Starting server on :8080")
	go func() {
		if err := http.ListenAndServe(":8080", app); err != nil {
			fmt.Println("app start err:", err)
		}
	}()

	app2 := nf.New()
	app2.Get("/proxy/hello", helloHandler)
	fmt.Println("Starting server on :8081")
	if err := http.ListenAndServe(":8081", app2); err != nil {
		fmt.Println("app2 start err:", err)
	}
}

func helloHandler(c *nf.Ctx) error {
	fmt.Println(c.Request.URL)
	_, _ = fmt.Fprintf(c.Writer, "Hello, World!")
	return nil
}
