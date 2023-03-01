package context

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_Context(t *testing.T) {
	ctx := context.Background()
	SetValue(&ctx, KeyRequestID, "123")
	correlationID := GetValue(ctx, KeyCorrelationID)
	requestID := GetValue(ctx, KeyRequestID)
	assert.Equal(t, "123", requestID)
	assert.Equal(t, "", correlationID)
	assert.NotNil(t, ctx)
}
