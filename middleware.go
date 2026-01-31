package ursa

import (
	"fmt"
	"os"
	"runtime/debug"
	"time"

	"github.com/loveuer/ursa/ursatool/log"
)

func NewRecover(enableStackTrace bool) HandlerFunc {
	return func(c *Ctx) error {
		defer func() {
			if r := recover(); r != nil {
				// Log detailed error internally only
				if enableStackTrace {
					os.Stderr.WriteString(fmt.Sprintf("recovered from panic: %v\nStack: %s", r, debug.Stack()))
				} else {
					os.Stderr.WriteString(fmt.Sprintf("recovered from panic: %v\n", r))
				}

				// Send generic error message to client to prevent information disclosure
				_ = c.Status(500).SendString("Internal Server Error")
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

		msg := fmt.Sprintf("URSA | %v | %15s | %3d | %s | %6s | %s", c.Context().Value(TraceKey), ip, c.StatusCode, HumanDuration(duration.Nanoseconds()), c.Method(), c.Path())

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
