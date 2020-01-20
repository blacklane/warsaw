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

// Type aliases to abstract internal zerolog
type LoggedEvent = zerolog.Event
type Context = zerolog.Context
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
