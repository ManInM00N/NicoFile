package util

import (
	"github.com/orandin/lumberjackrus"
	"github.com/sirupsen/logrus"
	"path/filepath"
)

var Log *logrus.Logger

func NewLog(prefix string) {
	Log = logrus.New()
	Log.AddHook(newRotateHook(prefix))
	return
}
func newRotateHook(prefix string) logrus.Hook {
	hook, err := lumberjackrus.NewHook(
		&lumberjackrus.LogFile{ // 通用日志配置
			Filename:   filepath.Join(prefix, "general.log"),
			MaxSize:    100,
			MaxBackups: 1,
			MaxAge:     1,
			Compress:   false,
			LocalTime:  true,
		},
		logrus.InfoLevel,
		&logrus.TextFormatter{DisableColors: true},
		&lumberjackrus.LogFileOpts{ // 针对不同日志级别的配置
			logrus.TraceLevel: &lumberjackrus.LogFile{
				Filename:   filepath.Join(prefix, "trace.log"),
				MaxSize:    10,   // 日志文件在轮转之前的最大大小，默认 100 MB
				MaxBackups: 10,   // 保留旧日志文件的最大数量
				MaxAge:     10,   // 保留旧日志文件的最大天数
				Compress:   true, // 是否使用 gzip 对日志文件进行压缩归档
				LocalTime:  true, // 是否使用本地时间，默认 UTC 时间
			},
			logrus.ErrorLevel: &lumberjackrus.LogFile{
				Filename:   filepath.Join(prefix, "error.log"),
				MaxSize:    10,   // 日志文件在轮转之前的最大大小，默认 100 MB
				MaxBackups: 10,   // 保留旧日志文件的最大数量
				MaxAge:     10,   // 保留旧日志文件的最大天数
				Compress:   true, // 是否使用 gzip 对日志文件进行压缩归档
				LocalTime:  true, // 是否使用本地时间，默认 UTC 时间
			},
		},
	)
	if err != nil {
		panic(err)
	}
	return hook
}
