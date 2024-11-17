package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	gremlingo "github.com/apache/tinkerpop/gremlin-go/v3/driver"
)

// "ws://localhost:8182/gremlin" to connect to local gremlin server; ws is websockets used for TCP connections
const database_url = "ws://localhost:8182/gremlin"

func main() {

	// Creating the connection to the server.
	driverRemoteConnection, err := gremlingo.NewDriverRemoteConnection(database_url,
		func(settings *gremlingo.DriverRemoteConnectionSettings) {
			settings.TraversalSource = "g"
		})
	if err != nil {
		fmt.Println(err)
		return
	}

	// Deffered Cleanup; Will be called when this function (main) reaches the end
	defer driverRemoteConnection.Close()

	// Creating graph traversal, this traverses the graph and initializes a traversal object
	//which is used to navigate, query, & manipulate the graph data stored remotely
	g := gremlingo.Traversal_().WithRemote(driverRemoteConnection)

	//initialize a reader (takes in user input)
	reader := bufio.NewReader(os.Stdin)

	//start for demo program
	for {
		printMenu()

		input, _ := reader.ReadString('\n') //reads string until \n
		// input = strings.TrimSpace(input)
		input = input[:len(input)-2] // Remove newline character & \r char

		choice, err := strconv.Atoi(input)

		if err != nil || choice < 1 || choice > 6 {
			fmt.Println("Invalid input. Please enter a number between 1 and 5.")
			continue
		}

		switch choice {
		case 1:
			fmt.Println("You chose Option 1")
			// Ask the user to input a name for a process
			fmt.Print("Enter a name for the process: ")
			processName, _ := reader.ReadString('\n')
			processName = processName[:len(processName)-2] // Remove newline character & /r
			insertNewProcesses(g, processName)
		case 2:
			fmt.Println("You chose Option 2")
			getAllProcesses(g)

		case 3:
			fmt.Println("You chose Option 3")
			fmt.Println("Enter the name of the Process")
			processName, _ := reader.ReadString('\n')
			processName = processName[:len(processName)-2] // Remove newline character
			getAllStages(g, processName)
		case 4:
			fmt.Println("You chose Option 4")
			fmt.Println("Enter the name of the Process")
			processName, _ := reader.ReadString('\n')
			processName = processName[:len(processName)-2] // Remove newline character
			getAllChildren(g, processName)
		case 5:
			fmt.Println("You chose Option 5")
			fmt.Println("Enter the name of the Process")
			processName, _ := reader.ReadString('\n')
			processName = processName[:len(processName)-2] // Remove newline character
			getAllMeasures(g, processName)
		case 6:
			fmt.Println("You chose Option 6")
			fmt.Println("Enter the name of the Process")

			processName, _ := reader.ReadString('\n')
			processName = processName[:len(processName)-2] // Remove newline character
			getAllResults(g, processName)
		}

		// Ask the user if they want to continue
		fmt.Print("Do you want to enter another menu option? (y/n): ")
		again, _ := reader.ReadString('\n')
		again = again[:len(again)-2] // Remove newline character

		if again != "y" {
			break
		}
	}

}

func printMenu() {
	fmt.Println("Choose an action (1-5):")
	fmt.Println("1. Insert a new process")
	fmt.Println("2. Query all Processes")
	fmt.Println("3. Query all Stages of a given Process")
	fmt.Println("4. Query all children of a given Process")
	fmt.Println("5. Query all measures for a given Process")
	fmt.Println("6. Query all results for a given Process")
	fmt.Print("Enter your choice: ")

}

