package logger

import (
	"sync"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var (
	log  *zap.Logger
	once sync.Once
)

// Initialize sets up the logger
func Initialize() {
	once.Do(func() {
		config := zap.NewProductionConfig()
		config.EncoderConfig.TimeKey = "timestamp"
		config.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder

		var err error
		log, err = config.Build()
		if err != nil {
			panic("failed to initialize logger: " + err.Error())
		}
	})
}

// GetLogger returns the logger instance
func GetLogger() *zap.Logger {
	// make sure logger initialized
	if log == nil {
		Initialize()
	}

	return log
}

// Sync flushes any buffered log entries
func Sync() error {
	if log == nil {
		return nil
	}

	return log.Sync()
}
