# http-server-golang

This a simple http server created using golang net package for learning how an http server works

## How to use the package

### Start listening
To start listening first you have to create an instance of the server with the CreateServer function and then use the ListenAndServer function
```go

func main(){
    myServer := server.CreateServer()

    if err := server.ListenAndServe(); err != nil {
        panic(err)
    }
}

```

### Adding a route 
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

    if err := server.ListenAndServe(); err != nil {
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

    if err := server.ListenAndServe(); err != nil {
        panic(err)
    }
}

```

### Sending html as text 
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

    if err := server.ListenAndServe(); err != nil {
        panic(err)
    }
}

```

### Sending html from a html file 
to parse and send an html file you can use RenderHtml function in the request struct.
you also have to provide the data as map[string]string which will then be mapped on to the template.
In the template wrap the key of the provided data wit {{ key }} for it to be replaced by its value.
 
```go 

func main(){
    myServer := server.CreateServer()

    myServer.AddRoute("/name/hello", func(req *request.Request) {
        
        data := map[string]string{
            "testData":           "Konnichiwa",
            "secondTestData":     "Second Test Data",
            "yetAnotherTestData": "Yet another test data huh ??",
        }

	    if _, err := req.RenderHtml("templates/index.html", data); err != nil {
		    println(err.Error())
	    }
    })

    if err := server.ListenAndServe(); err != nil {
        panic(err)
    }   
}

```

``` html

@<Component(main)>
<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Test Page</title>
</head>
<body>
    <h1>Test</h1>
    <div>This is a test {{ testData }}</div>
    <div>This is an other test {{ secondTestData }}</div>
    <div>This is yet an other test {{ yetAnotherTestData }} </div>
</body>
</html>
@<ComponentEnd>
    
```

### Handling different http methods
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

### 7. Adding middleware (Now deprecated in favor of Group and Group.UseMiddleWare)
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

### Creating groups 
To create a new group you can use the Server.CreateGroup function

```go

func main() {
    
    server := server.CreateServer()

    newGroup := server.CreateGroup("/name")

    newGroup.UseMiddleWare(func(req *request.Request) {
	    //.. middleware
    })

    newGroup.AddRoute("/hello", func(req *request.Request) {
	    //.. route handling
    })

    if err := server.ListenAndServe(); err != nil {
        panic(err)
    }
}

```

### Components In Html Template
You can now define component in your html fine using the @<Component(componentName)> content @<ComponentEnd> markdown 
All your html and outer boiler plate should now be a main component 
and to insert the component on an html doc you can use !<Component(componentName)>

```html

@<Component(main)>
<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Test Page</title>
</head>
<body>
    <h1>Test</h1>
    <div>This is a test {{ testData }}</div>
    <div>This is an other test {{ secondTestData }}</div>
    <div>This is yet an other test {{ yetAnotherTestData }} </div>
    !<Component(helloComponent)>
</body>
</html>
@<ComponentEnd>

@<Component(helloComponent)>
<div>
    <span>Hello there</span><br>
    !<Component(test)>
</div>
@<ComponentEnd>

@<Component(test)>
<span>Testing</span>
@<ComponentEnd>

```

### mapping over a slice in html 
you can use {% key -> variableName : <html>${variableName}</html> %} to map over a slice in the html template

``` html

@<Component(main)>
<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Test Page</title>
</head>
<body>
    <h1>Test</h1>
    <div>This is a test {{ testData }}</div>
    <div>This is an other test {{ secondTestData }}</div>
    {% value -> val : 
        <div> My name is ${val} </div> 
    %}
    <div>This is yet an other test {{ yetAnotherTestData }} </div>
    !<Component(helloComponent)>
</body>
</html>
@<ComponentEnd>

@<Component(helloComponent)>
<div>
    <span>Hello there</span><br>
    !<Component(test)>
</div>
@<ComponentEnd>

@<Component(test)>
<span>Testing</span>
@<ComponentEnd>

```
