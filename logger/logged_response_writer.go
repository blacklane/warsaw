package logger

import (
	"bytes"
	"net/http"
)

type LoggedResponseWriter struct {
	http.ResponseWriter
	statusCode int
	body       *bytes.Buffer
}

func (w LoggedResponseWriter) StatusCode() int {
	return w.statusCode
}

func (w LoggedResponseWriter) Body() string {
	w.body.Reset()
	return w.body.String()
}

func (w *LoggedResponseWriter) WriteHeader(statusCode int) {
	w.ResponseWriter.WriteHeader(statusCode)
	w.statusCode = statusCode
}

func (w *LoggedResponseWriter) Write(buf []byte) (int, error) {
	if w.statusCode < http.StatusOK || w.statusCode > 299 {
		w.body.Write(buf)
	}
	return w.ResponseWriter.Write(buf)
}

func NewLoggedResponseWriter(w http.ResponseWriter) *LoggedResponseWriter {
	return &LoggedResponseWriter{ResponseWriter: w, body: new(bytes.Buffer), statusCode: http.StatusOK}
}
