package main

import (
	"fmt"
	"os"
	"time"

	"github.com/go-ping/ping"
)

func main() {
	var arg string = "34.125.33.128"
	if len(os.Args) == 2 {
		arg = os.Args[1]
	}

	pinger, err := ping.NewPinger(arg)

	if err != nil {
		panic(err)
	}
	pinger.SetPrivileged(true)
	pinger.Count = 3 // send 3 pings
	if pinger.Privileged() {
		fmt.Println("Privileged")
	}
	pinger.Timeout = 5 * time.Second  // 5 seconds
	pinger.Interval = 1 * time.Second // 1 second

	pinger.OnFinish = func(stats *ping.Statistics) {
		fmt.Printf("\n--- %s ping statistics ---\n", stats.Addr)
		fmt.Printf("%d packets transmitted, %d packets received, %v%% packet loss\n", stats.PacketsSent, stats.PacketsRecv, stats.PacketLoss)
		fmt.Printf("round-trip min/avg/max/stddev = %v/%v/%v/%v\n", stats.MinRtt, stats.AvgRtt, stats.MaxRtt, stats.StdDevRtt)
	}
	err = pinger.Run() // blocks until finished
	stat := pinger.Statistics().PacketLoss
	if stat >= 60 {
		fmt.Println("Ping failed")
	}
	if err != nil {
		panic(err)
	}

}
