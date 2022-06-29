package main

import (
	"fmt"
	ct "project01/taskbatcher"
	"strconv"
	"strings"
)

func main() {
	arrs := []int32{}
	var inputStr string
	for {
		fmt.Println("请输入一组数(输入逗号间隔, exit 退出):")
		fmt.Scan(&inputStr)
		if inputStr == "exit" {
			break
		}
		strArray := strings.Split(inputStr, ",")

		//数组元素的添加
		for _, str := range strArray {
			num, _ := strconv.Atoi(str)
			arrs = append(arrs, int64(num))
		}

	}

	calTask := &ct.CalTask{
		Data: arrs,
	}

	fmt.Printf("calTask: %v\n", *calTask)
	fmt.Printf("最大值：%v; 最小值：%v; 平均值：%v\n", <-calTask.GetMax(), <-calTask.GetMin(), <-calTask.GetAverage())
}
