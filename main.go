package main

import (
	"digitalmmd-ppx-goapi-template/docs"
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"

	_ "embed"

	"github.com/labstack/echo/v4/middleware"
	echoSwagger "github.com/swaggo/echo-swagger"
)

//========================= Swagger API Header Fields =====================
// @title Golang API template
// @description golang API Template

// @contact.name Author Email
// @contact.email author@lehigh.edu

// @tag.name metadata
// @tag.description endpoints with metadata about the API

// @tag.name dataModels
// @tag.description data model endpoints

// @host localhost:8080

// ========================= End Swagger Definition =======================

//go:embed version.txt
var ver string
var version = fmt.Sprintf("ver 0.0.%v", ver)

func main() {

	docs.SwaggerInfo.Version = version
	e := echo.New()

	// optional - show version of API in log
	fmt.Printf("API Template - %v", version)

	// include if you are going to include swagger documenation
	e.GET("/swagger/*", echoSwagger.WrapHandler)

	// basepage endpoint
	e.GET("/", homeLink)

	// metadata endpoints
	e.GET("/version", versionEndpt)

	// model generation endpoints
	e.GET("/v1/dataModel/:name", dataModelEndpt)

	// report API output with detailed logs
	e.Use(middleware.Logger())
	// enable CORS support (required in Merck network)
	e.Use(middleware.CORS())
	// control port ID through environmental variable - useful for remote deploy
	e.Logger.Fatal(e.Start(":8080"))

}

func homeLink(c echo.Context) error {

	// put whatever HTML you'd like to appear on the welcome screen
	// this is totally optional

	welcome := fmt.Sprintln("<div style='width:100%;text-align:center;font-size:30px;font-family:sans-serif;margin-top:50px;'>" +
		"<p style='margin:0 200px;color:gray;border-bottom:1px solid gray;'>" +
		"Merck DB Model API" +
		"</p>" +
		"<p style='font-size:60%;color:gray;margin:0;'>" + version + "</p>" +
		"<p style='color:darkblue;font-size:180%;font-weight:400;margin-top:20px;'>" +
		"go for it!" +
		"</p>" +
		"<a href='/swagger/index.html' style='font-size:70%;color:darkgray;text-decoration:none;'>api documentation</a>" +
		"</div>")

	return c.HTML(http.StatusOK, welcome)
}
