package main

import (
	"fmt"
	"io"
	"log"
	"strconv"
	"strings"
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
WithTransportCredentials and insecure.NewCredentials()
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

	GetResults(client)
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
			req, err := stream.Recv()
			if err == io.EOF {
				fmt.Println("Server Closed")
				break
			}
			if err != nil {
				continue
			}
			fmt.Println("最大值:", req.GetMax(), "最小值:", req.GetMin(), "平均值:", req.GetAvg())
		}
	}()

	waitGroup.Add(1)
	go func() {
		defer waitGroup.Done()
		fmt.Println("请输入一组数(逗号间隔, exit 退出):")
		for {
			arr := []int32{}
			var inputStr string
			fmt.Scan(&inputStr)
			if inputStr == "exit" {
				break
			}
			strArray := strings.Split(inputStr, ",")

			//数组元素添加
			for _, str := range strArray {
				num, _ := strconv.Atoi(str)
				arr = append(arr, int32(num))
			}
			//发送数据给服务端计算
			err := stream.Send(&pb.Data{Array: arr})
			if err != nil {
				log.Printf("send error:%v\n", err)
			}
		}

		// 发送完毕关闭stream
		err := stream.CloseSend()
		if err != nil {
			log.Printf("send error:%v\n", err)
			return
		}
	}()
	waitGroup.Wait()
}
