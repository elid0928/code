package main

import logger "halo0201.live/code/golang/uberZap"

func main() {

	logger.Init(&logger.Conf{
		Path:    "./logs",
		Encoder: "json",
		LogConfig: &logger.LogConfig{
			MaxAge:  7,
			MaxSize: 10,
			MaxBack: 3,
		},
	})
	logger.Info("hello today is 2020-01-01")
	logger.Error("hello today is 2020-01-01, error")

}
