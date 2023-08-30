package coroutine

import (
	"fmt"
	"strings"
	"testing"
	"time"
)

func TestDelete(t *testing.T) {

	a := map[string]string{
		"name": "liujiadong",
	}
	a["addr"] = "shenzhen"
	t.Logf("%v", a)
	delete(a, "name")
	t.Logf("%v", a)
	delete(a, "age")
	t.Logf("%v", a)
}

func TestSubPub(t *testing.T) {
	p := NewPublisher(100*time.Millisecond, 10)
	defer p.Close()

	all := p.Subscribe()
	golang := p.SubscribeTopic(func(v interface{}) bool {
		if s, ok := v.(string); ok {
			return strings.Contains(s, "golang")
		}
		return false
	})

	p.Publish("hello,  world!")
	p.Publish("hello, golang!")

	go func() {
		t.Log("发送信息")
		for i := 10; i < 100; i++ {
			p.Publish(fmt.Sprintf("hello, golang, %d", i))
		}
	}()
	go func() {
		for msg := range all {
			fmt.Println("all:", msg)
		}
	}()

	go func() {
		for msg := range golang {
			fmt.Println("golang:", msg)
		}
	}()

	// 运行一定时间后退出
	time.Sleep(3 * time.Second)
}
