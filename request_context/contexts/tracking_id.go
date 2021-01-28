package contexts

import (
	"context"

	"github.com/google/uuid"
)

type trackingIDContext int

const trackingIDKey trackingIDContext = 0

// GetTrackingID allows to access the TrackingID details from the context.
func GetTrackingID(ctx context.Context) string {
	if trackingID, ok := ctx.Value(trackingIDKey).(string); ok {
		return trackingID
	}
	return uuid.New().String()
}

// GetRequestID allows to access the TrackingID details from the context. (deprecated)
func GetRequestID(ctx context.Context) string {
	return GetTrackingID(ctx)
}

func WithTrackingID(ctx context.Context, trackingID string) context.Context {
	return context.WithValue(ctx, trackingIDKey, trackingID)
}

func WithRequestID(ctx context.Context, requestID string) context.Context {
	return WithTrackingID(ctx, requestID)
}
