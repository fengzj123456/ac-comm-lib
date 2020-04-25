package zaplog

import (
	"go.uber.org/zap"
)

var (
	Config zap.Config
	Std    *Logger
)

func init() {
	Config = zap.NewDevelopmentConfig()
	if err := Reset(); err != nil {
		panic(err)
	}
}

func Reset() error {
	logger, err := Config.Build(zap.AddStacktrace(zap.DPanicLevel))
	if err != nil {
		return err
	}
	Std = NewLogger(logger)
	return nil
}
