package main

import (
	"tcp_http_server.com/server"
	"tcp_http_server.com/server/request"
)

func main() {
	server := server.CreateServer()
	server.AddRoute("/name/hello", func(req *request.Request) {
		switch req.Method {
		case "GET":
			GET(req)
		default:
			req.Write([]byte("Error 404"))
		}
	})
	if err := server.ListenAndServe(); err != nil {
		panic(err)
	}
}

func GET(req *request.Request) {
	if _, err := req.RenderHtml("templates/index.html"); err != nil {
		println(err.Error())
	}
}
