package apiserver

import (
	"io"
	"log"
)

type CustomLogger struct {
	Info    *log.Logger
	Warning *log.Logger
	Error   *log.Logger
}

func newLogger(infoHandle, warningHandle, errorHandle io.Writer) *CustomLogger {
	return &CustomLogger{
		Info:    log.New(infoHandle, "INFO: ", log.Ldate|log.Ltime|log.Lshortfile),
		Warning: log.New(warningHandle, "WARNING: ", log.Ldate|log.Ltime|log.Lshortfile),
		Error:   log.New(errorHandle, "ERROR: ", log.Ldate|log.Ltime|log.Lshortfile),
	}
}
