package nf

import "strings"

type router struct {
	roots    map[string]*_node
	handlers map[string][]HandlerFunc
}

func newRouter() *router {
	return &router{
		roots:    make(map[string]*_node),
		handlers: make(map[string][]HandlerFunc),
	}
}

// Only one * is allowed
func parsePattern(pattern string) []string {
	vs := strings.Split(pattern, "/")

	parts := make([]string, 0)
	for _, item := range vs {
		if item != "" {
			parts = append(parts, item)
			if item[0] == '*' {
				break
			}
		}
	}
	return parts
}

func (r *router) addRoute(method string, pattern string, handlers ...HandlerFunc) {
	parts := parsePattern(pattern)

	key := method + "-" + pattern
	_, ok := r.roots[method]
	if !ok {
		r.roots[method] = &_node{}
	}
	r.roots[method].insert(pattern, parts, 0)
	r.handlers[key] = handlers
}

func (r *router) getRoute(method string, path string) (*_node, map[string]string) {
	searchParts := parsePattern(path)
	params := make(map[string]string)
	root, ok := r.roots[method]

	if !ok {
		return &_node{}, nil
	}

	n := root.search(searchParts, 0)

	if n != nil {
		parts := parsePattern(n.pattern)
		for index, part := range parts {
			if part[0] == ':' {
				params[part[1:]] = searchParts[index]
			}
			if part[0] == '*' && len(part) > 1 {
				params[part[1:]] = strings.Join(searchParts[index:], "/")
				break
			}
		}

		return n, params
	}

	return root, nil
}

func (r *router) getRoutes(method string) []*_node {
	root, ok := r.roots[method]
	if !ok {
		return nil
	}
	nodes := make([]*_node, 0)
	root.travel(&nodes)
	return nodes
}

func (r *router) handle(c *Ctx) error {
	if err := c.verify(); err != nil {
		return err
	}

	node, params := r.getRoute(c.Method, c.path)
	if node != nil {
		c.params = params
		key := c.Method + "-" + node.pattern
		c.handlers = append(c.handlers, r.handlers[key]...)
		//c.handlers = append(r.handlers[key], c.handlers...)
	} else {
		return c.app.config.NotFoundHandler(c)
	}

	return c.Next()
}
