package main

import (
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

/*
// NewProduction builds a sensible production Logger that writes InfoLevel and
// above logs to standard error as JSON.
//
// It's a shortcut for NewProductionConfig().Build(...Option).

	func NewProduction(options ...Option) (*Logger, error) {
		return NewProductionConfig().Build(options...)
	}
*/

type (
	Conf struct {
		Path    string // 日志路径
		Encoder string // 编码器选择
	}
	logItem struct {
		FileName string
		Level    zap.LevelEnablerFunc
	}
	Encoder interface {
		Config() zapcore.Encoder
		WithKey(key string) Encoder
		WithField(key, val string) Encoder
		Debug(msg string)
		Debugf(format string, v ...interface{})
		Info(msg string)
		Infof(format string, v ...interface{})
		Warn(msg string)
		Warnf(format string, v ...interface{})
		Error(msg string)
		Errorf(format string, v ...interface{})
		Fatal(msg string)
		Fatalf(format string, v ...interface{})
	}
)

func main() {
	// 编译错误
	// var a string

	//
	// var _a string
	config := zap.NewProductionConfig()
	config.Encoding = "json" // or "json"
	// 设置时间格式
	config.EncoderConfig.EncodeTime = func(t time.Time, pae zapcore.PrimitiveArrayEncoder) {
		pae.AppendString("[" + t.Format("2006-01-02 15:04:05") + "]")
	}

	// 打印日志路径
	config.EncoderConfig.EncodeCaller = func(ec zapcore.EntryCaller, pae zapcore.PrimitiveArrayEncoder) {
		pae.AppendString("[" + ec.TrimmedPath() + "]")
	}

	// 日志级别显示
	config.EncoderConfig.EncodeLevel = func(l zapcore.Level, pae zapcore.PrimitiveArrayEncoder) {
		pae.AppendString("[" + l.CapitalString() + "]")
	}

	logger, _ := config.Build()
	logger.Info("info service start")
	logger.Info("info msg",
		zap.String("name", "xiaoming"),
		zap.Int("age", 18),
		zap.Bool("married", false),
	)
}
