package web

import (
	"context"
	"net/http"

	"github.com/pborman/uuid"
)

type ctxKey struct{}

const (
	RequestIDHeader = "X-Request-Id"
)

func GetRequestIDFromRequestHeader(r *http.Request) string {
	requestID := r.Header.Get(RequestIDHeader)
	if requestID != "" {
		return requestID
	}
	return uuid.New()
}

func GetRequestIDFromContext(ctx context.Context) string {
	reqID := ctx.Value(ctxKey{})
	if reqID != nil {
		return reqID.(string)
	}
	return ""
}

func SetRequestID(ctx context.Context, requestID string) context.Context {
	return context.WithValue(ctx, ctxKey{}, requestID)
}
