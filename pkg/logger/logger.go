package logger

import (
	"log"
	"os"
)

var (
	l *loggers
)

func Init() {
	flags := log.Ldate | log.Ltime | log.Lshortfile

	logInfo := log.New(os.Stdout, "INFO: ", flags)
	logWarn := log.New(os.Stdout, "WARN: ", flags)
	logErr := log.New(os.Stderr, "ERROR: ", flags)

	l = &loggers{
		logInfo: logInfo,
		logWarn: logWarn,
		logErr:  logErr,
	}
}

type loggers struct {
	logInfo *log.Logger
	logWarn *log.Logger
	logErr  *log.Logger
}

func Info(v ...interface{}) {
	l.logInfo.Println(v...)
}

func Warn(v ...interface{}) {
	l.logWarn.Println(v...)
}

func Err(v ...interface{}) {
	l.logErr.Println(v...)
}
