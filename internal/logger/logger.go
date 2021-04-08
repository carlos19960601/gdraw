package logger

import (
	"go.uber.org/zap"
)

var logger *zap.Logger

func Init() (err error) {
	logger, err = zap.NewProduction(zap.WithCaller(true))
	if err != nil {
		return err
	}

	zap.ReplaceGlobals(logger)
	return nil
}
