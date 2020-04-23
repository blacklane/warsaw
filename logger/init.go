package logger

import (
	"time"

	"github.com/blacklane/warsaw/constants"

	"github.com/rs/zerolog"
)

func init() {
	now = time.Now

	zerolog.TimestampFieldName = constants.FieldTimestamp
	zerolog.TimeFieldFormat = time.RFC3339Nano
	zerolog.TimestampFunc = now

	DefaultLogger = buildDefaultLogger(LogSink)
}
