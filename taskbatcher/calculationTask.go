package taskbatcher

import (
	"math"
)

// 计算任务 用于计算一组数的最大值、最小值、平均值
type CalTask struct {
	Data []int32
	// 结果
	Max, Min, Average, Sum int32
}

func (ct *CalTask) GetMin() int32 { // TODO 没必要chan，返回一个int32即可
	// ch := make(chan int32, 1)
	ct.Min = ct.Data[0] //假设第一个为最小值
	// go func() {
	//计算最小值
	for _, num := range ct.Data {
		if num < ct.Min {
			ct.Min = int32(math.Min(float64(num), float64(ct.Min)))
		}
	}
	// ch <- ct.Min
	// }()
	return ct.Min

}

func (ct *CalTask) GetMax() int32 {
	// ch := make(chan int32, 1)
	ct.Max = ct.Data[0] //假设第一个为最大值
	// go func() {
	//计算最大值
	for _, num := range ct.Data {
		if num > ct.Max {
			ct.Max = int32(math.Max(float64(num), float64(ct.Max)))
		}
	}
	// ch <- ct.Max
	// }()
	return ct.Max

}

func (ct *CalTask) GetAverage() int32 {
	// ch := make(chan int32, 1)
	ct.Sum = 0
	// go func() {
	//计算平均值
	for _, num := range ct.Data {
		ct.Sum += num
	}
	ct.Average = ct.Sum / int32(len(ct.Data))
	// ch <- ct.Average
	// }()
	return ct.Average

}
