package logger

type LogConfig struct {
	Encoder string // 编码器选择  json or console
	MaxAge  int    // 保存天数
	MaxSize int    // 单个文件大小，单位M
	MaxBack int    // 最多保留备份个数
}
