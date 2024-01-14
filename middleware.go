package nf

import (
	"fmt"
	"log"
	"os"
	"runtime/debug"
	"time"
)

func NewRecover(enableStackTrace bool) HandlerFunc {
	return func(c *Ctx) error {
		defer func() {
			if r := recover(); r != nil {
				if enableStackTrace {
					os.Stderr.WriteString(fmt.Sprintf("recovered from panic: %v\n%s\n", r, debug.Stack()))
				} else {
					os.Stderr.WriteString(fmt.Sprintf("recovered from panic: %v\n", r))
				}
			}
		}()

		return c.Next()
	}
}

func NewLogger() HandlerFunc {
	l := log.New(os.Stdout, "[NF] ", 0)

	durationFormat := func(num int64) string {
		var (
			unit = "ns"
		)

		if num > 1000 {
			num = num / 1000
			unit = "µs"
		}

		if num > 1000 {
			num = num / 1000
			unit = "ms"
		}

		if num > 1000 {
			num = num / 1000
			unit = " s"
		}

		return fmt.Sprintf("%v %s", num, unit)
	}

	return func(c *Ctx) error {
		start := time.Now()

		err := c.Next()

		var (
			duration = time.Now().Sub(start).Nanoseconds()
			status   = c.StatusCode
			path     = c.path
			method   = c.Request.Method
		)

		l.Printf("%s | %5s | %d | %s | %s",
			start.Format("06/01/02T15:04:05"),
			method,
			status,
			durationFormat(duration),
			path,
		)

		return err
	}
}
