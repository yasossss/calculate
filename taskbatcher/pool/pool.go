package pool

import "fmt"

type Pool struct {
	// 给外部的入口
	EntryChannel chan *Task

	//Task队列
	JobsChannel chan *Task

	// 协程池中worker的最大数量
	workerNum int
}

func NewPool(cap int) *Pool {
	return &Pool{
		EntryChannel: make(chan *Task),
		JobsChannel:  make(chan *Task),
		workerNum:    cap,
	}
}

// 协程池创建一个worker ，并且让这个worker去工作
func (p *Pool) Worker(workerID int) {
	//具体的工作

	// 1. 永久的从JobsChannel中取任务 取不到任务就会一直阻塞
	for task := range p.JobsChannel {
		// 2. 一旦取到任务就去执行
		task.Execute()
		fmt.Println("workerID", workerID, "执行完了一个任务")
	}
}

// 启动协程池
func (p *Pool) Run() {
	// 1. 根据workerNum 创建worker去工作
	for i := 0; i < p.workerNum; i++ {
		//每个worker都有一个goroutine
		go p.Worker(i)
	}

	// 2. 从EntryChannel中去取task，将取到的task，发给JobsChannel
	for task := range p.EntryChannel {
		//读到task，将其传给JobsChannel
		p.JobsChannel <- task
	}
}
