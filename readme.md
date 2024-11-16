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
Headers :  
type : 'build'  
name : name for 'build' (^)  
```
  
or  
  
```
Headers :  
type : 'query' (type of query; process, measures, or materials)  
proc : '***' / process name (either wildcard for all, or the process name to query under)  
```  
Backend should read JSON file and if type is build, then it should route to the /build directory which will run the build function.  
And if the type is query it should route to the /query directory with the query functions.  
  
__RETURN__:  

```
Headers :
staus : 'OK' or 'Failed'
time : str (how long did the query take)   
err : reason for failure, ex. 'no such process' \ 'process already exists'
```
  
#### HTTP Handler Outline
```GO
import (
    "encoding/json"
    "fmt"
    "net/http"
)
type Req struct {
    Type string
    Name string
}

func buildParse(w http.ResponseWriter, r *http.Request) {
    var build Req

    err := json.NewDecoder(r.body).Decode(&build)
    if err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }
    /* Call build func using build.Type and build.Name */
}
func queryParse(w http.ResponseWriter, r *http.Request) {
    var query Req

    err := json.NewDecoder(r.body).Decode(&query)
    if err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }
    /* Call query func using query.Type and query.Name */
}

func main(){
    srv := http.NewServeMux()
    srv.HandleFunc("/build", buildParse)
    srv.HandleFunc("/query", queryParse)
    
    err := http.ListenAndServe(":1010", srv)
    log.Fatal(err)
}
```