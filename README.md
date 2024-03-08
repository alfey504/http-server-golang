# http-server-golang

This a simple http server created using golang net package for learning how an http server works

## How to use the package

### 1. Start listening
To start listening first you have to create an instance of the server with the CreateServer function and then use the ListenAndServer function
```go

func main(){
    myServer := server.CreateServer()
    if err := server.ListenAndServe(); err != nil {
		panic(err)
	}
}

```

### 2. Adding a route 
To add a route you can use the AddRoute() function to add a route handler to your server

```go

func main(){
    myServer := server.CreateServer()
    myServer.AddRoute("/name/hello", func(req *request.Request) {
        someText := []byte("Hello every one how are you fine thank you, OHHHH MYYYYY GAAAHHHHH !!!")
        if _, err := req.Write(someText); err != nil {
            println(err.Error())
        }
    })
    if err := myServer.ListenAndServe(); err != nil {
		panic(err)
	}
}

```

### 3. Reading query from request 
The queries in the request are stored in request.Request struct as query

``` go 

func main(){
    myServer := server.CreateServer()
    myServer.AddRoute("/name/hello", func(req *request.Request) {
        for k, v := range req.Query {
            println(k, " -> ", v)
        }
    })
    if err := myServer.ListenAndServe(); err != nil {
		panic(err)
	}
}

```

### 4. Sending html as text 
to send simple html you can use Html function in the request struct

```go 

func main(){
    myServer := server.CreateServer()
    myServer.AddRoute("/name/hello", func(req *request.Request) {
        html := []byte("<h1>Hello! Konnichiwa </h1>")
        if _, err := req.Html(html); err != nil {
            println(err.Error())
        }
    })
    if err := myServer.ListenAndServe(); err != nil {
		panic(err)
	}
}

```

### 5. Sending html from a html file 
to parse and send an html file you can use RenderHtml function in the request struct

```go 

func main(){
    myServer := server.CreateServer()
    myServer.AddRoute("/name/hello", func(req *request.Request) {
        if _, err := req.RenderHtml("templates/index.html"); err != nil {
			println(err.Error())
		}
    })
    if err := myServer.ListenAndServe(); err != nil {
		panic(err)
	}
}

```

### 6. Handling different http methods
to handle different methods you can use the method field in the Request struct

```go 

func main(){
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

```

### 7. Adding middleware 
To add a middleware you can use the AddMiddleware function in Server struct

``` go

func main() {
	server := server.CreateServer()

	server.AddMiddleware("/name", func(req *request.Request) {
		//.. middleware functions
	})

	if err := server.ListenAndServe(); err != nil {
		panic(err)
	}
}

```