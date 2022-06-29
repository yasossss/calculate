package main

import (
	"fmt"
	"io"
	"log"
	"net"
	ct "project01/taskbatcher"
	pb "project01/utils/protocol"
	"sync"

	"golang.org/x/net/context"
	"google.golang.org/grpc"
)

/*
1）定义一个结构体，必须包含pb.UnimplementedGreeterServer 对象；
2）实现 .proto文件中定义的API；
3）将服务描述及其具体实现注册到 gRPC 中；
4）启动服务。
*/

const (
	PORT = ":50001"
)

// 定义空结构体，关联server服务
// 新版本 gRPC 要求必须嵌入 pb.UnimplementedConnectServer 结构体
type server struct {
	pb.UnimplementedConnectServer
}

// 实现.proto文件中定义SayHello方法
func (s *server) SayHello(ctx context.Context, in *pb.HelloRequest) (*pb.HelloReply, error) {
	log.Println("request: ", in.Name, in.Data)
	return &pb.HelloReply{Message: "Hello " + in.Name + in.Data.Message}, nil
}

//实现GetResults方法 双向流
//注意：服务端和客户端的这两个Stream 是独立的。
/*
 1. 建立连接 获取client
 2. 通过client调用方法获取stream
 3. 开两个goroutine（使用 chan 传递数据） 分别用于Recv()和Send()
 3.1 一直Recv()到err==io.EOF(即客户端关闭stream)
 3.2 Send()则自己控制什么时候Close 服务端stream没有close 只要跳出循环就算close了。 具体见https://github.com/grpc/grpc-go/issues/444
*/
func (s *server) GetResults(stream pb.Connect_GetResultsServer) error {
	// TODO
	var (
		waitGroup sync.WaitGroup
		msgCh     = make(chan []int32)
	)

	waitGroup.Add(1)
	go func() {
		defer waitGroup.Done()
		for arr := range msgCh {

			//获取客户端的数据v，传给calTask进行计算
			calTask := &ct.CalTask{
				Data: arr,
			}
			err := stream.Send(&pb.EchoResponse{
				Max: <-calTask.GetMax(),
				Min: <-calTask.GetMin(),
				Avg: <-calTask.GetAverage(),
			})

			if err != nil {
				fmt.Println("Send error:", err)
				continue
			}
		}
	}()

	waitGroup.Add(1)
	go func() {
		defer waitGroup.Done()
		for {
			req, err := stream.Recv()
			if err == io.EOF {
				break
			}
			if err != nil {
				log.Fatalf("recv error:%v", err)
			}
			fmt.Printf("Recved : %v\n", req.GetArray())
			msgCh <- req.GetArray()
		}
		close(msgCh)
	}()
	waitGroup.Wait()

	//返回nil表示已经完成响应
	return nil
}

func main() {
	lis, err := net.Listen("tcp", PORT) //port与client发起dial的一致
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	s := grpc.NewServer()
	pb.RegisterConnectServer(s, &server{}) //注册服务到grpc的sever端，RegisterConnectServer是proto里service生成

	log.Println("rpc服务已经开启")

	s.Serve(lis) //建立连接，开始服务
}
