package logger

import (
	"context"
)

type logger struct {
	log internalLogger
}

type Logger interface {
	Event(name string) *Event
}

// Get returns Logger instance from the argument representing current `context.Context`. Useful to get the logger
// instance downstream somewhere deep in your app. Then you pass just the context instance all way down and get
// the `Logger` out using this function.
func Get(ctx context.Context) Logger {
	return logger{internalLoggerFromContext(ctx)}
}
