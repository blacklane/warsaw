package contexts

import (
	"context"

	"github.com/google/uuid"
)

type requestIDContext int

const requestIDKey requestIDContext = 0

// GetRequestID allows to access the RequestID details from the context.
func GetRequestID(ctx context.Context) string {
	if requestID, ok := ctx.Value(requestIDKey).(string); ok {
		return requestID
	}
	return uuid.New().String()
}

func WithRequestID(ctx context.Context, requestID string) context.Context {
	return context.WithValue(ctx, requestIDKey, requestID)
}
