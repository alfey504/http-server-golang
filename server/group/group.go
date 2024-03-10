package group

import (
	"tcp_http_server.com/server/request"
	"tcp_http_server.com/server/router"
)

type Group struct {
	root *router.RouteNode
}

func CreateGroup(routeNode *router.RouteNode) Group {
	return Group{
		root: routeNode,
	}
}

func (group *Group) AddRoute(route string, handler func(req *request.Request)) {
	currentNode := group.root
	parent := currentNode
	subRoute := splitRoute(route)
	for _, r := range subRoute {
		println(r)
		parent = currentNode
		newNode, err := currentNode.GetSubRoute(r)
		if err != nil {
			println("Error -> ", err.Error())
			routeNode := router.CreateRouteNode(r, nil)
			newNode = parent.NewSubRoute(r, &routeNode)
		}
		currentNode = newNode
	}
	currentNode.SetHandler(handler)
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

func (group *Group) UseMiddleWare(middleware func(req *request.Request)) {
	group.root.AddMiddleware(middleware)
}
