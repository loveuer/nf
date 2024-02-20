package nf

import (
	"context"
	"crypto/tls"
	"errors"
	"fmt"
	"github.com/loveuer/nf/internal/bytesconv"
	"io"
	"log"
	"net"
	"net/http"
	"path"
	"regexp"
)

var (
	_ IRouter = (*App)(nil)

	regSafePrefix         = regexp.MustCompile("[^a-zA-Z0-9/-]+")
	regRemoveRepeatedChar = regexp.MustCompile("/{2,}")

	mimePlain = []string{"text/plain"}
)

type App struct {
	RouterGroup
	config *Config
	groups []*RouterGroup
	server *http.Server

	trees methodTrees

	maxParams   uint16
	maxSections uint16

	redirectTrailingSlash  bool // true
	redirectFixedPath      bool // false
	handleMethodNotAllowed bool // false
	useRawPath             bool // false
	unescapePathValues     bool // true
	removeExtraSlash       bool // false
}

func (a *App) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	var (
		err error
		c   = newContext(a, writer, request)
		nfe = new(Err)
	)

	if err = c.verify(); err != nil {
		if errors.As(err, nfe) {
			_ = c.Status(nfe.Status).SendString(nfe.Msg)
			return
		}

		_ = c.Status(500).SendString(err.Error())
		return
	}

	a.handleHTTPRequest(c)
}

func (a *App) run(ln net.Listener) error {
	srv := &http.Server{Handler: a}

	if a.config.DisableHttpErrorLog {
		srv.ErrorLog = log.New(io.Discard, "", 0)
	}

	a.server = srv

	if !a.config.DisableBanner {
		fmt.Println(banner + "nf serve at: " + ln.Addr().String() + "\n")
	}

	err := a.server.Serve(ln)
	if !errors.Is(err, http.ErrServerClosed) || a.config.ErrServeClose {
		return err
	}

	return nil
}

func (a *App) Run(address string) error {
	ln, err := net.Listen("tcp", address)
	if err != nil {
		return err
	}

	return a.run(ln)
}

func (a *App) RunTLS(address string, tlsConfig *tls.Config) error {
	ln, err := tls.Listen("tcp", address, tlsConfig)
	if err != nil {
		return err
	}

	return a.run(ln)
}

func (a *App) RunListener(ln net.Listener) error {
	a.server = &http.Server{Addr: ln.Addr().String()}

	return a.run(ln)
}

func (a *App) Shutdown(ctx context.Context) error {
	return a.server.Shutdown(ctx)
}

func (a *App) addRoute(method, path string, handlers ...HandlerFunc) {
	elsePanic(path[0] == '/', "path must begin with '/'")
	elsePanic(method != "", "HTTP method can not be empty")
	elsePanic(len(handlers) > 0, "without enable not implement, there must be at least one handler")

	if !a.config.DisableMessagePrint {
		fmt.Printf("[NF] Add Route: %-8s - %-25s (%2d handlers)\n", method, path, len(handlers))
	}

	root := a.trees.get(method)
	if root == nil {
		root = new(node)
		root.fullPath = "/"
		a.trees = append(a.trees, methodTree{method: method, root: root})
	}

	root.addRoute(path, handlers...)

	if paramsCount := countParams(path); paramsCount > a.maxParams {
		a.maxParams = paramsCount
	}

	if sectionsCount := countSections(path); sectionsCount > a.maxSections {
		a.maxSections = sectionsCount
	}
}

func (a *App) handleHTTPRequest(c *Ctx) {
	var (
		err error
	)

	httpMethod := c.Request.Method
	rPath := c.Request.URL.Path
	unescape := false
	if a.useRawPath && len(c.Request.URL.RawPath) > 0 {
		rPath = c.Request.URL.RawPath
		unescape = a.unescapePathValues
	}

	if a.removeExtraSlash {
		rPath = cleanPath(rPath)
	}

	// Find root of the tree for the given HTTP method
	t := a.trees
	for i, tl := 0, len(t); i < tl; i++ {
		if t[i].method != httpMethod {
			continue
		}
		root := t[i].root
		// Find route in tree
		value := root.getValue(rPath, c.params, c.skippedNodes, unescape)
		if value.params != nil {
			c.Params = *value.params
		}
		if value.handlers != nil {
			c.handlers = value.handlers
			c.fullPath = value.fullPath

			if err = c.Next(); err != nil {
				serveError(c, errorHandler)
			}

			c.writermem.WriteHeaderNow()
			return
		}
		if httpMethod != http.MethodConnect && rPath != "/" {
			if value.tsr && a.redirectTrailingSlash {
				redirectTrailingSlash(c)
				return
			}
			if a.redirectFixedPath && redirectFixedPath(c, root, a.redirectFixedPath) {
				return
			}
		}
		break
	}

	if a.handleMethodNotAllowed {
		// According to RFC 7231 section 6.5.5, MUST generate an Allow header field in response
		// containing a list of the target resource's currently supported methods.
		allowed := make([]string, 0, len(t)-1)
		for _, tree := range a.trees {
			if tree.method == httpMethod {
				continue
			}
			if value := tree.root.getValue(rPath, nil, c.skippedNodes, unescape); value.handlers != nil {
				allowed = append(allowed, tree.method)
			}
		}

		if len(allowed) > 0 {
			c.handlers = a.combineHandlers()

			serveError(c, a.config.MethodNotAllowedHandler)

			return
		}
	}

	c.handlers = a.combineHandlers()

	serveError(c, a.config.NotFoundHandler)
}

func errorHandler(c *Ctx) error {
	return c.Status(500).SendString(_500)
}

func serveError(c *Ctx, handler HandlerFunc) {
	err := c.Next()

	if c.writermem.Written() {
		return
	}

	_ = handler(c)
	_ = err
}

func redirectTrailingSlash(c *Ctx) {
	req := c.Request
	p := req.URL.Path
	if prefix := path.Clean(c.Request.Header.Get("X-Forwarded-Prefix")); prefix != "." {
		prefix = regSafePrefix.ReplaceAllString(prefix, "")
		prefix = regRemoveRepeatedChar.ReplaceAllString(prefix, "/")

		p = prefix + "/" + req.URL.Path
	}
	req.URL.Path = p + "/"
	if length := len(p); length > 1 && p[length-1] == '/' {
		req.URL.Path = p[:length-1]
	}

	redirectRequest(c)
}

func redirectFixedPath(c *Ctx, root *node, trailingSlash bool) bool {
	req := c.Request
	rPath := req.URL.Path

	if fixedPath, ok := root.findCaseInsensitivePath(cleanPath(rPath), trailingSlash); ok {
		req.URL.Path = bytesconv.BytesToString(fixedPath)
		redirectRequest(c)
		return true
	}
	return false
}

func redirectRequest(c *Ctx) {
	req := c.Request
	//rPath := req.URL.Path
	rURL := req.URL.String()

	code := http.StatusMovedPermanently // Permanent redirect, request with GET method
	if req.Method != http.MethodGet {
		code = http.StatusTemporaryRedirect
	}

	//debugPrint("redirecting request %d: %s --> %s", code, rPath, rURL)

	http.Redirect(c.writer, req, rURL, code)
	c.writermem.WriteHeaderNow()
}
