package main

import (
	"fmt"
	"time"

	ct "github.com/yasossss/calculate/taskbatcher"
)

type GrpcRequest struct {
	Reqs []Request
}

type Request struct {
	Id   int
	Data []int32
}

type Response struct {
	Id      int
	Data    []int32
	Max     int32
	Min     int32
	Average int32
}

func (rep *Response) String() string {
	return fmt.Sprintf(`response: {id: %d, data: %v, max: %d, min: %d, aver: %d}`, rep.Id, rep.Data, rep.Max, rep.Min, rep.Average)
}

func main() {
	reqCh := make(chan Request)
	rspCh := make(chan Response)
	stopCh := make(chan struct{})

	// 5 means numbers of goroutine
	for i := 0; i < 5; i++ {
		workerId := i
		go worker(workerId, reqCh, rspCh, stopCh)
	}

	// len(Request) = len(Response) = 4
	rsps := make([]Response, 0, 4)
	go func() {
		for {
			select {
			case rsp := <-rspCh:
				rsps = append(rsps, rsp)
				fmt.Println(rsp.String())
				if len(rsps) == 4 {
					close(stopCh)
					return
				}
			default:
				time.Sleep(200 * time.Millisecond)
			}
		}
	}()

	for i := 0; i < 4; i++ {
		arr := [][]int32{{1, 2, 3, 4, 5}, {6, 7, 8, 9, 10}, {11, 12, 13, 14, 15}, {-1, 2, 100, -5}}
		req := Request{
			Id:   i,
			Data: arr[i],
		}
		reqCh <- req
	}

	// wait for all goroutine exit
	select {
	case <-stopCh:
		fmt.Printf("finish task total: %+v\n", len(rsps))
	}

	// to see if all goroutine exit
	time.Sleep(time.Second * 2)
}

func worker(workId int, reqCh chan Request, rspCh chan Response, stopCh chan struct{}) {
	for {
		select {
		case req := <-reqCh:
			fmt.Printf("worker %d is working for task %d\n", workId, req.Id)
			calTask := ct.CalTask{
				Data: req.Data,
			}
			rspCh <- Response{
				Id:      req.Id,
				Data:    req.Data,
				Max:     <-calTask.GetMax(),
				Min:     <-calTask.GetMin(),
				Average: <-calTask.GetAverage(),
			}
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
