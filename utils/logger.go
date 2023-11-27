package utils

import (
	"log"
	"os"

	"github.com/newrelic/go-agent/v3/integrations/logcontext-v2/logWriter"
)

const (
	reset  = "\033[0m"
	red    = "\033[31m"
	yellow = "\033[33m"
	green  = "\033[32m"
)

type Logger interface {
	Info(v string)
	Error(v string)
	Warn(v string)
	Fatal(v string)
}

type logger struct {
	l *log.Logger
}

func (l logger) Fatal(v string) {
	l.l.Fatal(v)
}

func NewLogger() Logger {
	newRelicApp := NewRelicInstance()
	writer := logWriter.New(os.Stdout, newRelicApp.App)
	l := log.New(&writer, "", log.Default().Flags())

	return &logger{l}

}

func (l logger) Error(v string) {

	l.l.Printf("%sERROR: %s%s\n", red, v, reset)
}

func (l logger) Warn(v string) {
	l.l.Printf("%sWARN: %s%s\n", yellow, v, reset)
}

func (l logger) Info(v string) {

	l.l.Printf("%sINFO: %s%s\n", green, v, reset)
}
