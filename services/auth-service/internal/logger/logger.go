package logger

import "go.uber.org/zap"

var Log *zap.SugaredLogger

func InitLogger() {
	logger, _ := zap.NewProduction()
	defer logger.Sync()
	Log = logger.Sugar()
}
