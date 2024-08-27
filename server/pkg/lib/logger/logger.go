package logger

import (
	"log"

	"go.uber.org/zap"
)

type Logger interface {
	Info(msg string, fields ...zap.Field)
	Warn(msg string, fields ...zap.Field)
	Error(msg string, fields ...zap.Field)
	Fatal(msg string, fields ...zap.Field)
}

type logger struct {
	zap *zap.Logger
}

func (l logger) Info(msg string, fields ...zap.Field) {
	l.zap.Info(msg, fields...)
}

func (l logger) Warn(msg string, fields ...zap.Field) {
	l.zap.Warn(msg, fields...)
}

func (l logger) Error(msg string, fields ...zap.Field) {
	l.zap.Error(msg, fields...)
}

func (l logger) Fatal(msg string, fields ...zap.Field) {
	l.zap.Fatal(msg, fields...)
}

func New() Logger {
	zap, err := config.Build(zap.AddCallerSkip(1))
	if err != nil {
		log.Fatalf("error at init zap logger err is%s.\n", err.Error())
		panic(err)
	}
	return logger{
		zap: zap,
	}
}
