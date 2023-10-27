package utils

import (
	"log"
	"os"

	"github.com/newrelic/go-agent/v3/integrations/logcontext-v2/logWriter"
	"gorm.io/gorm/logger"
)

const (
	reset  = "\033[0m"
	red    = "\033[31m"
	yellow = "\033[33m"
	green  = logger.Green
)

type Logger struct {
	l *log.Logger
}

func NewLogger() *Logger {
	newRelicApp := NewRelicInstance()
	writer := logWriter.New(os.Stdout, newRelicApp.App)
	l := log.New(&writer, "", log.Default().Flags())

	return &Logger{l}

}

func (l Logger) Error(v string, exit bool) {

	l.l.Printf("%sERROR: %s%s\n", red, v, reset)

	if exit {
		os.Exit(1)
	}

}

func (l Logger) Warn(v string) {
	l.l.Printf("%sWARN: %s%s\n", yellow, v, reset)
}

func (l Logger) Info(v string) {

	l.l.Printf("%sINFO: %s%s\n", green, v, reset)
}
