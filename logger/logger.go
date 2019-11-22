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

func Get(ctx context.Context) Logger {
	return logger{internalLoggerFromContext(ctx)}
}