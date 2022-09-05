package main

import (
	"context"
	"fmt"
	"google.golang.org/grpc"
	"grpc_learning/grpc_stream_test/proto"
	"sync"
	"time"
)

func main() {
	// 服务端流模式
	dial, err := grpc.Dial("localhost:50052", grpc.WithInsecure())
	if err != nil {
		panic(err)
	}
	defer dial.Close()
	client := proto.NewGreeterClient(dial)
	stream, err := client.GetStream(context.Background(), &proto.StreamReqData{Data: "hahaha"})
	for {
		recv, err := stream.Recv()
		if err != nil {
			fmt.Println(err)
			break
		}
		fmt.Println(recv.Data)
	}

	// 客户端流模式
	putStream, _ := client.PutStream(context.Background())
	i := 0
	for {
		i++
		putStream.Send(&proto.StreamReqData{
			Data: fmt.Sprintf("发送%d", i),
		})
		time.Sleep(time.Second)
		if i > 10 {
			break
		}
	}

	// 双向流模式
	allStream, err := client.AllStream(context.Background())
	group := sync.WaitGroup{}
	group.Add(2)
	go func() {
		defer group.Done()
		for {
			recv, _ := allStream.Recv()
			fmt.Println("收到服务端的消息:" + recv.Data)
		}
	}()
	go func() {
		defer group.Done()
		for {
			_ = allStream.Send(&proto.StreamReqData{Data: "我是客户端"})
			time.Sleep(time.Second)
		}
	}()
	group.Wait()
}
