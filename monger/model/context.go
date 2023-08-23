package model

import (
	"context"

	"github.com/google/uuid"
)

type ContextKey string // can be unexported

// ContextKeyRequestID is the ContextKey for RequestID
const ContextKeyRequestID ContextKey = "requestID" // can be unexported

// AttachRequestID will attach a brand new request ID to a http request
func AssignRequestID(ctx context.Context) context.Context {

	reqID := uuid.New()

	return context.WithValue(ctx, ContextKeyRequestID, reqID.String())
}

// GetRequestID will get reqID from a http request and return it as a string
func GetRequestID(ctx context.Context) string {

	reqID := ctx.Value(ContextKeyRequestID)

	if ret, ok := reqID.(string); ok {
		return ret
	}

	return ""
}
