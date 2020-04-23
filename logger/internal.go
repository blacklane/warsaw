package logger

import (
	"context"
	"io"

	"github.com/rs/zerolog"
)

// internalLogger shares interface of the logger and encapsulates the internal zerolog used for it's implementation
type internalLogger interface {
	Info() *LoggedEvent
	Err(err error) *LoggedEvent
	WithContext(context.Context) context.Context
	UpdateContext(func(Context) Context)
	GetLevel() Level
}

// LoggedEvent is a type-alias for zerolog.Event
type LoggedEvent = zerolog.Event

// Context is a type-alias for zerolog.Context
type Context = zerolog.Context

// Level is a type-alias for zerolog.Level
type Level = zerolog.Level

const disabledLoggerLevel = zerolog.Disabled

func newInternalLogger(w io.Writer) internalLogger {
	log := zerolog.New(w)
	return &log
}

func internalLoggerFromContext(ctx context.Context) internalLogger {
	return zerolog.Ctx(ctx)
}
