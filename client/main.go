package main

import (
	"fmt"
	"io"
	"log"
	"sync"

	pb "project01/utils/protocol"

	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

const (
	address     = "localhost:50001"
	defaultName = "client01"
)

/*
1）首先使用 grpc.Dial() 与 gRPC 服务器建立连接；
2）使用 pb.NewConnectClient(conn)获取客户端；
3）通过客户端调用ServiceAPI方法client.SayHello。
*/
func main() {
	//client端主动发起grpc连接，dial对方
	conn, err := grpc.Dial(address, grpc.WithTransportCredentials(insecure.NewCredentials()), grpc.WithBlock())
	if err != nil {
		log.Fatalf("dial to server with error : %v", err)
	}

	defer conn.Close() //一定不要忘记关闭连接

	// 获取客户端
	client := pb.NewConnectClient(conn)
	if client == nil {
		log.Fatalf("connect to server with error : %v", err)
	}
	fmt.Println("连接成功")
	// GetResults(client)
	GetResults1(client)
}

// 双向流
/*
1. 建立连接 获取client
2. 通过client获取stream
3. 开两个goroutine 分别用于Recv()和Send()
	3.1 一直Recv()到err==io.EOF(即服务端关闭stream)
	3.2 Send()则由自己控制
4. 发送完毕调用 stream.CloseSend()关闭stream 必须调用关闭 否则Server会一直尝试接收数据 一直报错...
*/
func GetResults(client pb.ConnectClient) {

	var waitGroup sync.WaitGroup
	// reqCh := make(chan pb.Request)
	reqsCh := make(chan *pb.GrpcRequest)
	// 调用方法获取stream
	stream, err := client.GetResults(context.Background())
	if err != nil {
		panic(err)
	}
	// 开两个goroutine 分别用于Recv()和Send()
	waitGroup.Add(1)
	go func() {
		defer waitGroup.Done()
		for {
			rep, err := stream.Recv()
			if err == io.EOF {
				fmt.Println("Server Closed")
				break
			}
			if err != nil {
				continue
			}
			fmt.Println("最大值:", rep.GetMax(), "最小值:", rep.GetMin(), "平均值:", rep.GetAvg())
		}
	}()

	waitGroup.Add(1)
	go func() {
		defer waitGroup.Done()
		for i := 0; i < 4; i++ {
			arr := [][]int32{{1, 2, 3, 4, 5}, {6, 7, 8, 9, 10}, {11, 12, 13, 14, 15}, {-1, 2, 100, -5}}
			req := pb.Request{
				Id:   int32(i),
				Data: arr[i],
			}
			fmt.Println("req:", req.Id, "req.Data:", req.Data)
			reqs := pb.GrpcRequest{
				Reqs: []*pb.Request{&req},
			}
			fmt.Println("reqs:", reqs.Reqs)
			// TODO
			reqsCh <- &reqs
		}
		reqs := <-reqsCh
		err := stream.Send(&pb.GrpcRequest{Reqs: reqs.Reqs})
		fmt.Println("Send()执行完") // =========================无法进入
		if err != nil {
			log.Printf("send error:%v\n", err)
		}

		// 发送完毕关闭stream
		err1 := stream.CloseSend()
		if err1 != nil {
			log.Printf("send error:%v\n", err)
			return
		}
	}()
	waitGroup.Wait()
}

func GetResults1(client pb.ConnectClient) {

	var waitGroup sync.WaitGroup
	// reqCh := make(chan pb.Request)
	// reqsCh := make(chan *pb.GrpcRequest)
	// 调用方法获取stream
	stream, err := client.GetResults1(context.Background())
	if err != nil {
		panic(err)
	}
	// 开两个goroutine 分别用于Recv()和Send()
	waitGroup.Add(1)
	go func() {
		defer waitGroup.Done()
		for {
			rep, err := stream.Recv()
			if err == io.EOF {
				fmt.Println("Server Closed")
				break
			}
			if err != nil {
				continue
			}
			fmt.Println("最大值:", rep.GetMax(), "最小值:", rep.GetMin(), "平均值:", rep.GetAvg())
		}
	}()

	waitGroup.Add(1)
	go func() {
		defer waitGroup.Done()
		for i := 0; i < 4; i++ {
			arr := [][]int32{{1, 2, 3, 4, 5}, {6, 7, 8, 9, 10}, {11, 12, 13, 14, 15}, {-1, 2, 100, -5}}
			err := stream.Send(&pb.Request{
				Id:   int32(i),
				Data: arr[i],
			})
			if err != nil {
				log.Printf("send error:%v\n", err)
			}
			// TODO
		}
		// reqs := <-reqsCh
		// err := stream.Send(&pb.GrpcRequest{Reqs: reqs.Reqs})
		fmt.Println("Send()执行完") // =========================无法进入
		if err != nil {
			log.Printf("send error:%v\n", err)
		}

		// 发送完毕关闭stream
		err1 := stream.CloseSend()
		if err1 != nil {
			log.Printf("send error:%v\n", err)
			return
		}
	}()
	waitGroup.Wait()
}
