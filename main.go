package main

import (
	"tcp_http_server.com/server"
	"tcp_http_server.com/server/request"
)

func main() {

	server := server.CreateServer()

	newGroup := server.CreateGroup("/name")
	newGroup.UseMiddleWare(func(req *request.Request) {
		req.Query = map[string]string{
			"hello": "sudalamani",
			"vasu":  "mandan",
		}
	})
	newGroup.AddRoute("/hello", func(req *request.Request) {
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
	for k, v := range req.Query {
		println(k, v)
	}
	data := map[string]interface{}{
		"testData":           "Konnichiwa",
		"secondTestData":     "Second Test Data",
		"yetAnotherTestData": "Yet another test data huh ??",
		"value":              []string{"hello", "how are you", "fine thank you"},
		"names":              []string{"Bartholomew", "Brad", "Joe", "Sudalamami"},
	}
	if _, err := req.RenderHtml("templates/index.html", data); err != nil {
		println(err.Error())
	}
}
