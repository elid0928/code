package concurrency

import (
	"fmt"
	"runtime"
)

func GetsystemInfo() {

	fmt.Println("打印系统信息")
	runtime.GOMAXPROCS(4)
	fmt.Printf("MaxProcess %d", runtime.NumCPU())
}

// 构建自然数序列
func GenerateNatural() chan int {
	ch := make(chan int)

	// 开启协程
	go func() {
		for i := 2; ; i++ {
			ch <- i
		}
	}()

	return ch
}

// 过滤器, 删除能被素数整除的数, 从管道丢弃
func PrimeFilter(in <-chan int, prime int) chan int {
	out := make(chan int)
	go func() {
		for {
			if i := <-in; i%prime != 0 {
				out <- i
			}
		}
	}()
	return out
}
