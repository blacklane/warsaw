package main

import (
	"context"
	"errors"

	"github.com/blacklane/warsaw/logger"
)

func defaultLogger() {
	logger.Event("mySimpleEvent").
		Msg("log my line")

	logger.Event("myComplexEvent").
		Str("aString", "field").
		Int("num", 422).
		Bool("valid", false).
		Err(errors.New("something failed")).
		Msg("foo failed")

	logger.Error("An event", errors.New("error")).
		Msg("Failed to create something")
}

func standaloneLogger() {
	log, ctx := logger.New(context.TODO(), "myAppName")
	log.Event("An event").Msg("hello world")
	log.Error("An event", errors.New("error")).Msg("Failed to create something")

	useLoggerFromContext(ctx)
}

func useLoggerFromContext(ctx context.Context) {
	log := logger.Get(ctx)
	log.Event("useLoggerFromContext").Bool("aTruth", true).Msg("hello :)")
}
