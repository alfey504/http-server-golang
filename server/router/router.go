package router

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"

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

func (routeNode RouteNode) PrintRouteName() {
	println(routeNode.routeName)
}

func (routeNode *RouteNode) GetSubRoute(route string) (*RouteNode, error) {
	subNode, ok := routeNode.subRoute[route]
	if !ok {
		return nil, fmt.Errorf("sub route does not exist")
	}
	return subNode, nil
}

func (routeNode *RouteNode) NewSubRoute(route string, newRouteNode *RouteNode) *RouteNode {
	routeNode.subRoute[route] = newRouteNode
	return routeNode.subRoute[route]
}

func (routeNode *RouteNode) AddMiddleware(middleware func(req *request.Request)) {
	routeNode.middleware = append(routeNode.middleware, middleware)
}

func (routeNode *RouteNode) SetHandler(handler func(req *request.Request)) {
	routeNode.handler = handler
}

type Router struct {
	root *RouteNode
}

func CreateRouter() Router {
	return Router{
		root: nil,
	}
}

func (router *Router) GetRoot() *RouteNode {
	return router.root
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

	if len(routes) > 0 && routes[0] == "statics" {
		handleStaticRoute(req)
		return nil
	}

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

func (router *Router) GetNodeAtRoute(route string) *RouteNode {
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
	return currentNode
}

func handleStaticRoute(req *request.Request) {
	routes := splitRoute(req.Route)
	dir := getStaticRoute() + "/statics/" + strings.Join(routes[1:], "/")
	content, err := os.ReadFile(dir)
	if err != nil {
		fmt.Println(err.Error())
	}
	if _, err := req.Http(200, content); err != nil {
		fmt.Println(err.Error())

	}
}

func getStaticRoute() string {
	_, b, _, _ := runtime.Caller(0)
	return filepath.Join(filepath.Dir(b), "../..")
}
