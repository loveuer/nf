package nf

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net"
	"net/http"
	"strings"
)

type Ctx struct {
	writermem responseWriter
	// origin objects
	writer  http.ResponseWriter
	Request *http.Request
	// request info
	path   string
	method string
	// response info
	StatusCode int

	app          *App
	params       *Params
	index        int
	handlers     []HandlerFunc
	locals       map[string]interface{}
	skippedNodes *[]skippedNode
	fullPath     string
	Params       Params
}

func newContext(app *App, writer http.ResponseWriter, request *http.Request) *Ctx {

	skippedNodes := make([]skippedNode, 0, app.maxSections)
	v := make(Params, 0, app.maxParams)

	ctx := &Ctx{
		writer:     writer,
		writermem:  responseWriter{},
		Request:    request,
		path:       request.URL.Path,
		method:     request.Method,
		StatusCode: 200,

		app:          app,
		index:        -1,
		locals:       map[string]interface{}{},
		handlers:     make([]HandlerFunc, 0),
		skippedNodes: &skippedNodes,
		params:       &v,
	}

	ctx.writermem = responseWriter{
		ResponseWriter: ctx.writer,
		size:           -1,
		status:         0,
	}

	return ctx
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
	var (
		dv = ""
	)

	if len(defaultValue) > 0 {
		dv = defaultValue[0]
	}

	cookie, err := c.Request.Cookie(key)
	if err != nil || cookie.Value == "" {
		return dv
	}

	return cookie.Value
}

func (c *Ctx) Next() error {
	c.index++

	var err error

	for c.index < len(c.handlers) {
		if c.handlers[c.index] != nil {
			if err = c.handlers[c.index](c); err != nil {
				return err
			}
		}

		c.index++
	}

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

func (c *Ctx) Param(key string) string {
	return c.Params.ByName(key)
}

func (c *Ctx) Form(key string) string {
	return c.Request.FormValue(key)
}

func (c *Ctx) FormFile(key string) (*multipart.FileHeader, error) {
	_, fh, err := c.Request.FormFile(key)
	return fh, err
}

func (c *Ctx) Query(key string) string {
	return c.Request.URL.Query().Get(key)
}

func (c *Ctx) Get(key string, defaultValue ...string) string {
	value := c.Request.Header.Get(key)
	if value == "" && len(defaultValue) > 0 {
		return defaultValue[0]
	}

	return value
}

func (c *Ctx) IP() string {
	ip, _, err := net.SplitHostPort(strings.TrimSpace(c.Request.RemoteAddr))
	if err != nil {
		return ""
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
		bs, err := io.ReadAll(c.Request.Body)
		if err != nil {
			return err
		}
		_ = c.Request.Body.Close()

		c.Request.Body = io.NopCloser(bytes.NewReader(bs))

		return json.Unmarshal(bs, out)
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
	//v := reflect.ValueOf(out)
	//
	//if v.Kind() == reflect.Ptr && v.Elem().Kind() != reflect.Map {
	//}

	return parseToStruct("query", out, c.Request.URL.Query())
}

/* ===============================================================
|| Handle Ctx Response Part
=============================================================== */

func (c *Ctx) Status(code int) *Ctx {
	c.writermem.WriteHeader(code)
	c.StatusCode = c.writermem.status
	return c
}

func (c *Ctx) Set(key string, value string) {
	c.writermem.Header().Set(key, value)
}

func (c *Ctx) SetHeader(key string, value string) {
	c.writermem.Header().Set(key, value)
}

func (c *Ctx) SendString(data string) error {
	c.SetHeader("Content-Type", "text/plain")
	_, err := c.Write([]byte(data))
	return err
}

func (c *Ctx) Writef(format string, values ...interface{}) (int, error) {
	c.SetHeader("Content-Type", "text/plain")
	return c.writer.Write([]byte(fmt.Sprintf(format, values...)))
}

func (c *Ctx) JSON(data interface{}) error {
	c.SetHeader("Content-Type", MIMEApplicationJSON)

	encoder := json.NewEncoder(&c.writermem)

	if err := encoder.Encode(data); err != nil {
		return err
	}

	return nil
}

func (c *Ctx) RawWriter() http.ResponseWriter {
	return c.writer
}

func (c *Ctx) Write(data []byte) (int, error) {
	return c.writermem.Write(data)
}

func (c *Ctx) HTML(html string) error {
	c.SetHeader("Content-Type", "text/html")
	_, err := c.writer.Write([]byte(html))
	return err
}
