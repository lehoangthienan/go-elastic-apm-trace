package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"go.elastic.co/apm"
	"go.elastic.co/apm/module/apmechov4"
)

const thisServiceName = "service-b"

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Fail to load .env %v \n", err)
	}

	svcPort, okay := os.LookupEnv("PORT")
	if !okay {
		svcPort = "3001"
	}

	// Echo instance
	e := echo.New()
	// APM
	e.Use(apmechov4.Middleware())
	// Middleware
	e.Use(middleware.Logger())
	// e.Use(middleware.Recover()) // apm handle recover panic and send to APM

	// Routes
	e.GET("/go", DoRequestGet)

	// Start server
	e.Logger.Fatal(e.Start(fmt.Sprintf(":%s", svcPort)))
}

func DoRequestGet(c echo.Context) error {
	// get span context
	ctx := c.Request().Context()
	// init custom span
	span, _ := apm.StartSpan(ctx, "Queries data for service A", "test")
	defer span.End()
	return c.String(http.StatusOK, "Ok!")
}
