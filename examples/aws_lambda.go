package main

import (
	"context"
	"fmt"

	"github.com/blacklane/warsaw/logger"
)

type MyEvent struct {
	Name string `json:"name"`
}

func foo(ctx context.Context) {
	log := logger.Get(ctx)
	log.Event("I'm").Msg("insideMethod")
}

func handleRequest(ctx context.Context, name MyEvent) (string, error) {
	log, ctx := logger.NewLambdaLogger(ctx)
	log.Event("Called").Str("the_name", name.Name).Send()

	foo(ctx)
	return fmt.Sprintf("Hello %s!", name.Name), nil
}
