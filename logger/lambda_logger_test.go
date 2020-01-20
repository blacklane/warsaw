package logger

import (
	"context"
	"testing"

	"github.com/aws/aws-lambda-go/lambdacontext"
)

func TestNewLambdaLogger(t *testing.T) {

	t.Run("basic logger context is build", func(t *testing.T) {
		lambdacontext.FunctionName = "myFunctionName"
		lc := &lambdacontext.LambdaContext{
			AwsRequestID:       "a-request-id",
			InvokedFunctionArn: "lambda:arn",
		}
		ctx := lambdacontext.NewContext(context.Background(), lc)

		out := captureLogs(func() {
			log, _ := NewLambdaLogger(ctx)
			log.Event("aName").Send()
		})

		// lambda_function_name and lambda_function_version are now blank but in lambda environment they will be reported correctly, same as lambda_memory_limit_in_mb
		matchEvent(t, out, map[string]string{"application": "myFunctionName", "entry_point": "true", "request_id": "a-request-id", "event": "aName", "lambda_memory_limit_in_mb": "0", "lambda_function_arn": "lambda:arn", "lambda_function_version": ""})
	})
}
