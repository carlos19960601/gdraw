package logger

import (
	"go.uber.org/zap"
)

var logger *zap.Logger

func Init() (err error) {
	logger, err = zap.NewDevelopment(zap.AddCallerSkip(1))
	if err != nil {
		return err
	}
	return nil
}

func Info(msg string, fields ...zap.Field) {
	logger.Info(msg, fields...)
}

func SugerInfof(template string, args ...interface{}) {
	logger.Sugar().Infof(template, args)
}

func SugerInfo(args ...interface{}) {
	logger.Sugar().Info(args)
}
