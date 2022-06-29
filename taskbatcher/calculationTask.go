package taskbatcher

import (
	"math"
)

// 计算任务 用于计算一组数的最大值、最小值、平均值
type CalTask struct {
	Data []int32
	// 任务起始位置
	// 结果
	max, min, average, sum int32
}

func (ct *CalTask) GetMin() chan int32 {
	ch := make(chan int32, 1)
	ct.min = ct.Data[0]
	go func() {
		//计算最小值
		for _, num := range ct.Data {
			if num < ct.min {
				ct.min = int32(math.Min(float64(num), float64(ct.min)))
			}
		}
		ch <- ct.min
	}()
	return ch

}

func (ct *CalTask) GetMax() chan int32 {
	ch := make(chan int32, 1)
	ct.max = ct.Data[0]
	go func() {
		//计算最大值
		for _, num := range ct.Data {
			if num > ct.max {
				ct.max = int32(math.Max(float64(num), float64(ct.max)))
			}
		}
		ch <- ct.max
	}()
	return ch

}

func (ct *CalTask) GetAverage() chan int32 {
	ch := make(chan int32, 1)
	ct.sum = 0
	go func() {
		//计算平均值
		for _, num := range ct.Data {
			ct.sum += num
		}
		ct.average = ct.sum / int32(len(ct.Data))
		ch <- ct.average
	}()
	return ch

}
