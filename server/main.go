// package main

// import (
// 	"context"
// 	"fmt"
// 	"io"
// 	"log"
// 	"net"
// 	"time"

// 	pb "github.com/yasossss/calculate/utils/protocol"

// 	ct "github.com/yasossss/calculate/taskbatcher"

// 	"google.golang.org/grpc"
// )

// /*
// 1）定义一个结构体，必须包含pb.UnimplementedGreeterServer 对象；
// 2）实现 .proto文件中定义的API；
// 3）将服务描述及其具体实现注册到 gRPC 中；
// 4）启动服务。
// */

// const (
// 	PORT = ":50001"
// )

// // type GrpcRequest struct {
// // 	Reqs []Request
// // }

// // type Request struct {
// // 	Id   int
// // 	Data []int32
// // }

// // type Response struct {
// // 	Id      int
// // 	Data    []int32
// // 	Max     int32
// // 	Min     int32
// // 	Average int32
// // }

// // 定义空结构体，关联server服务
// // 新版本 gRPC 要求必须嵌入 pb.UnimplementedConnectServer 结构体
// type server struct {
// 	pb.UnimplementedConnectServer
// }

// // 实现.proto文件中定义SayHello方法
// func (s *server) SayHello(ctx context.Context, in *pb.HelloRequest) (*pb.HelloReply, error) {
// 	log.Println("request: ", in.Name, in.Data)
// 	return &pb.HelloReply{Message: "Hello " + in.Name + in.Data.String()}, nil
// }

// //实现GetResults方法 双向流
// //注意：服务端和客户端的这两个Stream 是独立的。
// /*
//  1. 建立连接 获取client
//  2. 通过client调用方法获取stream
//  3. 开两个goroutine（使用 chan 传递数据） 分别用于Recv()和Send()
//  3.1 一直Recv()到err==io.EOF(即客户端关闭stream)
//  3.2 Send()则自己控制什么时候Close 服务端stream没有close 只要跳出循环就算close了。 具体见https://github.com/grpc/grpc-go/issues/444
// */

// func (s *server) GetResults(stream pb.Connect_GetResultsServer) error {
// 	// TODO
// 	var (
// 		// waitGroup sync.WaitGroup
// 		reqCh  = make(chan *pb.Request)
// 		rspCh  = make(chan *pb.Response)
// 		stopCh = make(chan struct{})
// 	)
// 	// 5 means numbers of goroutine
// 	for i := 0; i < 5; i++ {
// 		workerId := i
// 		go worker(workerId, reqCh, rspCh, stopCh)
// 	}
// 	// len(Request) = len(Response) = 4
// 	rsps := make([]pb.Response, 0, 10)
// 	// waitGroup.Add(1)
// 	go func() {
// 		// defer waitGroup.Done()
// 		for rsp := range rspCh {
// 			fmt.Printf("Id:%d Data:%v max:%v min:%v avg:%v\n", rsp.Id, rsp.Data, rsp.Max, rsp.Min, rsp.Avg)
// 			// fmt.Println(rsp.String())
// 			err := stream.Send(&pb.Response{
// 				Id:   rsp.Id,
// 				Data: rsp.Data,
// 				Max:  rsp.Max,
// 				Min:  rsp.Min,
// 				Avg:  rsp.Avg,
// 			})
// 			if err != nil {
// 				log.Printf("send error:%v\n", err)
// 			}

// 			rsps = append(rsps, *rsp)
// 			if len(rsps) == 10 {
// 				close(stopCh)
// 				return
// 			}

// 			//获取客户端的数据arr，传给calTask进行计算
// 			// calTask := &ct.CalTask{
// 			// 	Data: req.Data,
// 			// }
// 			// err := stream.Send(&pb.Response{
// 			// 	Id:   req.Id,
// 			// 	Data: req.Data,
// 			// 	Max:  <-calTask.GetMax(),
// 			// 	Min:  <-calTask.GetMin(),
// 			// 	Avg:  <-calTask.GetAverage(),
// 			// })

// 			// if err != nil {
// 			// 	fmt.Println("Send error:", err)
// 			// 	continue
// 			// }
// 		}
// 	}()

// 	// waitGroup.Add(1)
// 	go func() {
// 		// defer waitGroup.Done()
// 		for {
// 			reqs, err := stream.Recv() // 接收请求的数据
// 			if err == io.EOF {
// 				break
// 			}
// 			if err != nil {
// 				log.Fatalf("recv error:%v", err)
// 			}

// 			fmt.Printf("Recved : %v\n", reqs.Data)
// 			i := 0
// 			reqCh <- reqs
// 			i++
// 		}
// 		// close(reqCh)
// 	}()
// 	// waitGroup.Wait()
// 	// wait for all goroutine exit
// 	select {
// 	case <-stopCh:
// 		fmt.Printf("finish task total: %+v\n", len(rsps))
// 	}

// 	// to see if all goroutine exit
// 	time.Sleep(time.Second * 2)

// 	//返回nil表示已经完成响应
// 	return nil

// }

// func main() {
// 	lis, err := net.Listen("tcp", PORT) //port与client发起dial的一致
// 	if err != nil {
// 		log.Fatalf("failed to listen: %v", err)
// 	}

// 	s := grpc.NewServer()
// 	pb.RegisterConnectServer(s, &server{}) //注册服务到grpc的sever端，RegisterConnectServer是proto里service生成

// 	log.Println("rpc服务已经开启")

// 	s.Serve(lis) //建立连接，开始服务
// }

// func worker(workId int, reqCh chan *pb.Request, rspCh chan *pb.Response, stopCh chan struct{}) {
// 	for {
// 		select {
// 		case req := <-reqCh:
// 			// fmt.Printf("worker %d is working for task %d\n", workId, req.Id)
// 			calTask := ct.CalTask{
// 				Data: req.Data,
// 			}
// 			rsp := pb.Response{
// 				Id:   req.Id,
// 				Data: req.Data,
// 				Max:  <-calTask.GetMax(),
// 				Min:  <-calTask.GetMin(),
// 				Avg:  <-calTask.GetAverage(),
// 			}
// 			rspCh <- &rsp
// 			fmt.Printf("Id:%d Data:%v max:%v min:%v avg:%v\n", rsp.Id, rsp.Data, rsp.Max, rsp.Min, rsp.Avg)
// 			time.Sleep(time.Millisecond * 100)

// 		case <-stopCh:
// 			fmt.Printf("worker %d exit\n", workId)
// 			return
// 		default:
// 			// fmt.Printf("worker %d can not get task, will sleep 200ms\n", workId)
// 			time.Sleep(200 * time.Millisecond)
// 		}
// 	}
// }
