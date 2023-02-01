package logger

import (
	"chat_room/utils/config"
	"go.uber.org/zap"
)

func Debug(msg string, fields ...zap.Field) {
	config.GLog.Debug(msg, fields...)
}

func Info(msg string, fields ...zap.Field) {
	config.GLog.Info(msg, fields...)
}

func Warn(msg string, fields ...zap.Field) {
	config.GLog.Warn(msg, fields...)
}

func Error(msg string, fields ...zap.Field) {
	config.GLog.Error(msg, fields...)
}

func Fatal(msg string, fields ...zap.Field) {
	config.GLog.Fatal(msg, fields...)
}
