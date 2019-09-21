package logger

import (
	"go.uber.org/zap"
)

type Config struct {
	ErrorOutputs []string
	Outputs      []string
	Debug        bool
}

func InitLogger(params Config) *zap.Logger {
	config := zap.NewProductionConfig()

	config.OutputPaths = params.Outputs
	config.ErrorOutputPaths = params.ErrorOutputs
	if params.Debug {
		config.Level.SetLevel(zap.DebugLevel)
	} else {
		config.Level.SetLevel(zap.InfoLevel)
	}

	zapLogger, err := config.Build()
	if err != nil {
		panic(err)
	}

	return zapLogger
}
