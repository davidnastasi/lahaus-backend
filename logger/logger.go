package logger

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"os"
	"sync"
)

var once sync.Once
var logger *zap.Logger
var atom = zap.NewAtomicLevel()

// GetInstance get instance
func GetInstance() *zap.Logger {
	if logger == nil {
		once.Do(
			func() {
				logger = zap.New(zapcore.NewCore(getEncoder(), zapcore.Lock(os.Stdout), atom), zap.AddCaller())
			})
	}
	return logger
}

// GetAtomLevel get atom level
func GetAtomLevel() zap.AtomicLevel {
	return atom
}

func getEncoder() zapcore.Encoder {
	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	encoderConfig.TimeKey = "time"
	return zapcore.NewConsoleEncoder(encoderConfig)
}
