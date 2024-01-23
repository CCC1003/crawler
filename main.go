package main

import (
	"crawler/log"
	"go.uber.org/zap/zapcore"
)

func main() {
	plugin, c := log.NewFilePlugin("./log.txt", zapcore.InfoLevel)
	defer c.Close()

	logger := log.NewLogger(plugin)
	logger.Info("log init end")
	logger.Error("log error")
}
