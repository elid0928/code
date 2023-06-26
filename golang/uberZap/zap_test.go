package logger

import (
	"testing"

	"github.com/sirupsen/logrus"
)

func BenchmarkZap(b *testing.B) {
	Init(&Conf{
		Path:    "/tmp",
		Encoder: "json",
	})
	logger.Info("hello world")
	// log.Info("hello world")
}

func BenchmarkLogrus(b *testing.B) {
	// runtime.MaxProcs(1)
	// runtime.GOMAXPROCS(1)
	// TODO: Initialize
	// for i := 0; i < b.N; i++ {

	logrus.Info("hello world")
	// }
}

func TestSwitch(t *testing.T) {
	var c string = "json"

	switch c {
	case "json":
		t.Log("json")
	case "console":
		t.Log("console")
	default:
		t.Log("default")
	}
}
