package main

import (
	"fmt"
	"net/http"

	"github.com/blacklane/warsaw/logger"
)

func pingHandler(w http.ResponseWriter, req *http.Request) {
	log := logger.Get(req.Context())

	log.Event("ping_started").Str("some_field", "value").Int("some_int", 123).Send()

	_, _ = fmt.Fprint(w, "ping")
}

func addLogMiddleware() {
	loggerMiddleware := logger.NewKievRequestLogger("MyAppName")

	http.HandleFunc("/ping", loggerMiddleware(pingHandler))
	if err := http.ListenAndServe(":8080", nil); err != nil {
		panic(err)
	}
}
