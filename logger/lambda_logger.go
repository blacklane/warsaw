package logger

import (
	"context"

	"github.com/aws/aws-lambda-go/lambdacontext"

	"github.com/blacklane/warsaw/constants"
)

// NewLambdaLogger returns a logger and enhanced context which is ready to log details of request in JSON responses
// compatible with Kiev format.
func NewLambdaLogger(ctx context.Context) (Logger, context.Context) {
	log, loggingContext := New(ctx, lambdacontext.FunctionName)
	setupFromLambdaContext(ctx, log)
	return log, loggingContext
}

func setupFromLambdaContext(ctx context.Context, log Logger) {
	lc, _ := lambdacontext.FromContext(ctx)

	log.WithScope(map[string]interface{}{
		constants.FieldRequestID:    lc.AwsRequestID,
		constants.FieldEntryPoint:   true,
		"lambda_function_arn":       lc.InvokedFunctionArn,
		"lambda_function_version":   lambdacontext.FunctionVersion,
		"lambda_memory_limit_in_mb": lambdacontext.MemoryLimitInMB,
	})
}
