package context

import (
	"context"
)

type key int

const (
	KeyCorrelationID key = iota
	KeyRequestID     key = iota
)

func SetValue(ctx *context.Context, key key, value string) {
	*ctx = context.WithValue(*ctx, key, value)
}

func GetValue(ctx context.Context, key key) string {
	if v := ctx.Value(key); v != nil {
		return v.(string)
	}
	return ""
}
