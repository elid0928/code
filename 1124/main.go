package main

import (
	"context"
	"os/exec"
	"sync"
	"time"

	"github.com/go-ping/ping"
	"github.com/panjf2000/ants/v2"
	"github.com/sirupsen/logrus"
)

const (
	retryCount = 3
)

var (
	wg *sync.WaitGroup
)

func main() {
	//
	wg = &sync.WaitGroup{}
	ctx, cancel := context.WithTimeout(context.Background(), 50*time.Second)
	defer cancel()

	// 协程池
	now := time.Now()
	pinerPool, err := ants.NewPoolWithFunc(1000, func(i interface{}) {
		defer wg.Done()
		PingContext(ctx, i.(string))
	})
	if err != nil {
		panic(err)
	}
	defer pinerPool.Release()
	for i := 0; i < 50000; i++ {
		wg.Add(1)
		pinerPool.Invoke("192.168.0.105")
	}
	wg.Wait()
	logrus.Infof("used time: %v,   ping %s", time.Since(now), "192.168.0.105")

}

func PingContext(ctx context.Context, addr string) bool {
	// var isSuccess bool
	// now := time.Now()
	// logrus.Infof("time: %v,   ping %s", now.String(), addr)

	for i := 0; i < retryCount; i++ {
		// ping

		select {
		case <-ctx.Done():
			logrus.Errorf("ping %s timeout", addr)
			return false
		default:
		}
		pinger, err := ping.NewPinger(addr)

		if err != nil {
			logrus.Errorf("ping %s error: %v", addr, err)
			return false
		}
		// ping -c 3 -i 0.5 -w 3 ip
		pinger.Count = 3
		pinger.Interval = 500 * time.Millisecond
		pinger.Timeout = 3 * time.Second
		// if success, return nil
		// if fail, continue
		err = pinger.Run()
		if err != nil {
			logrus.Errorf("ping %s error: %v", addr, err)
			continue
		}

		if pinger.Statistics().PacketLoss > 90 {
			logrus.Errorf("ping %s packetLoss: %v", addr, pinger.Statistics().PacketLoss)
			continue
		}
		// logrus.Infof("ping %s success, packetRecv: %d", addr, pinger.Statistics().PacketsRecv)
		return true
	}

	return false
}

func PingDetect(ctx context.Context, ipAddr string) bool {

	// 使用系统命令ping
	// ping -c 3 -i 0.5 -w 3 ip
	var err error
	for i := 0; i < retryCount; i++ {
		select {
		case <-ctx.Done():
			logrus.Errorf("ping %s timeout", ipAddr)
			return false
		default:
		}
		cmd := exec.CommandContext(ctx, "pingtest", "-c", "3", "-i", "500ms", "-t", "3s", ipAddr)
		_, err = cmd.CombinedOutput()
		if err != nil {
			// logrus.Errorf("ping %s error: %v", ipAddr, err)
			continue
		}
		// logrus.Infof("ping %s success, packetRecv: %d", ipAddr, out)
		return true
	}
	logrus.Errorf("ping %s error: %v", ipAddr, err)
	return false
}
