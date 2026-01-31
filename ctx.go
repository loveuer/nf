package ursa

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"html/template"
	"io"
	"mime/multipart"
	"net"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/google/uuid"
	"github.com/loveuer/ursa/internal/sse"
)

var forwardHeaders = []string{"CF-Connecting-IP", "X-Forwarded-For", "X-Real-Ip"}

type Ctx struct {
	lock       sync.Mutex
	writermem  responseWriter
	Writer     ResponseWriter
	Request    *http.Request
	path       string
	method     string
	StatusCode int

	app          *App
	params       *Params
	index        int
	handlers     []HandlerFunc
	locals       map[string]interface{}
	skippedNodes *[]skippedNode
	fullPath     string
}

func (c *Ctx) reset(w http.ResponseWriter, r *http.Request) {
	traceId := r.Header.Get(TraceKey)
	if traceId == "" {
		traceId = uuid.Must(uuid.NewV7()).String()
	}

	c.writermem.reset(w)

	c.Request = r.WithContext(context.WithValue(r.Context(), TraceKey, traceId))
	c.Writer = &c.writermem
	c.handlers = nil
	c.index = -1
	c.path = r.URL.Path
	c.method = r.Method
	c.StatusCode = 200

	c.fullPath = ""
	*c.params = (*c.params)[:0]
	*c.skippedNodes = (*c.skippedNodes)[:0]
	for key := range c.locals {
		delete(c.locals, key)
	}
	c.writermem.Header().Set(TraceKey, traceId)
}

func (c *Ctx) Drop() error {
	if h, ok := c.Writer.(http.Hijacker); ok {
		conn, _, _ := h.Hijack()
		return conn.Close()
	}

	return errors.New("hijack failed")
}

func (c *Ctx) Locals(key string, value ...interface{}) interface{} {
	data := c.locals[key]
	if len(value) > 0 {
		c.locals[key] = value[0]
	}

	return data
}

func (c *Ctx) Method(overWrite ...string) string {
	method := c.Request.Method

	if len(overWrite) > 0 && overWrite[0] != "" {
		c.Request.Method = overWrite[0]
	}

	return method
}

func (c *Ctx) Path(overWrite ...string) string {
	path := c.Request.URL.Path
	if len(overWrite) > 0 && overWrite[0] != "" {
		c.Request.URL.Path = overWrite[0]
	}

	return path
}

func (c *Ctx) Cookies(key string, defaultValue ...string) string {
	dv := ""

	if len(defaultValue) > 0 {
		dv = defaultValue[0]
	}

	cookie, err := c.Request.Cookie(key)
	if err != nil || cookie.Value == "" {
		return dv
	}

	return cookie.Value
}

func (c *Ctx) Context() context.Context {
	return c.Request.Context()
}

func (c *Ctx) SetContext(ctx context.Context) {
	c.Request = c.Request.WithContext(ctx)
}

func (c *Ctx) Next() error {
	c.index++

	if c.index >= len(c.handlers) {
		return nil
	}

	var (
		err     error
		handler = c.handlers[c.index]
	)

	if handler != nil {
		if err = handler(c); err != nil {
			return err
		}
	}

	c.index++

	return nil
}

/* ===============================================================
|| Handle Ctx Request Part
=============================================================== */

func (c *Ctx) verify() error {
	// 验证 body size
	if c.app.config.BodyLimit != -1 && c.Request.ContentLength > c.app.config.BodyLimit {
		return NewNFError(413, "Content Too Large")
	}

	return nil
}

func (c *Ctx) Body() []byte {
	bs, _ := io.ReadAll(c.Request.Body)
	return bs
}

func (c *Ctx) Param(key string) string {
	return c.params.ByName(key)
}

func (c *Ctx) SetParam(key, value string) {
	c.lock.Lock()
	defer c.lock.Unlock()

	params := append(*c.params, Param{Key: key, Value: value})
	c.params = &params
}

func (c *Ctx) Form(key string) string {
	return c.Request.FormValue(key)
}

// FormValue fiber ctx function
func (c *Ctx) FormValue(key string, defaultValue ...string) string {
	return defaultString(c.Request.FormValue(key), defaultValue)
}

func (c *Ctx) FormFile(key string) (*multipart.FileHeader, error) {
	_, fh, err := c.Request.FormFile(key)
	return fh, err
}

func (c *Ctx) MultipartForm() (*multipart.Form, error) {
	if err := c.Request.ParseMultipartForm(c.app.config.BodyLimit); err != nil {
		return nil, err
	}

	return c.Request.MultipartForm, nil
}

func (c *Ctx) Query(key string, defaultValue ...string) string {
	return defaultString(c.Request.URL.Query().Get(key), defaultValue)
}

func (c *Ctx) Queries() map[string]string {
	qs := c.Request.URL.Query()
	m := make(map[string]string, len(qs))
	for k, vs := range qs {
		if len(vs) > 0 {
			m[k] = vs[len(vs)-1]
		}
	}
	return m
}

func (c *Ctx) Get(key string, defaultValue ...string) string {
	value := c.Request.Header.Get(key)
	if value == "" && len(defaultValue) > 0 {
		return defaultValue[0]
	}
	return value
}

func (c *Ctx) Scheme() string {
	if c.Request.TLS != nil {
		return "https"
	}
	if scheme := c.Get("X-Forwarded-Proto"); scheme != "" {
		return scheme
	}
	return "http"
}

func (c *Ctx) Protocol() string {
	return c.Request.Proto
}