// Create a new process hierarchy using the ModelData function.
// This function calls a helper function that inserts the structured data into a graph database
func insertNewProcesses(g *gremlingo.GraphTraversalSource, name string) {
	fmt.Println("Inserting new Process: " + name + "\nLoading data...")

	start := time.Now()

	//creating a random material num for the current process
	matnr := StringReverse(fmt.Sprintf("%v", time.Now().Unix()))

	//creating data for the process
	process, err := ModelDataParent(name, matnr)
	if err != nil {
		log.Fatal(err)
	}

	//Store process, the stages, operations, & actions and any measures into variables

	//This will keep track of the time it requires to upload the data to database

	// a single ModelData(name) process will output various batches of a single parentRawMaterialNum(corresponds to the process name) (designated by unique parentBatchID )
	// this process has 4 parent raw materials named childMaterialName (which are all differentiated with childBatchID)(rawmat-1, rawmat-2...)
	// this means that the 4 raw-mats will be OUTPUTTED BY 4 other unique processes that will also have its own 4 parent raw materials
	// EACH child INPUT RAW MATERIAL IS UNIQUE, given by the childBatchID
	// Then ModelData(childInput) will produce 4 more raw-mat inputs with unique childBatchIDS

	// so generally speaking, process A will output material 123, & process A will have inputs (retrieved from raw materials) materials 345, 567
	// these inputs (materials 345, 456) will be outputted from their respective processes B & C that will also have inputs (678, 890)

	// the way to differentiate through processes is by inputs & output meta tag(attribute)
	// the way to differentiate the outputs (childmaterials) is with the childBatchID

	// so iterating through raw mat, we will be creating n number of parent processes depending on the number of unique childMaterialName:
	// e.g childMatName: RawMat-0 --> new ModelData(RawMat-0) which has new unique inputs (various childMaterialName) & outputs a new parentMaterialNum

	rawmats := process.RawMaterials

	//inputs will store k:childMaterialName -> v:childMaterialNum
	inputs := getInputsMap(rawmats)
	// our first few loops in raw mat, we will be able to extract the unique inputs to the current process (childMaterialNum)
	// along with their name (childMaterialName), and the output (parentMaterialNum)

	//first we insert the bottom most process which has the name the user has inputted
	//remember this will out contain stages, actions, measures, results etc.
	//so we must pass in the orignial ModelData (named process)
	insertNewProcess(g, name, process)

	//With the inputs we collected (k:name, v: num) of the current process, we can go ahead and create unique processes for each input

	for inputName, inputNum := range inputs {
		newProcess, err := ModelDataParent(inputName, inputNum)
		if err != nil {
			log.Fatal(err)
		}
		insertNewProcess(g, inputName, newProcess)
		//within here once a new process has been created, then we do:
		//insertInputEdgeForProcess(newProcess, originalProcess)
		//to connect an input edge from the new process to the original process
	}

	//TODO figure a way to connect the new ModelDataParent Process as the input to original process

	// rawmats := output.RawMaterials

	// for _, rawmat := range rawmats {
	// 	//err = insertRawMat(g, rawmat)
	// }

	//from a process give all measures,
	//linking raw materials to processes (jumping from one process to another)
	//

	// result, err := g.AddE("has").From(g.V("test-1")).To(g.V("test-1-1-1-M1")).Next()

	// fmt.Println("here is result:")

	// fmt.Println(result)
	// fmt.Println("here is err")

	// fmt.Println(err)

	// Perform traversal
	// results, err := g.V().Limit(5).ToList()
	// if err != nil {
	// 	fmt.Println(err)
	// 	return
	// }
	// // Print results
	// fmt.Println("printing here")
	// for _, r := range results {
	// 	fmt.Println(r.GetString())
	// }
	end := time.Now()

	fmt.Println("Database Updated!")
	elapsedTime := end.Sub(start)
	fmt.Println("Elapsed Time:", elapsedTime)
}

