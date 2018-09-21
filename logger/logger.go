package logger

const (
	VERBOSE_LOGGING = 1
)

type Logger interface {
	Printf(format string, v ...interface{})
}
