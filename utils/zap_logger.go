package utils

import (
	"go.uber.org/zap"
)

type ZapLogger struct {
	zap *zap.Logger
}

func (z ZapLogger) Info(v string) {
	defer func(zap *zap.Logger) {
		_ = zap.Sync()
	}(z.zap)
	z.zap.Sugar().Info(v)
}

func (z ZapLogger) Error(v string) {
	defer func(zap *zap.Logger) {
		_ = zap.Sync()
	}(z.zap)
	z.zap.Sugar().Error(v)
}

func (z ZapLogger) Warn(v string) {
	defer func(zap *zap.Logger) {
		_ = zap.Sync()
	}(z.zap)
	z.zap.Sugar().Warn(v)
}

func (z ZapLogger) Fatal(v string) {
	defer func(zap *zap.Logger) {
		_ = zap.Sync()
	}(z.zap)
	z.zap.Sugar().Fatal(v)
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
