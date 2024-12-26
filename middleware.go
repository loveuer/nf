package nf

import (
	"fmt"
	"os"
	"runtime/debug"
	"time"

	"github.com/loveuer/nf/nft/log"
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

				// serveError(c, 500, []byte(fmt.Sprint(r)))
				_ = c.Status(500).SendString(fmt.Sprint(r))
			}
		}()

		return c.Next()
	}
}

func NewLogger() HandlerFunc {
	return func(c *Ctx) error {
		var (
			now   = time.Now()
			logFn func(msg string, data ...any)
			ip    = c.IP()
		)

		err := c.Next()
		duration := time.Since(now)

		msg := fmt.Sprintf("NF | %v | %15s | %3d | %s | %6s | %s", c.Context().Value(TraceKey), ip, c.StatusCode, HumanDuration(duration.Nanoseconds()), c.Method(), c.Path())

		switch {
		case c.StatusCode >= 500:
			logFn = log.Error
		case c.StatusCode >= 400:
			logFn = log.Warn
		default:
			logFn = log.Info
		}

		logFn(msg)

		return err
	}
}
