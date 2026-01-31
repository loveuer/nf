package ursa

import (
	"errors"
	"fmt"
	"github.com/loveuer/ursa/internal/schema"
	"io"
	"net"
	"os"
	"reflect"
	"runtime"
	"strings"
	"sync"
)

const (
	MIMETextXML         = "text/xml"
	MIMETextHTML        = "text/html"
	MIMETextPlain       = "text/plain"
	MIMETextJavaScript  = "text/javascript"
	MIMEApplicationXML  = "application/xml"
	MIMEApplicationJSON = "application/json"
	MIMEApplicationForm = "application/x-www-form-urlencoded"
	MIMEOctetStream     = "application/octet-stream"
	MIMEMultipartForm   = "multipart/form-data"

	MIMETextXMLCharsetUTF8         = "text/xml; charset=utf-8"
	MIMETextHTMLCharsetUTF8        = "text/html; charset=utf-8"
	MIMETextPlainCharsetUTF8       = "text/plain; charset=utf-8"
	MIMETextJavaScriptCharsetUTF8  = "text/javascript; charset=utf-8"
	MIMEApplicationXMLCharsetUTF8  = "application/xml; charset=utf-8"
	MIMEApplicationJSONCharsetUTF8 = "application/json; charset=utf-8"
	// Deprecated: use MIMETextJavaScriptCharsetUTF8 instead
	MIMEApplicationJavaScriptCharsetUTF8 = "application/javascript; charset=utf-8"
)

// parseVendorSpecificContentType check if content type is vendor specific and
// if it is parsable to any known types. If it's not vendor specific then returns
// the original content type.
func parseVendorSpecificContentType(cType string) string {
	plusIndex := strings.Index(cType, "+")

	if plusIndex == -1 {
		return cType
	}

	var parsableType string
	if semiColonIndex := strings.Index(cType, ";"); semiColonIndex == -1 {
		parsableType = cType[plusIndex+1:]
	} else if plusIndex < semiColonIndex {
		parsableType = cType[plusIndex+1 : semiColonIndex]
	} else {
		return cType[:semiColonIndex]
	}

	slashIndex := strings.Index(cType, "/")

	if slashIndex == -1 {
		return cType
	}

	return cType[0:slashIndex+1] + parsableType
}

func parseToStruct(aliasTag string, out interface{}, data map[string][]string) error {
	schemaDecoder := schema.NewDecoder()
	schemaDecoder.SetAliasTag(aliasTag)

	if err := schemaDecoder.Decode(out, data); err != nil {
		return fmt.Errorf("failed to decode: %w", err)
	}

	return nil
}

func elsePanic(guard bool, text string) {
	if !guard {
		panic(text)
	}
}

func cleanPath(p string) string {
	const stackBufSize = 128
	// Turn empty string into "/"
	if p == "" {
		return "/"
	}

	// Reasonably sized buffer on stack to avoid allocations in the common case.
	// If a larger buffer is required, it gets allocated dynamically.
	buf := make([]byte, 0, stackBufSize)

	n := len(p)

	// Invariants:
	//      reading from path; r is index of next byte to process.
	//      writing to buf; w is index of next byte to write.

	// path must start with '/'
	r := 1
	w := 1

	if p[0] != '/' {
		r = 0

		if n+1 > stackBufSize {
			buf = make([]byte, n+1)
		} else {
			buf = buf[:n+1]
		}
		buf[0] = '/'
	}

	trailing := n > 1 && p[n-1] == '/'

	// A bit more clunky without a 'lazybuf' like the path package, but the loop
	// gets completely inlined (bufApp calls).
	// loop has no expensive function calls (except 1x make)		// So in contrast to the path package this loop has no expensive function
	// calls (except make, if needed).

	for r < n {
		switch {
		case p[r] == '/':
			// empty path element, trailing slash is added after the end
			r++

		case p[r] == '.' && r+1 == n:
			trailing = true
			r++

		case p[r] == '.' && p[r+1] == '/':
			// . element
			r += 2

		case p[r] == '.' && p[r+1] == '.' && (r+2 == n || p[r+2] == '/'):
			// .. element: remove to last /
			r += 3

			if w > 1 {
				// can backtrack
				w--

				if len(buf) == 0 {
					for w > 1 && p[w] != '/' {
						w--
					}
				} else {
					for w > 1 && buf[w] != '/' {
						w--
					}
				}
			}

		default:
			// Real path element.
			// Add slash if needed
			if w > 1 {
				bufApp(&buf, p, w, '/')
				w++
			}

			// Copy element
			for r < n && p[r] != '/' {
				bufApp(&buf, p, w, p[r])
				w++
				r++
			}
		}
	}

	// Re-append trailing slash
	if trailing && w > 1 {
		bufApp(&buf, p, w, '/')
		w++
	}

	// If the original string was not modified (or only shortened at the end),
	// return the respective substring of the original string.
	// Otherwise return a new string from the buffer.
	if len(buf) == 0 {
		return p[:w]
	}
	return string(buf[:w])
}

