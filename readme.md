# Starting Jupyter Notebook

## Notes for Set up

Have python 3.10 installed

Have apache tinkerpop gremlin server installed (to host a graph database locally)

### install the package

pip install graph-notebook

## Jupyter Classic Notebook

#### Enable the visualization widget

jupyter nbextension enable  --py --sys-prefix graph_notebook.widgets

# start jupyter notebook
In root directory: python -m graph_notebook.start_notebook --notebooks-dir notebook/destination/dir

## Starting Gremlin Server
For Windows: in root directory merck: ./gremlin-server/bin/gremlin-server.bat

## Golang Setup

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
