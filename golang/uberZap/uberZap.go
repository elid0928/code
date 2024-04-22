package logger

import (
	"fmt"
	"os"
	"path"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/buffer"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
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

var (
	logger *zap.Logger
	c      Conf
	_pool  = buffer.NewPool()
)

type (
	Conf struct {
		Path    string // 日志路径
		Encoder string // 编码器选择
		*LogConfig
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

func Init(conf *Conf) {
	infoPath := path.Join(conf.Path, "info")
	errPath := path.Join(conf.Path, "error")

	c = *conf
	items := []logItem{
		{
			FileName: fmt.Sprintf("%s/%s.log", infoPath, "info"),
			Level:    func(l zapcore.Level) bool { return l <= zapcore.InfoLevel },
		}, {
			FileName: fmt.Sprintf("%s/%s.log", errPath, "error"),
			Level:    func(l zapcore.Level) bool { return l >= zapcore.WarnLevel },
		},
	}
	newLogger(items, conf.LogConfig)
	return
}

func newLogger(items []logItem, conf *LogConfig) {
	var (
		cfg   zapcore.Encoder
		cores []zapcore.Core
	)

	switch c.Encoder {
	case "json":
		cfg = JsonConfig()
	case "console":
		cfg = ConsoleConfig()
	default:
		cfg = ConsoleConfig()
	}
	// 装配分割器
	for _, item := range items {
		hook := getLumberJackLogger(conf, item.FileName)
		core := zapcore.NewCore(cfg,
			zapcore.NewMultiWriteSyncer(zapcore.AddSync(hook), zapcore.AddSync(os.Stdout)),
			item.Level)
		cores = append(cores, core)
	}

	caller := zap.AddCaller()
	callerSkip := zap.AddCallerSkip(1)
	development := zap.Development()
	logger = zap.New(zapcore.NewTee(cores...), caller, development, callerSkip)
	return

}

/*
 日志分割器
 将日志文件切割， 设置大小， 最长保留时间，最多备份， 是否压缩等等
*/

func getLumberJackLogger(c *LogConfig, fileName string) *lumberjack.Logger {
	var ljL *lumberjack.Logger
	if c != nil {
		ljL = &lumberjack.Logger{
			Filename:   fileName,
			MaxSize:    c.MaxSize, // MB
			MaxBackups: c.MaxBack, // 最多保留3个备份
			MaxAge:     c.MaxAge,  // 保留7天
			Compress:   true,      // 是否压缩
			LocalTime:  true,
		}
	} else {
		ljL = &lumberjack.Logger{
			Filename:   "test.log",
			MaxSize:    10,   // MB
			MaxBackups: 3,    // 最多保留3个备份
			MaxAge:     7,    // 保留7天
			Compress:   true, // 是否压缩
			LocalTime:  true,
		}
	}
	return ljL
}

// Json 格式的输出配置
func JsonConfig() zapcore.Encoder {
	var (
		cfg = zap.NewProductionEncoderConfig()
	)

	// 时间格式自定义
	cfg.EncodeTime = func(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
		enc.AppendString(t.Format("2006-01-02 15:04:05"))
	}
	// 打印路径自定义
	cfg.EncodeCaller = func(caller zapcore.EntryCaller, encoder zapcore.PrimitiveArrayEncoder) {
		encoder.AppendString(getFilePath(caller))
	}
	// 级别显示自定义
	cfg.EncodeLevel = func(level zapcore.Level, encoder zapcore.PrimitiveArrayEncoder) {
		encoder.AppendString(level.String())
	}
	return zapcore.NewJSONEncoder(cfg)
}

// 终端输出配置
func ConsoleConfig() zapcore.Encoder {
	var (
		cfg = zap.NewProductionEncoderConfig()
	)

	// 时间格式自定义
	cfg.EncodeTime = func(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
		enc.AppendString("[" + t.Format("2006-01-02 15:04:05") + "]")
	}
	// 打印路径自定义
	cfg.EncodeCaller = func(caller zapcore.EntryCaller, encoder zapcore.PrimitiveArrayEncoder) {
		encoder.AppendString("[" + getFilePath(caller) + "]")
	}
	// 级别显示自定义
	cfg.EncodeLevel = func(level zapcore.Level, encoder zapcore.PrimitiveArrayEncoder) {
		encoder.AppendString("[" + level.String() + "]")
	}
	return zapcore.NewConsoleEncoder(cfg)
}

// getFilePath 自定义获取文件路径.
func getFilePath(ec zapcore.EntryCaller) string {
	if !ec.Defined {
		return "undefined"
	}
	buf := _pool.Get()
	buf.AppendString(ec.Function)
	buf.AppendByte(':')
	buf.AppendInt(int64(ec.Line))
	caller := buf.String()
	buf.Free()
	return caller
}
