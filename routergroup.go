package ursa

import (
	"math"
	"net/http"
	"path"
	"regexp"
)

var (
	// regEnLetter matches english letters for http method name
	regEnLetter = regexp.MustCompile("^[A-Z]+$")

	// anyMethods for RouterGroup Any method
	anyMethods = []string{
		http.MethodGet, http.MethodPost, http.MethodPut, http.MethodPatch,
		http.MethodHead, http.MethodOptions, http.MethodDelete, http.MethodConnect,
		http.MethodTrace,
	}
)

// IRouter defines all router handle interface includes single and group router.
type IRouter interface {
	IRoutes
	Group(string, ...HandlerFunc) *RouterGroup
}

// IRoutes defines all router handle interface.
type IRoutes interface {
	Use(...HandlerFunc) IRoutes

	Handle(string, string, ...HandlerFunc) IRoutes
	Any(string, ...HandlerFunc) IRoutes
	Get(string, ...HandlerFunc) IRoutes
	Post(string, ...HandlerFunc) IRoutes
	Delete(string, ...HandlerFunc) IRoutes
	Patch(string, ...HandlerFunc) IRoutes
	Put(string, ...HandlerFunc) IRoutes
	Options(string, ...HandlerFunc) IRoutes
	Head(string, ...HandlerFunc) IRoutes
	Match([]string, string, ...HandlerFunc) IRoutes

	//StaticFile(string, string) IRoutes
	//StaticFileFS(string, string, http.FileSystem) IRoutes
	//Static(string, string) IRoutes
	//StaticFS(string, http.FileSystem) IRoutes
}

type RouterGroup struct {
	Handlers []HandlerFunc
	basePath string
	app      *App
	root     bool
}

var _ IRouter = (*RouterGroup)(nil)

func (group *RouterGroup) Use(middlewares ...HandlerFunc) IRoutes {
	group.Handlers = append(group.Handlers, middlewares...)

	return group.returnObj()
}

func (group *RouterGroup) Group(relativePath string, middlewares ...HandlerFunc) *RouterGroup {
	return &RouterGroup{
		Handlers: group.combineHandlers(middlewares...),
		basePath: group.calculateAbsolutePath(relativePath),
		app:      group.app,
	}
}

func (group *RouterGroup) BasePath() string {
	return group.basePath
}

func (group *RouterGroup) handle(httpMethod, relativePath string, handlers ...HandlerFunc) IRoutes {
	absolutePath := group.calculateAbsolutePath(relativePath)
	handlers = group.combineHandlers(handlers...)
	group.app.addRoute(httpMethod, absolutePath, handlers...)

	return group.returnObj()
}

func (group *RouterGroup) Handle(httpMethod, relativePath string, handlers ...HandlerFunc) IRoutes {
	if matched := regEnLetter.MatchString(httpMethod); !matched {
		panic("http method " + httpMethod + " is not valid")
	}

	return group.handle(httpMethod, relativePath, handlers...)
}

func (group *RouterGroup) Post(relativePath string, handlers ...HandlerFunc) IRoutes {
	return group.handle(http.MethodPost, relativePath, handlers...)
}

func (group *RouterGroup) Get(relativePath string, handlers ...HandlerFunc) IRoutes {
	return group.handle(http.MethodGet, relativePath, handlers...)
}

func (group *RouterGroup) Delete(relativePath string, handlers ...HandlerFunc) IRoutes {
	return group.handle(http.MethodDelete, relativePath, handlers...)
}

func (group *RouterGroup) Patch(relativePath string, handlers ...HandlerFunc) IRoutes {
	return group.handle(http.MethodPatch, relativePath, handlers...)
}

func (group *RouterGroup) Put(relativePath string, handlers ...HandlerFunc) IRoutes {
	return group.handle(http.MethodPut, relativePath, handlers...)
}

func (group *RouterGroup) Options(relativePath string, handlers ...HandlerFunc) IRoutes {
	return group.handle(http.MethodOptions, relativePath, handlers...)
}

func (group *RouterGroup) Head(relativePath string, handlers ...HandlerFunc) IRoutes {
	return group.handle(http.MethodHead, relativePath, handlers...)
}

// Any registers a route that matches all the HTTP methods.
// GET, POST, PUT, PATCH, HEAD, OPTIONS, DELETE, CONNECT, TRACE.
func (group *RouterGroup) Any(relativePath string, handlers ...HandlerFunc) IRoutes {
	for _, method := range anyMethods {
		group.handle(method, relativePath, handlers...)
	}

	return group.returnObj()
}

func (group *RouterGroup) Match(methods []string, relativePath string, handlers ...HandlerFunc) IRoutes {
	for _, method := range methods {
		group.handle(method, relativePath, handlers...)
	}

	return group.returnObj()
}

const abortIndex int8 = math.MaxInt8 >> 1

func (group *RouterGroup) combineHandlers(handlers ...HandlerFunc) []HandlerFunc {
	finalSize := len(group.Handlers) + len(handlers)
	elsePanic(finalSize < int(abortIndex), "too many handlers")
	mergedHandlers := make([]HandlerFunc, finalSize)
	copy(mergedHandlers, group.Handlers)
	copy(mergedHandlers[len(group.Handlers):], handlers)
	return mergedHandlers
}

func (group *RouterGroup) calculateAbsolutePath(relativePath string) string {
	return path.Join(group.basePath, relativePath)
}

func (group *RouterGroup) returnObj() IRoutes {
	if group.root {
		return group.app
	}
	return group
}
