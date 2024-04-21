package main

import (
	"context"
	"fmt"
	"net"
	"sync"
	"time"
	// "google.golang.org/appengine/socket"
)

var (
	wg *sync.WaitGroup
)

// 通过tcp连接, 并发送数据
func main() {

	timeout := 5 * time.Second
	ctx, _ := context.WithTimeout(context.Background(), timeout)
	wg = &sync.WaitGroup{}
	wg.Add(1)
	go Dial(ctx)
	wg.Wait()

	// socket.Dial()
	// socket.Options
}

func Dial(ctx context.Context) {

	defer wg.Done()
	for {
		select {
		case <-ctx.Done():
			fmt.Println("context timeout,")
			return
		default:

			dialer := net.Dialer{Timeout: 10 * time.Second}
			conn, err := dialer.Dial("tcp", "34.125.33.128:8001")
			if err != nil {
				panic(err)
			}
			conn.SetDeadline(time.Now().Add(time.Second * 10))

			defer conn.Close()
			conn.Write([]byte("SYN"))
			fmt.Println("[SENT] SYN message sent to server.")

			data := make([]byte, 1024)
			_, err = conn.Read(data)
			if err != nil {
				panic(err)
			}
			if string(data) == "SYN-ACK" {
				fmt.Println(string(data))
				conn.Write([]byte("ACK"))
				fmt.Println("[SENT] ACK message sent to server.")
			}
			return
			// time.Sleep(time.Second * 14)

		}

	}

}
