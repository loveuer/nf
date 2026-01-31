package ursa

import "fmt"

type HandlerFunc func(*Ctx) error

func ToDoHandler(c *Ctx) error {
	return c.Status(501).SendString(fmt.Sprintf("%s - %s Not Implemented", c.Method(), c.Path()))
}
