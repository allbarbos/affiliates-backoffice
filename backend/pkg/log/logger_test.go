package log

import (
	"context"
	"testing"

	ctxpkg "affiliates-backoffice-backend/pkg/context"

	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

func Test_New(t *testing.T) {
	assert.NotNil(t, New())
}

func Test_NewWithZap(t *testing.T) {
	zl, _ := zap.NewProduction()
	l := NewWithZap(zl)
	assert.NotNil(t, l)
}

func Test_LoggerWith(t *testing.T) {
	ctx := context.Background()
	ctxpkg.SetValue(&ctx, ctxpkg.KeyRequestID, "123")
	ctxpkg.SetValue(&ctx, ctxpkg.KeyCorrelationID, "321")
	loggerTest, observer := NewForTest()

	loggerTest.With(ctx).Info("mock message")

	ctxMap := observer.All()[0].ContextMap()
	assert.Equal(t, observer.All()[0].Entry.Message, "mock message")
	assert.Equal(t, ctxMap["requestID"], "123")
	assert.Equal(t, ctxMap["correlationID"], "321")
}

func Test_LoggerWith_NoArgs(t *testing.T) {
	ctx := context.Background()
	loggerTest, observer := NewForTest()

	loggerTest.With(ctx).Info("mock message")

	assert.Equal(t, observer.All()[0].Entry.Message, "mock message")
}

func Test_NewForTest(t *testing.T) {
	logger, entries := NewForTest()
	assert.Equal(t, 0, entries.Len())
	logger.Info("msg 1")
	assert.Equal(t, 1, entries.Len())
	logger.Info("msg 2")
	logger.Info("msg 3")
	assert.Equal(t, 3, entries.Len())
	entries.TakeAll()
	assert.Equal(t, 0, entries.Len())
	logger.Info("msg 4")
	assert.Equal(t, 1, entries.Len())
}
