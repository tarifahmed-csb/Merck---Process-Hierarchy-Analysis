package main

// import (
// 	"log"

// 	"github.com/northwesternmutual/grammes"
// )

// func main() {
// 	// Creates a new client with the localhost IP.
// 	client, err := grammes.DialWithWebSocket("ws://172.31.64.105:8182")
// 	if err != nil {
// 		log.Fatalf("Error while creating client: %s\n", err.Error())
// 	}

// 	_ = client
// }

import (
	"fmt"

	gremlingo "github.com/apache/tinkerpop/gremlin-go/v3/driver"
)

func main() {
	//MODIFYING TO USE API DIRECTLY FROM THE CODE AS OPPOSED TO HTTP

	//Goal 1: single out only the hierarchy (process) from the output
	//note output only contains one hierarchy/process
	name := "test"
	output, err := ModelData(name)
	if err != nil {
		fmt.Println("error")
	}
	// fmt.Println(output.Hierarchy)

	//Goal 2: Store the name of the process, the stages, and any measures
	processName := output.Hierarchy.Process
	fmt.Println("Here is process name:" + processName)

	process_stages := output.Hierarchy.Stages

	//this outer loop traverses stages
	for _, stage := range process_stages {
		//fmt.Println("Element at index", i, ":", value)
		//this outer loop traverses stages
		operations := stage.Operations

		//this loop traverses operations
		for _, operation := range operations {
			//fmt.Println("Element at index", i, ":", value)
			actions := operation.Actions
			//this loop traverses operations
			for _, action := range actions {
				//fmt.Println("Element at index", i, ":", value)
				//insertAction(action)
				measures := action.Measures

				for _, measure := range measures {
					//insertMeasure(measure)
				}
			}

		}
	}

	// Creating the connection to the server.
	driverRemoteConnection, err := gremlingo.NewDriverRemoteConnection("wss://neptunedbinstance-6sb274y70mqt.c7a08meoiuv1.us-east-2.neptune.amazonaws.com:8182/gremlin",
		func(settings *gremlingo.DriverRemoteConnectionSettings) {
			settings.TraversalSource = "g"
		})
	if err != nil {
		fmt.Println(err)
		return
	}
	// Cleanup
	defer driverRemoteConnection.Close()

	// Creating graph traversal
	g := gremlingo.Traversal_().WithRemote(driverRemoteConnection)

	result := g.AddV("test").Property("test", "test2")

	fmt.Println(result)
	fmt.Println(err)

	// Perform traversal
	results, err := g.V().ValueMap().ToList()
	if err != nil {
		fmt.Println(err)
		return
	}
	// Print results
	fmt.Println("printing here")
	for _, r := range results {
		fmt.Println(r.GetString())
	}
}
