package main

import (
	"fmt"
	"time"
)

func main() {
	fmt.Println("主函数开始执行!")

	Gou()
	time.Sleep(time.Second * 10)
	fmt.Println("主函数执行结束")
}

func Gou() {

	go func() {
		var count int
		for {
			count += 2
			fmt.Printf("goroutine 睡觉中: %d\n", count)
			time.Sleep(time.Second * 2)
		}
	}()
	fmt.Println("Gou 函数运行完成!!")
}
