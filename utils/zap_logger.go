package utils

import (
	"go.uber.org/zap"
)

type ZapLogger struct {
	zap *zap.Logger
}

func (z ZapLogger) Info(v string) {
	z.zap.Info(v)
}

func (z ZapLogger) Error(v string) {
	z.zap.Error(v)
}

func (z ZapLogger) Warn(v string) {
	z.zap.Warn(v)
}

func (z ZapLogger) Fatal(v string) {
	z.zap.Fatal(v)
}

func NewZapLogger(config Config) *ZapLogger {
	env := config.Environment()
	if env == "development" {
		z, _ := zap.NewDevelopment()
		return &ZapLogger{zap: z}
	} else if env == "production" {
		z, _ := zap.NewProduction()
		return &ZapLogger{zap: z}
	} else {
		z, _ := zap.NewProduction()
		return &ZapLogger{zap: z}

	}

}
