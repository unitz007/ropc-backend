package logger

import (
	"log"
	"os"

	"gorm.io/gorm/logger"
)

const (
	reset  = "\033[0m"
	red    = "\033[31m"
	yellow = "\033[33m"
	green  = logger.Green
)

func Error(v string, exit bool) {
	log.Printf("%sERROR: %s%s\n", red, v, reset)

	if exit {
		os.Exit(1)
	}

}

func Warn(v string) {
	log.Printf("%sWARN: %s%s\n", yellow, v, reset)
}

func Info(v string) {
	log.Printf("%sINFO: %s%s\n", green, v, reset)
}
