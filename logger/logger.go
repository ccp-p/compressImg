package logger

import (
	"log"
	"os"
)

var (
	Logger *log.Logger
)

func InitLogger(logFile *os.File) *log.Logger {
	Logger = log.New(logFile, "prefix", log.LstdFlags)
	return Logger
}
