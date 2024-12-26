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

type Loading struct {
	Content string
	Type    Type
}

func Print(ctx context.Context, ch <-chan *Loading) {
	var (
		ok      bool
		frames  = []string{"|", "/", "-", "\\"}
		start   = time.Now()
		loading = &Loading{}
	)

	for {
		for _, frame := range frames {
			select {
			case <-ctx.Done():
				return
			case loading, ok = <-ch:
				if !ok || loading == nil {
					return
				}

				if loading.Content == "" {
					time.Sleep(100 * time.Millisecond)
					continue
				}

				switch loading.Type {
				case TypeInfo,
					TypeSuccess,
					TypeWarning,
					TypeError:
					// Clear the loading animation
					fmt.Printf("\r\033[K")
					fmt.Printf("%s%s\n", loading.Type.Symbol(), loading.Content)
					loading.Content = ""
				}
			default:
				elapsed := time.Since(start).Seconds()
				if loading.Content != "" {
					fmt.Printf("\r\033[K%s  %s (%.2fs)", frame, loading.Content, elapsed)
				}
				time.Sleep(100 * time.Millisecond)
			}
		}
	}
}
