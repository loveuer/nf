package nf

import (
	"encoding/json"
	"fmt"
)

func (c *Ctx) Status(code int) *Ctx {
	c.StatusCode = code
	c.Writer.WriteHeader(code)
	return c
}

func (c *Ctx) SetHeader(key string, value string) {
	c.Writer.Header().Set(key, value)
}

func (c *Ctx) SendString(data string) error {
	c.SetHeader("Content-Type", "text/plain")
	_, err := c.Write([]byte(data))
	return err
}

func (c *Ctx) Writef(format string, values ...interface{}) (int, error) {
	c.SetHeader("Content-Type", "text/plain")
	return c.Writer.Write([]byte(fmt.Sprintf(format, values...)))
}

func (c *Ctx) JSON(data interface{}) error {
	c.SetHeader("Content-Type", "application/json")

	encoder := json.NewEncoder(c.Writer)

	if err := encoder.Encode(data); err != nil {
		return err
	}

	return nil
}

func (c *Ctx) Write(data []byte) (int, error) {
	return c.Writer.Write(data)
}

func (c *Ctx) HTML(html string) error {
	c.SetHeader("Content-Type", "text/html")
	_, err := c.Writer.Write([]byte(html))
	return err
}
