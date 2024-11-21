package main

import (
	"fmt"
	"net/http"
	"os"

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

func fakeQuery(s string, b string) queryResponse {
	return queryResponse{"8ms", "measure-1-1-1-M2"}
}

type queryResponse struct {
	time   string `json:"time"`
	result string `json:"result"`
}

// Handlers use echo, refer to https://echo.labstack.com/docs for documentation

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
		// append response data to ./assets/logDat.txt
		logFile, err := os.OpenFile("./src/assets/logDat.txt", os.O_APPEND, 0644)
		if err != nil {
			fmt.Println(err)
		}
		defer logFile.Close()
		_, err = logFile.WriteString("{status: " + Response.Status + ", time: " + Response.Time + ", error: " + Response.Err + "}\n")
		if err != nil {
			fmt.Println(err)
		}

		return c.JSON(http.StatusOK, &Response)
	})

	e.POST("/query", func(c echo.Context) error {
		b := new(Req)
		if err := c.Bind(b); err != nil {
			return err
		}
		fmt.Println("Type rec : " + b.Type + "\nName : " + b.Name)
		ResultStruc := fakeQuery(b.Name, b.Type)
		var Response struct {
			Status  string `json:"status"`
			Time    string `json:"time"`
			Results string `json:"results"`
			Err     string `json:"error"`
		}
		Response.Status = "OK"
		Response.Time = ResultStruc.time
		Response.Results = ResultStruc.result
		Response.Err = "N/a"
		fmt.Println(Response)

		return c.JSON(http.StatusOK, &Response)
	})

	e.Logger.Fatal(e.Start(":1010"))
}
