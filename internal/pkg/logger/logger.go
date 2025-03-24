package logger

import (
	"go.uber.org/zap"
)

var log *zap.Logger

func Init() {
	config := zap.NewProductionConfig()

	var err error
	log, err = config.Build()
	if err != nil {
		panic(err)
	}
}

func GetLogger() *zap.Logger {
	return log
}
