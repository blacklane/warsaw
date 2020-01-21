package logger

import (
	"bytes"
	"net/http"
)

type loggedResponseWriter struct {
	http.ResponseWriter
	statusCode int
	body       *bytes.Buffer
}

func (w loggedResponseWriter) StatusCode() int {
	return w.statusCode
}

func (w loggedResponseWriter) Body() string {
	w.body.Reset()
	return w.body.String()
}

func (w *loggedResponseWriter) WriteHeader(statusCode int) {
	w.ResponseWriter.WriteHeader(statusCode)
	w.statusCode = statusCode
}

func (w *loggedResponseWriter) Write(buf []byte) (int, error) {
	if w.statusCode < http.StatusOK || w.statusCode > 299 {
		w.body.Write(buf)
	}
	return w.ResponseWriter.Write(buf)
}

func newLoggedResponseWriter(w http.ResponseWriter) *loggedResponseWriter {
	return &loggedResponseWriter{ResponseWriter: w, body: new(bytes.Buffer), statusCode: http.StatusOK}
}