func (c *Ctx) IP(useProxyHeader ...bool) string {
	ip, _, err := net.SplitHostPort(strings.TrimSpace(c.Request.RemoteAddr))
	if err != nil {
		return ""
	}

	if len(useProxyHeader) > 0 && useProxyHeader[0] {
		for _, h := range forwardHeaders {
			if rip := c.Request.Header.Get(h); rip != "" {
				for _, part := range strings.Split(rip, ",") {
					realIP := strings.TrimSpace(part)
					if check := net.ParseIP(realIP); check != nil {
						return realIP
					}
				}
			}
		}
	}

	return ip
}

func (c *Ctx) BodyParser(out interface{}) error {
	var (
		err   error
		ctype = strings.ToLower(c.Request.Header.Get("Content-Type"))
	)

	ctype = parseVendorSpecificContentType(ctype)

	ctypeEnd := strings.IndexByte(ctype, ';')
	if ctypeEnd != -1 {
		ctype = ctype[:ctypeEnd]
	}

	if strings.HasSuffix(ctype, "json") {
		// Check if body has already been read
		if c.Request.Body == nil || c.Request.Body == http.NoBody {
			return NewNFError(400, "Request body is empty")
		}

		// Read body once and cache it
		bs, err := io.ReadAll(c.Request.Body)
		if err != nil {
			return err
		}
		_ = c.Request.Body.Close()
		
		// Allow multiple parsing by restoring the body
		c.Request.Body = io.NopCloser(bytes.NewReader(bs))
		
		// Use json.Decoder for better performance with large bodies
		decoder := json.NewDecoder(bytes.NewReader(bs))
		if err := decoder.Decode(out); err != nil {
			return NewNFError(400, "Invalid JSON: "+err.Error())
		}
		
		return nil
	}

	if strings.HasPrefix(ctype, MIMEApplicationForm) {
		if err = c.Request.ParseForm(); err != nil {
			return NewNFError(400, err.Error())
		}
		return parseToStruct("form", out, c.Request.Form)
	}

	if strings.HasPrefix(ctype, MIMEMultipartForm) {
		if err = c.Request.ParseMultipartForm(c.app.config.BodyLimit); err != nil {
			return NewNFError(400, err.Error())
		}
		return parseToStruct("form", out, c.Request.PostForm)
	}

	return NewNFError(422, "Unprocessable Content")
}

func (c *Ctx) QueryParser(out interface{}) error {
	return parseToStruct("query", out, c.Request.URL.Query())
}

func (c *Ctx) SaveFile(fh *multipart.FileHeader, path string) (err error) {
	var (
		f  multipart.File
		ff *os.File
	)

	f, err = fh.Open()
	if err != nil {
		return err
	}
	defer f.Close()

	// Clean path to prevent path traversal
	path = filepath.Clean(path)

	if ff, err = os.Create(path); err != nil {
		return err
	}
	defer ff.Close()

	_, err = copyZeroAlloc(ff, f)
	return err
}

/* ===============================================================
|| Handle Ctx Response Part
=============================================================== */

func (c *Ctx) Status(code int) *Ctx {
	c.lock.Lock()
	defer c.lock.Unlock()

	c.Writer.WriteHeader(code)
	c.StatusCode = c.writermem.status

	return c
}

// Set set response header
func (c *Ctx) Set(key string, value string) {
	c.Writer.Header().Set(key, value)
}

// AddHeader add response header
func (c *Ctx) AddHeader(key string, value string) {
	c.Writer.Header().Add(key, value)
}

// SetHeader set response header
func (c *Ctx) SetHeader(key string, value string) {
	c.Writer.Header().Set(key, value)
}

func (c *Ctx) SendStatus(code int) error {
	c.Status(code)
	c.Writer.WriteHeaderNow()
	return nil
}

func (c *Ctx) SendString(data string) error {
	c.SetHeader("Content-Type", "text/plain")
	_, err := c.Write([]byte(data))
	return err
}

func (c *Ctx) Writef(format string, values ...interface{}) (int, error) {
	c.SetHeader("Content-Type", "text/plain")
	return c.Write([]byte(fmt.Sprintf(format, values...)))
}

func (c *Ctx) JSON(data interface{}) error {
	c.SetHeader("Content-Type", MIMEApplicationJSON)

	encoder := json.NewEncoder(c.Writer)

	if err := encoder.Encode(data); err != nil {
		return err
	}

	return nil
}

func (c *Ctx) SSEvent(event string, data interface{}) error {
	c.Set("Content-Type", "text/event-stream")
	c.Set("Cache-Control", "no-cache")
	c.Set("Transfer-Encoding", "chunked")

	return sse.Encode(c.Writer, sse.Event{Event: event, Data: data})
}

func (c *Ctx) Flush() error {
	if f, ok := c.Writer.(http.Flusher); ok {
		f.Flush()
		return nil
	}

	return errors.New("http.Flusher is not implemented")
}

func (c *Ctx) HTML(html string) error {
	c.SetHeader("Content-Type", "text/html")
	_, err := c.Write([]byte(html))
	return err
}

func (c *Ctx) RenderHTML(name, html string, obj any) error {
	c.SetHeader("Content-Type", "text/html")
	t, err := template.New(name).Parse(html)
	if err != nil {
		return err
	}

	return t.Execute(c.Writer, obj)
}

func (c *Ctx) Redirect(url string, code int) error {
	http.Redirect(c.Writer, c.Request, url, code)
	return nil
}

func (c *Ctx) SendStream(stream io.Reader, size ...int) (err error) {
	if len(size) > 0 && size[0] > 0 {
		_, err = io.CopyN(c.Writer, stream, int64(size[0]))
	}

	_, err = io.Copy(c.Writer, stream)

	return err
}

func (c *Ctx) SendStreamWriter(streamWriter func(*bufio.Writer)) error {
	bw := bufio.NewWriter(c.Writer)
	streamWriter(bw)
	return nil
}

func (c *Ctx) Write(data []byte) (int, error) {
	return c.Writer.Write(data)
}
