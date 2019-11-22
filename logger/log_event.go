package logger

import (
	"github.com/blacklane/warsaw/logger/kiev_fields"
)

func (logger logger) Event(name string) *Event {
	return logger.log.Info().Timestamp().Str(kiev_fields.Event, name)
}
