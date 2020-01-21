package logger

import (
	"context"
	"io"
	"time"

	"github.com/rs/zerolog"

	"github.com/blacklane/warsaw/logger/kiev_fields"
)

// internalLogger shares interface of the logger and encapsulates the internal zerolog used for it's implementation
type internalLogger interface {
	Info() *LoggedEvent
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
	zerolog.TimestampFieldName = kiev_fields.Timestamp
	zerolog.TimeFieldFormat = time.RFC3339Nano

	log := zerolog.New(w)
	return &log
}

func internalLoggerFromContext(ctx context.Context) internalLogger {
	return zerolog.Ctx(ctx)
}
