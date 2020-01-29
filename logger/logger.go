package logger

import (
	"go.uber.org/zap"
)

//Init should be called from main it instantiate a zap sugared logger
func Init(isDebug bool) *zap.SugaredLogger {
	zapLogger, err := zap.NewProduction()

	if isDebug {
		zapLogger, err = zap.NewDevelopment()
	}

	if err != nil {
		panic(err)
	}
	return zapLogger.Sugar()
}
