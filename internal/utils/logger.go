package utils

import (
	"go.uber.org/zap"
	"log"
)

var Log *zap.SugaredLogger

func InitLogger() error {
	logger, err := zap.NewProduction()
	if err != nil {
		return err
	}

	Log = logger.Sugar()

	return nil
}

func DeferLoggerClose() {
	err := Log.Sync()
	if err != nil {
		log.Panicln("could not flush buffer")
	}
}
