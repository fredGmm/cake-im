package zapLog

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"runtime"
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
		EncodeTime:     zapcore.ISO8601TimeEncoder,     // ISO8601 UTC 时间格式
		EncodeDuration: zapcore.SecondsDurationEncoder,
		EncodeCaller:   zapcore.FullCallerEncoder,      // 全路径编码器
	}

	level := zap.NewAtomicLevelAt(zap.DebugLevel)
	var defaultPath string
	if runtime.GOOS == "linux" {
		defaultPath = "/tmp/chong-linux.log"
	}else {
		defaultPath = "\\code\\golang\\module\\cake-im\\runtime\\log\\cake-im.log"
	}
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
