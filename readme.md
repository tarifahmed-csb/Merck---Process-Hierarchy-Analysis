# Golang template for ppx applications

## Starting Notebook
In root directory: python -m graph_notebook.start_notebook --notebooks-dir notebook/destination/dir 

## Starting Gremlin Server
For Windows: in root directory merck: ./gremlin-server/bin/gremlin-server.bat

### Setup
Initialize code first by running the following
`go mod init`

if using GoLand make sure the `Enable Go modules integration` setting is selected under
Settings > Go > Go Modules

Then run `go mod tidy` this shoud load the appropriate libraries for use

## Run code
your api can be run using `go run .` in the terminal window

to control the port being use add following environmental variable to the command `go run .`  
In this case the server will be using port 8080 and can be accessed at the following URL
`http:\\localhost:8080`


### Build application
The `buildNew.sh` script contains commands need to build new application - including creation of the swagger
documentation and compiling for linux deployment

To use automatic versioning the uprev script needs to be initialized. 
This is done taking the following steps in the terminal  
`cd uprev`  
`go build`  
The uprev executible shold now be present in the child directory and 
will run when the deploy script is executed

### HTTP Standards
- Neptune : port 1010
- Dynamo : port 1011
- PostgreSQL : port 1012
All backends will prepare a /build and /query route for HTTP request, both will use HTTP POST method.  
The HTTP request from front-end will send a JSON in the following configuration:  
```
Headers {
    'application-type' : 'json'
    }, 
body{
    type : 'build'  
    name : name for 'build' (^)  
}
```
  
or  
  
```
Headers {
    'application-type' : 'json'
    }, 
body {
    type : 'query' (type of query; process, measures, or materials)  
    proc : '***' / process name (either wildcard for all, or the process name to query under)  
}
```  
Backend should read JSON file and if type is build, then it should route to the /build directory which will run the build function.  
And if the type is query it should route to the /query directory with the query functions.  
  
__RETURN__:  

```
Headers {
    'application-type' : 'json'
    },
body{
    status : 'OK' or 'Failed'
    time : str (how long did the query take)   
    err : reason for failure, ex. 'no such process' \ 'process already exists'
}
```
  
#### HTTP Handler Outline
```GO
package main

import (
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

type Req struct {
	Type string `json:"Type"`
	Name string `json:"Name"`
}

func enableCORS(w *http.ResponseWriter) {
	(*w).Header().Set("Access-Control-Allow-Origin", "*")
	(*w).Header().Set("Access-Control-Allow-Methods", "POST,GET")
	(*w).Header().Set("Access-Control-Allow-Headers", "Content-Type")
}

func fakeBuild(s string) string {
	return "5ms"
}

func main() {
	e := echo.New()
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},
		AllowHeaders: []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept},
		AllowMethods: []string{http.MethodGet, http.MethodHead, http.MethodPut, http.MethodPatch, http.MethodPost, http.MethodDelete},
	}))
	e.Use(middleware.Logger())

	e.POST("/build", func(c echo.Context) error {
		b := new(Req)
		if err := c.Bind(b); err != nil {
			return err
		}
		fmt.Println("Type rec : " + b.Type + "\nName : " + b.Name)
		rTime := fakeBuild(b.Name)
		var Response struct {
			Status string `json:"status"`
			Time   string `json:"time"`
			Err    string `json:"error"`
		}
		Response.Status = "OK"
		Response.Time = rTime
		Response.Err = "N/a"
		fmt.Println(Response)
		return c.JSON(http.StatusOK, &Response)
	})
	e.Logger.Fatal(e.Start(":1010"))
}
```