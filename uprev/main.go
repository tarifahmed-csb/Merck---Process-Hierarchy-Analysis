package main

import (
	"fmt"
	"os"
	"strconv"
)

func main() {
	ver, err := os.ReadFile("version.txt")

	if err == nil {

		numVer, err2 := strconv.ParseInt(string(ver), 0, 64)
		if err2 == nil {

			newVer := fmt.Sprintf("%v", numVer+1)
			err3 := os.WriteFile("version.txt", []byte(newVer), 0777)
			if err3 != nil {
				fmt.Print(err3)
			}
		}
	}
}
