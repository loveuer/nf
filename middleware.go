package nf

import (
	"fmt"
	"github.com/google/uuid"
	"github.com/loveuer/nf/nft/log"
	"os"
	"runtime/debug"
	"strings"
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

				//serveError(c, 500, []byte(fmt.Sprint(r)))
				_ = c.Status(500).SendString(fmt.Sprint(r))
			}
		}()

		return c.Next()
	}
}

func NewLogger(traceHeader ...string) HandlerFunc {
	Header := "X-Trace-ID"
	if len(traceHeader) > 0 && traceHeader[0] != "" {
		Header = traceHeader[0]
	}

	return func(c *Ctx) error {
		var (
			now   = time.Now()
			trace = c.Get(Header)
			logFn func(msg string, data ...any)
			ip    = c.IP()
		)

		if trace == "" {
			trace = uuid.Must(uuid.NewV7()).String()
		}

		c.SetHeader(Header, trace)

		traces := strings.Split(trace, "-")
		shortTrace := traces[len(traces)-1]

		err := c.Next()
		duration := time.Since(now)

		msg := fmt.Sprintf("NF | %s | %15s | %3d | %s | %6s | %s", shortTrace, ip, c.StatusCode, HumanDuration(duration.Nanoseconds()), c.Method(), c.Path())

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
