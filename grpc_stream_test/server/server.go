package main

import (
	"fmt"
	"google.golang.org/grpc"
	"grpc_learning/grpc_stream_test/proto"
	"log"
	"net"
	"sync"
	"time"
)

const PORT = ":50052"

type server struct {
	proto.UnimplementedGreeterServer
}

func (s *server) GetStream(req *proto.StreamReqData, res proto.Greeter_GetStreamServer) error {
	i := 0
	for {
		i++
		_ = res.Send(&proto.StreamResData{
			Data: fmt.Sprintf("%v", time.Now().Unix()),
		})
		time.Sleep(time.Second)
		if i > 10 {
			break
		}
	}
	return nil
}

func (s *server) PutStream(cliStr proto.Greeter_PutStreamServer) error {
	for {
		if recv, err := cliStr.Recv(); err != nil {
			fmt.Println(err)
			break
		} else {
			fmt.Println(recv.Data)
		}
	}
	return nil
}

func (s *server) AllStream(allStr proto.Greeter_AllStreamServer) error {

	group := sync.WaitGroup{}
	group.Add(2)
	go func() {
		defer group.Done()
		for {
			recv, _ := allStr.Recv()
			fmt.Println("收到客户端的消息:" + recv.Data)
		}
	}()

	go func() {
		defer group.Done()
		for {
			_ = allStr.Send(&proto.StreamResData{Data: "我是服务器"})
			time.Sleep(time.Second)
		}
	}()
	group.Wait()
	return nil
}

func main() {
	listen, err := net.Listen("tcp", PORT)
	if err != nil {
		panic(err)
	}
	newServer := grpc.NewServer()
	proto.RegisterGreeterServer(newServer, &server{})
	log.Printf("server listening at %v", listen.Addr())
	if err = newServer.Serve(listen); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
