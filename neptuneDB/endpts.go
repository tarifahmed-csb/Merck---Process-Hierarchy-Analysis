package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"time"

	"github.com/labstack/echo/v4"

	"github.com/gocarina/gocsv"
)

// return JSON or CSV output based on encoding selected
func OutputEncode[T any](c echo.Context, fileText, encode string, input T) error {

	switch {
	case strings.HasPrefix(encode, "csv"):

		fileName := fmt.Sprintf("%v_%v.csv",
			fileText,
			time.Now().Format("20060102150405"),
		)

		outputFile, err := ioutil.TempFile("", fileName)
		if err != nil {
			return (c.JSON(http.StatusBadRequest, err.Error()))
		}
		defer outputFile.Close()
		gocsv.MarshalFile(input, outputFile)

		if encode == "csv" {
			return c.File(outputFile.Name())
		} else {
			return c.Attachment(outputFile.Name(), fileName)
		}

	default:
		return (c.JSON(http.StatusOK, input))
	}
}

// Handle calls to /version
// @Summary      API Version number
// @Description  Returns information about a given hierarchy element or measurement
// @Tags         metadata
// @Produce      json
// @Router       /version [get]
func versionEndpt(c echo.Context) error {

	fmt.Printf("%v\n", c.Request().URL)

	return (c.JSON(http.StatusOK, "Merck DB Model - "+version))
}

// Handle calls to /dataModel
// @Summary      Returns data model for Merck DB Modeling Project
// @Description  Returns JSON with 5 tables - Product Hierarchy, Xpath Connections, MetaData, Raw Materials, Data
// @Tags         dataModels
// @Produce      json
// @Produce      plain
// @Param        name path string true "name of product"
// @Param        encode query string false "output encoding of ('json','csv','csvFile')" default(json)
// @Router       /v1/dataModel/{name} [get]
func dataModelEndpt(c echo.Context) error {

	// log request url (aids in troubleshooting frontend and backend)
	fmt.Printf("%v\n", c.Request().URL)

	// get query parameters from input URL
	name := c.Param("name")
	matNum := c.Param("matNum")
	encode := c.QueryParam("encode")

	// run model
	outPut, err := ModelDataParent(name, matNum)
	if err != nil {
		fmt.Println(err)
		return (c.JSON(http.StatusBadRequest, err.Error()))
	}

	// return model output to user
	OutputEncode(c, "dbModel", encode, outPut)
	return err
}
