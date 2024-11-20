package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

func main() {
	url := "http://localhost:8080/v1/dataModel/lehigh?encode=json"

	// Send GET request
	resp, err := http.Get(url)
	if err != nil {
		log.Fatalf("Error sending GET request: %v", err)
	}
	defer resp.Body.Close()

	// Check response status code
	if resp.StatusCode != http.StatusOK {
		log.Fatalf("Unexpected response status code: %d", resp.StatusCode)
	}

	// Read response body
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalf("Error reading response body: %v", err)
	}

	// Log JSON string
	fmt.Println(string(body))
}