func getInputsOutput(rawmats []RawMaterials) ([]string, string) {
	inputs := make(map[string]string)
	output := rawmats[0].ParentMaterialNum
	for _, rawmat := range rawmats {
		inputName := rawmat.ChildMaterialName
		inputNum := rawmat.ChildMaterialNum

		//if inputName exists in map then break; we have reached all possible inputs
		if _, exists := inputs[inputName]; exists {
			break
		} else {
			//else if key does not exist then append kv pair because it signifies a new input
			inputs[inputName] = inputNum
		}
	}
	var inputsArray []string
	for _, inputNum := range inputs {
		inputsArray = append(inputsArray, inputNum)
	}

	return inputsArray, output
}

func getInputsMap(rawmats []RawMaterials) map[string]string {
	inputs := make(map[string]string)
	count := 0
	for _, rawmat := range rawmats {
		inputName := rawmat.ChildMaterialName
		inputNum := rawmat.ChildMaterialNum
		child := rawmat.ChildBatchID
		println("here is input num:" + inputNum + "\n&input name:" + inputName + "& child" + child)

		if count == 15 {
			break
		}
		count++
	}
	for _, rawmat := range rawmats {
		inputName := rawmat.ChildMaterialName
		inputNum := rawmat.ChildMaterialNum
		println("here is input num:" + inputNum)

		//if inputName exists in map then break; we have reached all possible inputs
		if _, exists := inputs[inputName]; exists {
			break
		} else {
			//else if key does not exist then append kv pair because it signifies a new input
			inputs[inputName] = inputNum
		}
	}

	return inputs
}

func insertNewProcess(g *gremlingo.GraphTraversalSource, name string, process ModelOutput) {
	//first get inputs & outputs for the new process
	inputs, output := getInputsOutput(process.RawMaterials)
	//Here we store process into the DB
	processName, err := insertProcessDB(g, name, inputs, output)
	if err != nil {
		log.Fatal(err)
	}

	//Here we iterate through processes to store the next levels in the hierarchy
	process_stages := process.Hierarchy.Stages

	//this outer loop traverses stages
	for _, stage := range process_stages {
		//this inserts the stage along with any connected measures
		stageID, err := insertStage(g, stage)
		if err != nil {
			log.Fatal(err)
		}
		//this connects newly inserted stage with the process
		err = edgeProcessStage(g, processName, stageID)
		if err != nil {
			log.Fatal(err)
		}
		operations := stage.Operations

		//this loop traverses operations
		for _, operation := range operations {
			//this inserts the operation along with any connected measures
			operationID, err := insertOperation(g, operation)
			if err != nil {
				log.Fatal(err)
			}
			//this connects newly inserted operation with the stage
			err = edgeStageOperation(g, stageID, operationID)
			if err != nil {
				log.Fatal(err)
			}

			actions := operation.Actions

			//this loop traverses actions
			for _, action := range actions {
				//this inserts the action along with any connected measures
				actionID, err := insertAction(g, action)
				if err != nil {
					log.Fatal(err)
				}
				//this connects newly inserted action with operation
				err = edgeOperationAction(g, operationID, actionID)
				if err != nil {
					log.Fatal(err)
				}
				measures := action.Measures

				//this loop traverses measures
				for _, measure := range measures {

					measureID, err := insertMeasure(g, measure)
					if err != nil {
						log.Fatal(err)
					}
					err = edgeActionMeasure(g, actionID, measureID)
					if err != nil {
						log.Fatal(err)
					}
				}

			}
		}
	}
	//Now iterate through x_paths & adds the vertexes
	//Plus corresponding edges from measure to xpath
	x_paths := process.Xpath

	for _, x_path := range x_paths {
		err := insertXPathnEdge(g, x_path)
		if err != nil {
			log.Fatal(err)
		}

	}

	ress := process.Results

	for _, result := range ress {
		err = insertResult(g, result)
		//fmt.Println("in here")
		if err != nil {
			log.Fatal(err)
		}
		// fmt.Println(result.Result)
		// fmt.Println(result.ResultName)
		//RESULT NAME HOLDS MEAURE ID

		//insertXPath(g, x_path)
	}
}
