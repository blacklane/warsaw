package main

import (
	"net/http"

	"github.com/go-chi/chi"

	"github.com/blacklane/warsaw/logger"
)

func newLoggerMiddleware(appName string) func(handler http.Handler) http.Handler {
	logger := logger.NewKievRequestLogger(appName)

	// wrap logger to match middleware func signature
	fn := func(next http.Handler) http.Handler {
		return logger(next.(http.HandlerFunc))
	}
	return fn
}

func addRouterMiddleware() {
	router := chi.NewRouter()
	router.Use(newLoggerMiddleware("MyApp"))
}
