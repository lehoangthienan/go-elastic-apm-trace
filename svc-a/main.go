package main

import (
	"fmt"
	"net/http"
	"os"

	httpCustom "github.com/lehoangthienan/go-elastic-apm-trace/utils/http"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"go.elastic.co/apm"

	"github.com/sirupsen/logrus"
	"go.elastic.co/apm/module/apmechov4"
)

var log = &logrus.Logger{
	Out:   os.Stderr,
	Hooks: make(logrus.LevelHooks),
	Level: logrus.DebugLevel,
	Formatter: &logrus.JSONFormatter{
		FieldMap: logrus.FieldMap{
			logrus.FieldKeyTime:  "@timestamp",
			logrus.FieldKeyLevel: "log.level",
			logrus.FieldKeyMsg:   "message",
			logrus.FieldKeyFunc:  "function.name", // non-ECS
		},
	},
}

func init() {
	apm.DefaultTracer.SetLogger(log)
}

func main() {
	svcPort := os.Getenv("PORT")
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
	span, _ := apm.StartSpan(ctx, "Call http request to service B", "test")
	defer span.End()

	svcBPort := os.Getenv("SERVICE_B_PORT")
	svcBHost := fmt.Sprintf("http://localhost:%s", svcBPort)

	req, err := http.NewRequest("GET", fmt.Sprintf("%s/go", svcBHost), nil)
	if err != nil {
		return nil
	}

	_, err = httpCustom.Do(req, ctx, log)

	if err != nil {
		return nil
	}

	// test panic
	// panic("non-ASCII name!")

	return c.String(http.StatusOK, "Ok!")
}
