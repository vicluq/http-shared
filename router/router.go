package router

import (
	"fmt"
	"net/http"
	"slices"
	"strings"

	mid "github.com/vicluq/http-shared/middleware"
)

type routeMap = map[string]http.Handler

type Router struct {
	basePath    string
	routes      routeMap
	middlewares []mid.Middleware // Router middleware
}

func (r *Router) GetBasePath() string {
	return r.basePath
}

func (r *Router) Register(root *http.ServeMux) {
	for path, route := range r.routes {
		root.Handle(path, route)
	}
}

func (r *Router) AddRoute(path string, handler http.Handler, middlewares ...mid.Middleware) {
	pathData := strings.Split(path, " ")
	routePath := fmt.Sprintf("%v %v%v", pathData[0], r.basePath, pathData[1])

	outRoute := handler
	mids := append(slices.Clone(r.middlewares), middlewares...)
	slices.Reverse(mids)
	for _, mid := range mids {
		outRoute = mid(outRoute)
	}

	r.routes[routePath] = outRoute
}

func (r *Router) AddMiddleware(middleware mid.Middleware) {
	r.middlewares = append(r.middlewares, middleware)
}
