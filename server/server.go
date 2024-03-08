package server

import (
	"net"

	"tcp_http_server.com/server/request"
	"tcp_http_server.com/server/router"
)

type Server struct {
	router *router.Router
}

func CreateServer() Server {
	router := router.CreateRouter()
	return Server{
		router: &router,
	}
}

func (server *Server) AddRoute(route string, handler func(req *request.Request)) {
	server.router.AddRoute(route, handler)
}

func (server *Server) AddMiddleware(route string, middleware func(req *request.Request)) {
	server.router.AddMiddleware(route, middleware)
}

func (server *Server) ListenAndServe() error {
	listener, err := net.Listen("tcp", ":8080")
	if err != nil {
		return err
	}
	defer listener.Close()
	for {
		conn, err := listener.Accept()
		if err != nil {
			println("Error (server/server.go ListenerAndServe()) -> ", err.Error())
		}
		server.handleConn(&conn)
	}
}

func (server *Server) handleConn(conn *net.Conn) {
	defer (*conn).Close()

	requestString := make([]byte, 1024)
	_, err := (*conn).Read(requestString)
	if err != nil {
		println(err.Error())
	}

	requestMap := request.ParseRequest(requestString)
	request := request.CreateRequest(requestMap, conn)

	println(request.Route)

	if err != nil {
		println(err.Error())
	}

	server.router.ExecRoute(request.Route, &request)
}
