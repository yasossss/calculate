package pool

import (
	"fmt"
	"time"

	ct "github.com/yasossss/calculate/taskbatcher"
	pb "github.com/yasossss/calculate/utils/protocol"
)

// 协程池
func WorkerPool(workerNum int, RequestNum int, reqCh chan *pb.Request, rspCh chan *pb.Response, stopCh chan struct{}) {
	// 5 (workerNum) means numbers of goroutine
	for i := 0; i < workerNum; i++ {
		workerId := i
		go worker(workerId, reqCh, rspCh, stopCh)
	}

	// len(Request) = len(Response) = 20 (RequestNum)
	rsps := make([]*pb.Response, 0, RequestNum)
	go func() {
		for {
			select {
			case rsp := <-rspCh:
				rsps = append(rsps, rsp)
				if len(rsps) == RequestNum {
					close(stopCh)
					return
				}
			default:
				time.Sleep(200 * time.Millisecond)
			}
		}
	}()

	for i := 0; i < 20; i++ {
		req := pb.Request{Id: int32(i)}
		reqCh <- &req
	}

	// wait for all goroutine exit
	select {
	case <-stopCh:
		fmt.Printf("finish task total: %+v\n", len(rsps))
	}

	// to see if all goroutine exit
	time.Sleep(time.Second * 2)
}

func worker(workId int, reqCh chan *pb.Request, rspCh chan *pb.Response, stopCh chan struct{}) {
	for {
		select {
		case req := <-reqCh:
			fmt.Printf("worker %d is working for task %d\n", workId, req.Id)
			calTask := &ct.CalTask{
				Data: req.Data,
			}
			rspCh <- &pb.Response{Id: req.Id,
				Data: req.Data,
				Max:  <-calTask.GetMax(),
				Min:  <-calTask.GetMin(),
				Avg:  <-calTask.GetAverage()}
			time.Sleep(time.Millisecond * 100)
		case <-stopCh:
			fmt.Printf("worker %d exit\n", workId)
			return
		default:
			fmt.Printf("worker %d can not get task, will sleep 200ms\n", workId)
			time.Sleep(200 * time.Millisecond)
		}
	}
}
