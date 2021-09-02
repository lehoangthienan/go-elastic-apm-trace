package services

import (
	"context"
	"fmt"
	"net/http"
	"os"

	"github.com/labstack/echo/v4"
	"github.com/lehoangthienan/go-elastic-apm-trace/proto"
	httpCustom "github.com/lehoangthienan/go-elastic-apm-trace/utils/http"
	"github.com/sirupsen/logrus"
	"go.elastic.co/apm"
	"go.elastic.co/apm/module/apmgrpc"
	"google.golang.org/grpc"
)

func GetServiceBData(c echo.Context, ctx context.Context, log *logrus.Logger) error {
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

func GetServiceBGRPC(c echo.Context, ctx context.Context, log *logrus.Logger) error {
	// init custom span
	span, _ := apm.StartSpan(ctx, "Call http request to service B GRPC", "test")
	defer span.End()

	svcBPort := os.Getenv("SERVICE_B_GRPC_PORT")

	var conn *grpc.ClientConn
	conn, err := grpc.Dial(fmt.Sprintf(":%s", svcBPort),
		grpc.WithUnaryInterceptor(apmgrpc.NewUnaryClientInterceptor()),
		grpc.WithStreamInterceptor(apmgrpc.NewStreamClientInterceptor()),
		grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %s", err)
	}
	defer conn.Close()

	cClient := proto.NewUserServiceClient(conn)

	response, err := cClient.SayHello(ctx, &proto.UserReq{Id: 222})
	if err != nil {
		log.Fatalf("Error when calling SayHello: %s", err)
	}
	log.Printf("Response from server: %s", response)
	return c.String(http.StatusOK, "Ok!")
}
