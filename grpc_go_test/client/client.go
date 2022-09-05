package main

import (
	"context"
	"flag"
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"grpc_learning/grpc_go_test/proto"
	"log"
)

var (
	addr = flag.String("addr", "localhost:50051", "the address to connect to")
)

func main() {
	conn, err := grpc.Dial(*addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	client := proto.NewGreeterClient(conn)
	hello, err := client.SayHello(context.Background(), &proto.HelloRequest{Name: "xtbo"})
	if err != nil {
		panic(err)
	}
	fmt.Println(hello.Message)
}
