package nf

import (
	"log"
	"net/http"
)

type RouterGroup struct {
	prefix      string
	middlewares []HandlerFunc // support middleware
	parent      *RouterGroup  // support nesting
	app         *App          // all groups share a Engine instance
}

// Group is defined to create a new RouterGroup
// remember all groups share the same Engine instance
func (group *RouterGroup) Group(prefix string) *RouterGroup {
	app := group.app
	newGroup := &RouterGroup{
		prefix: group.prefix + prefix,
		parent: group,
		app:    app,
	}
	app.groups = append(app.groups, newGroup)
	return newGroup
}

func (group *RouterGroup) addRoute(method string, comp string, handlers ...HandlerFunc) {
	verifyHandlers(comp, handlers...)
	pattern := group.prefix + comp
	log.Printf("Add Route %4s - %s", method, pattern)
	group.app.router.addRoute(method, pattern, handlers...)
}

func (group *RouterGroup) Get(pattern string, handlers ...HandlerFunc) {
	group.addRoute(http.MethodGet, pattern, handlers...)
}

func (group *RouterGroup) Post(pattern string, handlers ...HandlerFunc) {
	group.addRoute(http.MethodPost, pattern, handlers...)
}

func (group *RouterGroup) Put(pattern string, handlers ...HandlerFunc) {
	group.addRoute(http.MethodPut, pattern, handlers...)
}

func (group *RouterGroup) Delete(pattern string, handlers ...HandlerFunc) {
	group.addRoute(http.MethodDelete, pattern, handlers...)
}

func (group *RouterGroup) Patch(pattern string, handlers ...HandlerFunc) {
	group.addRoute(http.MethodPatch, pattern, handlers...)
}

func (group *RouterGroup) Head(pattern string, handlers ...HandlerFunc) {
	group.addRoute(http.MethodHead, pattern, handlers...)
}

func (group *RouterGroup) Use(middlewares ...HandlerFunc) {
	group.middlewares = append(group.middlewares, middlewares...)
}
