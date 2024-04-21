package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	// "github.com/coreos/etcd/clientv3/concurrency"
	clientv3 "go.etcd.io/etcd/client/v3"
	"go.etcd.io/etcd/client/v3/concurrency"
)

const (
	// 服务注册目录
	registryDir = "/registry/master"
)

// 创建一个分布式系统
func main() {

	if len(os.Args) < 2 {
		fmt.Println("Usage: 请输入id")
		return
	}

	masterId := os.Args[1]
	// mid, err := strconv.Atoi(masterId)
	// if err != nil {
	// 	fmt.Println("Usage: 请输入整数")
	// }
	etcdcli, err := clientv3.New(clientv3.Config{
		Endpoints:   []string{"http://127.0.0.1:2380"},
		DialTimeout: 5 * time.Second,
	})
	if err != nil {
		panic(err)
	}

	s, err := concurrency.NewSession(etcdcli, concurrency.WithTTL(5))
	if err != nil {
		fmt.Println("创建session失败", err)
	}
	defer s.Close()

	// 创建一个新的etcd选举
	e := concurrency.NewElection(s, registryDir)
	leaderCh := make(chan error)
	// 阻塞等待选举结果
	go func(e *concurrency.Election, leaderCh chan error) {
		err := e.Campaign(context.Background(), masterId)
		leaderCh <- err
	}(e, leaderCh)

	leaderChanged := e.Observe(context.Background())

	select {
	case resp := <-leaderChanged:
		log.Printf("watch leader chaged, leader: %s\n", string(resp.Kvs[0].Value))
	}

	for {
		select {
		case err := <-leaderCh:
			if err != nil {
				fmt.Printf("leader elect failed: %v\n", err)
				go func(e *concurrency.Election, leaderCh chan error) {
					err := e.Campaign(context.Background(), masterId)
					leaderCh <- err
				}(e, leaderCh)
			} else {
				fmt.Println("master change to leader")
				// if !m.IsLeader() {
				// 	m.BecomeLeader()
				// 	}
			}
		case resp := <-leaderChanged:
			if len(resp.Kvs) > 0 {
				fmt.Printf("watch leader change, leader: key: %s value: %s\n", string(resp.Kvs[0].Key), string(resp.Kvs[0].Value))
			}
		}
	}

}
