package utils

import "go.uber.org/zap"

func GetLogger() *zap.Logger {
	var logger *zap.Logger
	cfg := zap.NewProductionConfig()
	cfg.Level = zap.NewAtomicLevelAt(zap.DebugLevel)
	var err error
	logger, err = cfg.Build()
	if err != nil {
		panic("failed to initialize logger: " + err.Error())
	}
	return logger
}
