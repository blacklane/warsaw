package logger

import (
	"context"

	"github.com/aws/aws-lambda-go/lambdacontext"

	"github.com/blacklane/warsaw/logger/kiev_fields"
)

func NewLambdaLogger(ctx context.Context) (Logger, context.Context) {
	log, loggingContext := New(ctx, lambdacontext.FunctionName)
	setupFromLambdaContext(log, ctx)
	return log, loggingContext
}

func setupFromLambdaContext(log Logger, ctx context.Context) {
	lc, _ := lambdacontext.FromContext(ctx)

	log.WithScope(map[string]interface{}{
		kiev_fields.RequestID:       lc.AwsRequestID,
		kiev_fields.EntryPoint:      true,
		"lambda_function_arn":       lc.InvokedFunctionArn,
		"lambda_function_version":   lambdacontext.FunctionVersion,
		"lambda_memory_limit_in_mb": lambdacontext.MemoryLimitInMB,
	})
}
