package grpc

import (
	"context"
	"log"

	"github.com/lehoangthienan/go-elastic-apm-trace/proto"
)

type Server struct {
	proto.UnimplementedUserServiceServer
}

func (s *Server) SayHello(ctx context.Context, in *proto.UserReq) (*proto.User, error) {
	log.Printf("Receive message body from client: %s", in.Id)
	return &proto.User{Name: "An Le HIHI!"}, nil
}
