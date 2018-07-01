package stdout_logger

import "log"

type StandardLogger struct {
}

func (_ StandardLogger) Printf(format string, v ...interface{}) {
	log.Printf(format, v...)
}
