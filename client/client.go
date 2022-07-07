package main

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"math/rand"
	"os"
	"strconv"
	"time"

	"github.com/go-kit/kit/endpoint"
	"github.com/go-kit/kit/sd"
	"github.com/go-kit/kit/sd/etcdv3"
	"github.com/go-kit/kit/sd/lb"
	kitLog "github.com/go-kit/log"
	pb "github.com/yasossss/calculate/utils/protocol"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

const (
	DefaultName = "casa"
	SayHello    = "SayHello"
	GetResults  = "GetResults"
)

type EndPointRequest struct {
	Method string
	Req    interface{}
}

func main() {

	var (
		etcdServer = "127.0.0.1:2379"   // 注册中心地址
		prefix     = "/services/hello/" // 监听的服务前缀
		ctx        = context.Background()
	)

	// etcd连接参数
	options := etcdv3.ClientOptions{
		DialTimeout:   time.Second * 3,
		DialKeepAlive: time.Second * 3,
	}

	// 连接注册中心
	client, err := etcdv3.NewClient(ctx, []string{etcdServer}, options)
	if err != nil {
		panic(err)
	}

	// 创建实例管理器, 此管理器会Watch监听etc中prefix的目录变化更新缓存的服务实例数据
	logger := kitLog.NewNopLogger()
	instancer, err := etcdv3.NewInstancer(client, prefix, logger)
	if err != nil {
		panic(err)
	}

	// 创建端点管理器， 此管理器根据Factory和监听的到实例创建endPoint并订阅instancer的变化动态更新Factory创建的endPoint
	endpointer := sd.NewEndpointer(instancer, reqFactory, logger) // reqFactory 是自定义的业务处理函数

	// (1)创建轮询负载均衡器
	balancer := lb.NewRoundRobin(endpointer)

	// (2)创建随机负载均衡器
	// balancer := lb.NewRandom(endpointer, 100)

	// 获取Endpoint用来进行grpc调用
	retry := lb.Retry(3, 3*time.Second, balancer)

	scanner := bufio.NewScanner(os.Stdin)
	fmt.Println("请输入每次生成数组的个数")
	scanner.Scan()
	num, _ := strconv.Atoi(scanner.Text())
	fmt.Println("请输入生成数组的大小")
	scanner.Scan()
	size, _ := strconv.Atoi(scanner.Text())

	var i = 0
	for {
		var reqs []*pb.Request
		var arrs [][]int32
		for j := 0; j < num; j++ {
			arr := GenRandIntArr(size) // 随机生成数组（0-100）
			arrs = append(arrs, arr)
		}

		for _, arr := range arrs {
			reqs = append(reqs, &pb.Request{Id: int32(i), Data: arr})
			i++
		}

		// 调用GetResults方法
		calRequest := EndPointRequest{
			Method: GetResults,
			Req: &pb.GrpcRequest{
				Reqs: reqs,
			},
		}
		calRsp, err := retry(ctx, calRequest)
		if err != nil {
			fmt.Println(err)
			continue
		}
		response, _ := calRsp.(*pb.GrpcResponse)
		for _, v := range response.Resps {
			fmt.Printf("id:%v\n data:%v\n max:%v min:%v avg:%v\n", v.Id, v.Data, v.Max, v.Min, v.Avg)
		}
		time.Sleep(time.Second) // 方便演示
	}
}

// 通过传入的实例地址  创建对应的请求endPoint
func reqFactory(instanceAddr string) (endpoint.Endpoint, io.Closer, error) {
	return func(ct context.Context, request interface{}) (interface{}, error) {
		fmt.Println("请求服务：", instanceAddr)
		conn, err := grpc.Dial(instanceAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
		if err != nil {
			return nil, err
		}
		defer conn.Close()

		client := pb.NewConnectClient(conn)
		ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
		defer cancel()
		// context.Background()
		req, _ := request.(EndPointRequest)
		switch req.Method {
		case GetResults:
			calRequest, _ := req.Req.(*pb.GrpcRequest)
			return client.GetResults(ctx, calRequest)
		default:
			return nil, fmt.Errorf("unsupport method: %s", req.Method)
		}
	}, nil, nil
}

// 通过传入数组的长度， 随机生成一个传入长度的数组（数字在0-100之间）
func GenRandIntArr(length int) []int32 {
	nums := make([]int32, length)
	rand.Seed(time.Now().UnixNano())
	for i := 0; i < length; i++ {
		nums[i] = int32(rand.Intn(100))
	}
	return nums
}
