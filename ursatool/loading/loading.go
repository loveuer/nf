package loading

import (
	"context"
	"fmt"
	"time"
)

type Type int

const (
	TypeProcessing Type = iota
	TypeInfo
	TypeSuccess
	TypeWarning
	TypeError
)

func (t Type) Symbol() string {
	switch t {
	case TypeSuccess:
		return "✔️  "
	case TypeWarning:
		return "❗ "
	case TypeError:
		return "❌ "
	case TypeInfo:
		return "❕ "
	default:
		return ""
	}
}

type _msg struct {
	msg string
	t   Type
}

var frames = []string{"|", "/", "-", "\\"}

func Do(ctx context.Context, fn func(ctx context.Context, print func(msg string, types ...Type)) error) (err error) {
	start := time.Now()
	ch := make(chan *_msg)

	defer func() {
		fmt.Printf("\r\033[K")
	}()

	go func() {
		var (
			m          *_msg
			ok         bool
			processing string
		)

		for {
			for _, frame := range frames {
				select {
				case <-ctx.Done():
					return
				case m, ok = <-ch:
					if !ok || m == nil {
						return
					}

					switch m.t {
					case TypeProcessing:
						if m.msg != "" {
							processing = m.msg
						}
					case TypeInfo,
						TypeSuccess,
						TypeWarning,
						TypeError:
						// Clear the loading animation
						fmt.Printf("\r\033[K")
						fmt.Printf("%s%s\n", m.t.Symbol(), m.msg)
					}
				default:
					elapsed := time.Since(start).Seconds()
					if processing != "" {
						fmt.Printf("\r\033[K%s  %s (%.2fs)", frame, processing, elapsed)
					}
					time.Sleep(100 * time.Millisecond)
				}
			}
		}
	}()

	printFn := func(msg string, types ...Type) {
		if msg == "" {
			return
		}

		m := &_msg{
			msg: msg,
			t:   TypeProcessing,
		}

		if len(types) > 0 {
			m.t = types[0]
		}

		ch <- m
	}

	done := make(chan struct{})
	go func() {
		if err = fn(ctx, printFn); err != nil {
			ch <- &_msg{msg: err.Error(), t: TypeError}
		}

		close(ch)
		done <- struct{}{}
	}()

	select {
	case <-ctx.Done():
	case <-done:
	}

	return err
}
