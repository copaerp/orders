package services

import (
	"context"
	"strings"

	"github.com/aws/aws-lambda-go/events"
)

type HandlerFunc func(context.Context, events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error)

type Route struct {
	Method  string
	Path    string
	Handler HandlerFunc
}

type Router struct {
	routes []Route
}

func NewRouter() *Router {
	return &Router{}
}

func (r *Router) Add(method, path string, handler HandlerFunc) {
	r.routes = append(r.routes, Route{
		Method:  method,
		Path:    path,
		Handler: handler,
	})
}

func (r *Router) Find(method, path string) (HandlerFunc, bool) {
	path = "/" + strings.Trim(path, "/")

	for _, route := range r.routes {
		if strings.EqualFold(route.Method, method) && route.Path == path {
			return route.Handler, true
		}
	}
	return nil, false
}
