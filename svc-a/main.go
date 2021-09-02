package main

import (
	"fmt"
	"net"
	"net/http"
	"os"
	"reflect"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/lehoangthienan/go-elastic-apm-trace/proto"
	grpcHandlers "github.com/lehoangthienan/go-elastic-apm-trace/utils/grpc"
	"github.com/lehoangthienan/go-elastic-apm-trace/utils/services"
	"github.com/sirupsen/logrus"
	"go.elastic.co/apm"
	"go.elastic.co/apm/module/apmechov4"
	"go.elastic.co/apm/module/apmgrpc"
	"google.golang.org/grpc"
)

// APM logger
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
	exitChan := make(chan error)
	go func() {
		osSignalChn := make(chan os.Signal, 1)
		log.Printf("exit by sign: %v\n", <-osSignalChn)
		exitChan <- nil
	}()

	go func() {
		err := startGRPC()
		if err != nil {
			exitChan <- err
		}
	}()

	go func() {
		err := startHTTP()
		if err != nil {
			exitChan <- err
		}
	}()

	if err := <-exitChan; err != nil {
		log.Printf("server error: %v\n", err)
	}
	log.Printf("server stopped\n")
}

func startHTTP() error {
	svcPort := os.Getenv("PORT")
	// Echo instance
	e := echo.New()
	// APM
	e.Use(apmechov4.Middleware())
	// Middleware
	e.Use(middleware.Logger())
	// e.Use(middleware.Recover()) // apm handle recover panic and send to APM

	// Routes
	e.GET("/go", wrapHandler(services.GetServiceBData))
	e.GET("/go-grpc", wrapHandler(services.GetServiceBGRPC))

	// Start server
	e.Logger.Fatal(e.Start(fmt.Sprintf(":%s", svcPort)))

	return nil
}

func startGRPC() error {
	grpcPort := os.Getenv("GRPC_PORT")
	lis, err := net.Listen("tcp", fmt.Sprintf(":%s", grpcPort))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	fmt.Printf("grpc server started on [::]:", grpcPort)

	grpcServer := grpc.NewServer(
		grpc.UnaryInterceptor(apmgrpc.NewUnaryServerInterceptor()),
		grpc.StreamInterceptor(apmgrpc.NewStreamServerInterceptor()),
	)

	userServer := grpcHandlers.Server{}
	proto.RegisterUserServiceServer(grpcServer, &userServer)

	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %s", err)
		return err
	}

	return nil
}

func wrapHandler(customHandlerFunc interface{}) echo.HandlerFunc {
	return func(c echo.Context) error {
		customHandlerValue := reflect.ValueOf(customHandlerFunc)
		ctx := c.Request().Context()
		res := customHandlerValue.Call([]reflect.Value{reflect.ValueOf(c), reflect.ValueOf(ctx), reflect.ValueOf(log)})
		return c.JSON(http.StatusOK, res[0].Interface())
	}
}
