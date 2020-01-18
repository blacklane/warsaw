package logger

import (
	"io"
	"os"
)

// LogSink points to the writer that store all logs, default one is os.Stdout
// If you want to have it different, just set it prior using any `logger.New*(...)` function
var LogSink io.Writer

func init() {
	LogSink = os.Stdout
}
