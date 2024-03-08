package router

import (
	"fmt"

	"tcp_http_server.com/server/request"
)

type RouteNode struct {
	routeName  string
	handler    func(req *request.Request)
	middleware []func(req *request.Request)
	subRoute   map[string]*RouteNode
}

func CreateRouteNode(routeName string, handler func(req *request.Request)) RouteNode {
	return RouteNode{
		routeName:  routeName,
		handler:    handler,
		middleware: []func(req *request.Request){},
		subRoute:   make(map[string]*RouteNode),
	}
}

type Router struct {
	root *RouteNode
}

func CreateRouter() Router {
	return Router{
		root: nil,
	}
}

func (router *Router) AddRoute(route string, handler func(req *request.Request)) {
	if router.root == nil {
		routerNode := CreateRouteNode("/", nil)
		router.root = &routerNode
	}

	routeArr := splitRoute(route)

	currentNode := router.root
	for _, r := range routeArr {
		parent := currentNode
		currentNode = currentNode.subRoute[r]
		if currentNode == nil {
			routerNode := CreateRouteNode(r, nil)
			parent.subRoute[r] = &routerNode
			currentNode = parent.subRoute[r]
		}
	}
	currentNode.handler = handler
}

func (router *Router) AddMiddleware(route string, middleware func(req *request.Request)) {
	if router.root == nil {
		routerNode := CreateRouteNode("/", nil)
		router.root = &routerNode
	}

	routeArr := splitRoute(route)

	currentNode := router.root
	for _, r := range routeArr {
		parent := currentNode
		currentNode = currentNode.subRoute[r]
		if currentNode == nil {
			routerNode := CreateRouteNode(r, nil)
			parent.subRoute[r] = &routerNode
			currentNode = parent.subRoute[r]
		}
	}
	currentNode.middleware = append(currentNode.middleware, middleware)
}

func (router Router) ExecRoute(route string, req *request.Request) error {
	routes := splitRoute(route)

	if router.root == nil {
		return fmt.Errorf("route does not exist")
	}

	currentRouter := router.root
	for _, r := range routes {
		currentRouter = currentRouter.subRoute[r]
		if currentRouter == nil {
			return fmt.Errorf("route does not exist")
		}
		for _, middleware := range currentRouter.middleware {
			middleware(req)
		}
	}
	if currentRouter.handler == nil {
		return fmt.Errorf("route does not have a handler")
	}
	currentRouter.handler(req)
	return nil
}

func splitRoute(route string) []string {
	splitRoute := []string{}
	currentRoute := ""
	for _, r := range route {
		if r == '/' {
			splitRoute = append(splitRoute, currentRoute)
			currentRoute = ""
			continue
		}
		currentRoute += string(r)
	}
	splitRoute = append(splitRoute, currentRoute)
	return splitRoute[1:]
}
