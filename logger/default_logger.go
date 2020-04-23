package logger

import (
	"io"
)

// DefaultLogger represents everywhere available logger with default LogSink.
// It is also used with package logger.Event(...) method executions.
var DefaultLogger defaultLogger

type defaultLogger interface {
	eventLogger
}

func buildDefaultLogger(sink io.Writer) defaultLogger {
	return logger{newInternalLogger(sink)}
}
