package web_pusher

import "log"

type Logger interface {
	Printf(format string, v ...interface{})
}

type StandardLogger struct {
}

func (_ StandardLogger) Printf(format string, v ...interface{}) {
	log.Printf(format, v...)
}

type NopeLogger struct {
}

func (_ NopeLogger) Printf(format string, v ...interface{}) {
}
