package zapLog

import (
	"fmt"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"os"
	"path/filepath"
	"strings"
	"time"
)

func Logger() *zap.Logger{
	encoderConfig := zapcore.EncoderConfig{
		TimeKey:        "time",
		LevelKey:       "level",
		NameKey:        "logger",
		//CallerKey:      "caller",
		MessageKey:     "msg",
		//StacktraceKey:  "stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.LowercaseLevelEncoder,  // 小写编码器
		//EncodeTime:     zapcore.ISO8601TimeEncoder,     // ISO8601 UTC 时间格式
		EncodeTime:     TimeEncoder,     // ISO8601 UTC 时间格式
		EncodeDuration: zapcore.SecondsDurationEncoder,
		EncodeCaller:   zapcore.FullCallerEncoder,      // 全路径编码器
	}

	level := zap.NewAtomicLevelAt(zap.DebugLevel)
	var defaultPath string
	//if runtime.GOOS == "linux" {
	//	defaultPath = "/tmp/chong-linux.log"
	//}else {
	//	defaultPath = "\\code\\golang\\module\\cake-im\\runtime\\log\\cake-im.log"
	//}
	defaultPath = loggerPath()
	config := zap.Config{
		Level:level,
		Development:true,
		Encoding:"json",
		EncoderConfig:encoderConfig,
		InitialFields: map[string]interface{}{"serviceName": "chongliao-01"},
		OutputPaths:      []string{"stdout", defaultPath},
		ErrorOutputPaths: []string{"stderr"},
	}

	//构建日志
	logger, err := config.Build()
	if err != nil {
		panic(err)
	}
	return logger
}

func TimeEncoder(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
	enc.AppendString(t.Format("2006-01-02 15:04:05.000"))
}

//后期优化 对日志文件切割
func loggerPath() string {
	dir,_ := os.Getwd()
	logDir := dir + string(filepath.Separator) + "runtime" + string(filepath.Separator) + "logs"
	err := os.MkdirAll(logDir, 0755)
	if err != nil {
		fmt.Printf("日志目录创建失败:%s", logDir)
		panic(err)
	}
	logFile := dir + string(filepath.Separator) + "runtime" + string(filepath.Separator) +
		"logs" + string(filepath.Separator) + "cake-im.log"
	_, createErr := os.Create(logFile)
	if err != nil {
		panic(createErr)
	}
	return strings.TrimLeft(logFile, "E:")
}
