package common

import (
	"math/rand"
	"time"
)

// 通过传入数组的长度， 随机生成一个传入长度的数组
func GenRandIntArr(length int) []int32 {
	nums := make([]int32, length)
	rand.Seed(time.Now().UnixNano())
	for i := 0; i < length; i++ {
		nums[i] = int32(rand.Intn(100))
	}
	return nums
}
