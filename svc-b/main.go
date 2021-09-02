package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"reflect"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/lehoangthienan/go-elastic-apm-trace/proto"
	grpcHandlers "github.com/lehoangthienan/go-elastic-apm-trace/utils/grpc"
	"go.elastic.co/apm"
	"go.elastic.co/apm/module/apmechov4"
	"go.elastic.co/apm/module/apmgrpc"
	"google.golang.org/grpc"
)

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
	e.GET("/go", wrapHandler(ResSample))

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
		res := customHandlerValue.Call([]reflect.Value{reflect.ValueOf(c), reflect.ValueOf(ctx)})
		return c.JSON(http.StatusOK, res[0].Interface())
	}
}

func ResSample(c echo.Context, ctx context.Context) error {
	// init custom span
	span, _ := apm.StartSpan(ctx, "Queries data for service A", "test")
	defer span.End()
	return c.String(http.StatusOK, "Ok!")
}
