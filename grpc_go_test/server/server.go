package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net"

	"google.golang.org/grpc"

	"grpc_learning/grpc_go_test/proto"
)

var (
	port = flag.Int("port", 50051, "The server port")
)

type Server struct {
	proto.UnimplementedGreeterServer
}

func (s *Server) SayHello(ctx context.Context, request *proto.HelloRequest) (*proto.HelloReply, error) {
	return &proto.HelloReply{
		Message: "hello " + request.Name,
	}, nil
}

func main() {
	g := grpc.NewServer()
	proto.RegisterGreeterServer(g, &Server{})
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", *port))
	log.Printf("server listening at %v", lis.Addr())
	if err = g.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
