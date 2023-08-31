package controllers

import (
	"go.uber.org/zap"
	"log"
)

func NewLogger() *zap.Logger {
	logger, err := zap.NewDevelopment()
	zap.AddStacktrace(logger.Level())
	if err != nil {
		log.Fatalf("Unable to construct logger: %v", err)
	}
	return logger
}
