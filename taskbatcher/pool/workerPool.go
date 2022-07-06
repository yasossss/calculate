package pool

import (
	"fmt"
	"log"
	"time"

	ct "github.com/yasossss/calculate/taskbatcher"
	pb "github.com/yasossss/calculate/utils/protocol"
)

// 协程池
func WorkerPool(workerNum int, RequestNum int, in *pb.GrpcRequest) pb.GrpcResponse {
	var (
		reqCh  = make(chan *pb.Request)
		rspCh  = make(chan *pb.Response)
		stopCh = make(chan struct{})
	)
	// 5 (workerNum) means numbers of goroutine
	for i := 0; i < workerNum; i++ {
		workerId := i
		go worker(workerId, reqCh, rspCh, stopCh)
	}

	// len(Request) = len(Response) = 20 (RequestNum)
	resps := make([]*pb.Response, 0, RequestNum)
	go func() {
		for {
			select {
			case rsp := <-rspCh:
				resps = append(resps, rsp)
				if len(resps) == RequestNum {
					close(stopCh)
					return
				}
			default:
				time.Sleep(200 * time.Millisecond)
			}
		}
	}()

	for i := 0; i < RequestNum; i++ {
		req := in.Reqs[i]
		log.Print("req", req.String())
		reqCh <- req
	}

	// wait for all goroutine exit
	select {
	case <-stopCh:
		fmt.Printf("finish task total: %+v\n", len(resps))
	}

	return pb.GrpcResponse{Resps: resps}
	// // to see if all goroutine exit 看到后面的log加的，实际中并不需要
	// time.Sleep(time.Second * 2)
}

func worker(workId int, reqCh chan *pb.Request, rspCh chan *pb.Response, stopCh chan struct{}) {
	for {
		select {
		case req := <-reqCh:
			fmt.Printf("worker %d is working for task %d\n", workId, req.Id)
			calTask := &ct.CalTask{
				Data: req.Data,
			}
			rspCh <- &pb.Response{
				Id:   req.Id,
				Data: req.Data,
				Max:  <-calTask.GetMax(),
				Min:  <-calTask.GetMin(),
				Avg:  <-calTask.GetAverage()}
			time.Sleep(time.Millisecond * 100)
		case <-stopCh:
			fmt.Printf("worker %d exit\n", workId)
			return
		default:
			// fmt.Printf("worker %d can not get task, will sleep 200ms\n", workId)
			time.Sleep(100 * time.Millisecond)
		}
	}
}
