package middlewares

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"time"

	"github.com/blacklane/warsaw/logger"
	"github.com/blacklane/warsaw/request_context/contexts"
)

func ExampleNewHttpHandlerLogger_simple() {
	loggerMiddleware := NewHttpHandlerLogger()

	h := loggerMiddleware(
		http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) { _, _ = fmt.Fprint(w, "Hello, world") }))
	http.Handle("example", h)
}

func ExampleNewHttpHandlerLogger_complete() {
	sec := -1
	// Set current time function so we can control the request duration
	logger.SetNowFunc(func() time.Time {
		sec++
		return time.Date(2009, time.November, 10, 23, 0, sec, 0, time.UTC)
	})

	r := httptest.NewRequest(http.MethodGet, "http://example.com/foo", nil)
	w := httptest.NewRecorder()

	ctx := contexts.WithRequestID(r.Context(), "42")

	// It's needed otherwise 'go test' will read empty string from stdout when verifying the output
	logger.LogSink = os.Stdout
	_, ctx = logger.New(ctx, "")

	rr := r.WithContext(ctx)
	loggerMiddleware := NewHttpHandlerLogger()

	h := loggerMiddleware(http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) { return }))

	h.ServeHTTP(w, rr)

	// Output:
	// {"level":"info","application":"","entry_point":true,"host":"example.com","ip":"192.0.2.1","params":"","path":"/foo","request_depth":0,"request_id":"42","route":"","tree_path":"","user_agent":"","verb":"GET","request_duration":1000,"status":0,"timestamp":"2009-11-10T23:00:02Z","event":"request_finished","message":"GET /foo"}
}
