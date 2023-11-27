package utils

import (
	"go.uber.org/zap"
)

type zapLogger struct {
	zap *zap.Logger
}

func (z zapLogger) Info(v string) {
	z.zap.Info(v)
}

func (z zapLogger) Error(v string) {
	z.zap.Error(v)
}

func (z zapLogger) Warn(v string) {
	z.zap.Warn(v)
}

func (z zapLogger) Fatal(v string) {
	z.zap.Fatal(v)
}

func NewZapLogger(config Config) Logger {
	env := config.Environment()
	if env == "development" {
		z, _ := zap.NewDevelopment()
		return zapLogger{zap: z}
	} else if env == "production" {
		z, _ := zap.NewProduction()
		return zapLogger{zap: z}
	} else {
		z, _ := zap.NewProduction()
		return zapLogger{zap: z}

	}

}
