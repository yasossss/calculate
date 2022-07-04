// package main

// import (
// 	"encoding/json"
// 	"fmt"
// 	ct "project01/taskbatcher"
// 	"project01/taskbatcher/pool"
// 	"time"
// )

// type Request struct {
// 	Arrs []Array
// }
// type Array struct {
// 	Data []int32
// }

// func main() {
	// arrs := []int32{}
	// var inputStr string
	// for {
	// 	fmt.Println("请输入一组数(输入逗号间隔, exit 退出):")
	// 	fmt.Scan(&inputStr)
	// 	if inputStr == "exit" {
	// 		break
	// 	}
	// 	strArray := strings.Split(inputStr, ",")

	// 	//数组元素的添加
	// 	for _, str := range strArray {
	// 		num, _ := strconv.Atoi(str)
	// 		arrs = append(arrs, int32(num))
	// 	}

	// }

	// calTask := &ct.CalTask{
	// 	Data: arrs,
	// }

	// fmt.Printf("calTask: %v\n", *calTask)
	// fmt.Printf("最大值：%v; 最小值：%v; 平均值：%v\n", <-calTask.GetMax(), <-calTask.GetMin(), <-calTask.GetAverage())

	//启动协程池============================
// 	ArrChan := make(chan Array)

// 	arr1 := Array{
// 		Data: []int32{1, 2, 3, 4},
// 	}
// 	arr2 := Array{
// 		Data: []int32{5, 6, 7, 8},
// 	}
// 	request := &Request{
// 		Arrs: []Array{arr2, arr1},
// 	}
// 	go func() {
// 		for _, v := range request.Arrs {
// 			ArrChan <- v
// 		}
// 	}()

// 	// 2. 创建协程池(容量为5个worker)
// 	p := pool.NewPool(5)

// 	arr := <-ArrChan
// 	// 1. 创建一些任务
// 	calTask := &ct.CalTask{
// 		Data: arr.Data,
// 	}
// 	task := pool.NewTask(func() {

// 		<-calTask.GetMax()
// 		<-calTask.GetMin()
// 		<-calTask.GetAverage()
// 		time.Sleep(2 * time.Second)
// 		data, err := json.Marshal(calTask)
// 		if err != nil {
// 			fmt.Println("json encode stu failed,err:", err)
// 			return
// 		}
// 		// fmt.Println(data)         //字节数组形式输出
// 		fmt.Println(string(data)) //转换成字符串输出
// 	})

// 	// 3. 将这些任务交给协程池
// 	go func() {
// 		for {
// 			p.EntryChannel <- task
// 		}
// 	}()

// 	// 4. 启动pool
// 	p.Run()
// }