// Internal helper to lazily create a buffer if necessary.
// Calls to this function get inlined.
func bufApp(buf *[]byte, s string, w int, c byte) {
	b := *buf
	if len(b) == 0 {
		// No modification of the original string so far.
		// If the next character is the same as in the original string, we do
		// not yet have to allocate a buffer.
		if s[w] == c {
			return
		}

		// Otherwise use either the stack buffer, if it is large enough, or
		// allocate a new buffer on the heap, and copy all previous characters.
		length := len(s)
		if length > cap(b) {
			*buf = make([]byte, length)
		} else {
			*buf = (*buf)[:length]
		}
		b = *buf

		copy(b, s[:w])
	}
	b[w] = c
}

func HumanDuration(nano int64) string {
	duration := float64(nano)
	unit := "ns"
	if duration >= 1000 {
		duration /= 1000
		unit = "us"
	}

	if duration >= 1000 {
		duration /= 1000
		unit = "ms"
	}

	if duration >= 1000 {
		duration /= 1000
		unit = " s"
	}

	return fmt.Sprintf("%6.2f%s", duration, unit)
}

func defaultString(value string, defaultValue []string) string {
	if len(value) == 0 && len(defaultValue) > 0 {
		return defaultValue[0]
	}
	return value
}

var copyBufPool = sync.Pool{
	New: func() any {
		return make([]byte, 4096)
	},
}

func copyZeroAlloc(w io.Writer, r io.Reader) (int64, error) {
	var readerIsFile, readerIsConn bool

	switch r := r.(type) {
	case *os.File:
		readerIsFile = true
	case *net.TCPConn:
		readerIsConn = true
	case io.WriterTo:
		return r.WriteTo(w)
	}

	switch w := w.(type) {
	case *os.File:
		if readerIsConn {
			return w.ReadFrom(r)
		}
	case *net.TCPConn:
		if readerIsFile {
			// net.WriteTo requires go1.22 or later
			// Benchmark tests show that on Windows, WriteTo performs
			// significantly better than ReadFrom. On Linux, however,
			// ReadFrom slightly outperforms WriteTo. When possible,
			// copyZeroAlloc aims to perform  better than or as well
			// as io.Copy, so we use WriteTo whenever possible for
			// optimal performance.
			if rt, ok := r.(io.WriterTo); ok {
				return rt.WriteTo(w)
			}
			return w.ReadFrom(r)
		}
	case io.ReaderFrom:
		return w.ReadFrom(r)
	}

	vbuf := copyBufPool.Get()
	buf := vbuf.([]byte)
	n, err := copyBuffer(w, r, buf)
	copyBufPool.Put(vbuf)
	return n, err
}

func copyBuffer(dst io.Writer, src io.Reader, buf []byte) (written int64, err error) {
	for {
		nr, er := src.Read(buf)
		if nr > 0 {
			nw, ew := dst.Write(buf[0:nr])
			if nw < 0 || nr < nw {
				nw = 0
				if ew == nil {
					ew = errors.New("invalid write result")
				}
			}
			written += int64(nw)
			if ew != nil {
				err = ew
				break
			}
			if nr != nw {
				err = io.ErrShortWrite
				break
			}
		}
		if er != nil {
			if er != io.EOF {
				err = er
			}
			break
		}
	}
	return written, err
}

// getFunctionName 获取函数或方法的完整名称
func getFunctionName(i interface{}) string {
	v := reflect.ValueOf(i)
	if v.Kind() != reflect.Func {
		return "<not a function>"
	}
	pc := v.Pointer()
	if pc == 0 {
		return "<nil function>"
	}
	fn := runtime.FuncForPC(pc)
	if fn == nil {
		return "<unknown function>"
	}
	return fn.Name()
}

func _map[T, D any](items []T, fn func(item T, index int) D) []D {
	values := make([]D, len(items))
	for idx, item := range items {
		values[idx] = fn(item, idx)
	}

	return values
}

func _last[T any](items []T) (v T) {
	if len(items) == 0 {
		return
	}

	return items[len(items)-1]
}
