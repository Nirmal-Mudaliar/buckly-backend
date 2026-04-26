package utils

import (
	"context"
	"log"

	"go.uber.org/zap"
)

var logger *zap.Logger

func InitLog(isDevelopment bool) *zap.Logger {
	var err error
	logger, err = GetLogger(isDevelopment)
	if err != nil {
		log.Fatalf("Error occured while initializing logger: %v", err)
	}
	return logger
}

func GetLogger(isDevelopment bool) (*zap.Logger, error) {
	var logger *zap.Logger
	var err error
	if isDevelopment {
		logger, err = zap.NewDevelopment()
	} else {
		logger, err = zap.NewProduction()
	}
	if err != nil {
		log.Fatalf("Error occured while initializing logger: %v", err)
		return nil, err
	}
	return logger, nil
}

func GetLoggerFromContext(ctx context.Context) *zap.Logger {
	return logger.With(zap.Any("context", ctx))
}

func SetLoggerInContext(ctx context.Context, logger *zap.Logger) context.Context {
	return context.WithValue(ctx, "logger", logger)
}
