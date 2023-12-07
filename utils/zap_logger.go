package utils

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
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
	var logger *zap.Logger
	if env == "development" {
		c := zap.NewDevelopmentConfig()
		c.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
		logger, _ := c.Build()
		return &ZapLogger{zap: logger}
	}

	logger, _ = zap.NewProduction()
	return &ZapLogger{zap: logger}

}
