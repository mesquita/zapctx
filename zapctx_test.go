package zapctx_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"

	"github.com/mesquita/zapctx"
)

func TestFrom(t *testing.T) {
	t.Run("should return a non-nil zap logger given there isn't a logger in the context", func(t *testing.T) {
		assert.NotNil(t, zapctx.From(context.Background()))
	})

	t.Run("should return the zap logger we put into it", func(t *testing.T) {
		var (
			logger = zap.NewNop()
			ctx    = zapctx.With(context.Background(), logger)
		)
		assert.Equal(t, logger, zapctx.From(ctx))
	})
}

func TestWith(t *testing.T) {
	t.Run("should return the zap logger we put into it", func(t *testing.T) {
		var (
			logger = zap.NewNop()
			ctx    = zapctx.With(context.Background(), logger)
		)
		assert.Equal(t, logger, zapctx.From(ctx))
	})
}

func TestWithFields(t *testing.T) {
	t.Run("should return the original context if there's nothing to do", func(t *testing.T) {
		ctx := context.Background()
		assert.Equal(t, ctx, zapctx.WithFields(ctx))
	})

	t.Run("should add fields to the logger in the context", func(t *testing.T) {
		const (
			key   = "key"
			value = "value"
		)

		var (
			logger = zap.NewNop()
			ctx    = zapctx.WithFields(zapctx.With(context.Background(), logger), zap.String(key, value))
		)

		assert.Equal(t, logger.With(zap.String(key, value)), zapctx.From(ctx))
	})
}

func TestWithTraceID(t *testing.T) {
	t.Run("should return the context with the trace id field", func(t *testing.T) {
		const (
			aTraceID = "a-trace-id"
		)
		var (
			logger = zap.NewNop()
			ctx    = zapctx.WithTraceID(zapctx.With(context.Background(), logger), aTraceID)
		)

		assert.Equal(t, logger.With(zap.String(zapctx.TraceID, aTraceID)), zapctx.From(ctx))
	})
}

func TestSetLogLevel(t *testing.T) {
	t.Run("should change log level in runtime", func(t *testing.T) {
		var (
			ctx    = context.Background()
			logger = zapctx.From(ctx)
		)

		assert.False(t, logger.Core().Enabled(zapcore.DebugLevel))
		assert.True(t, logger.Core().Enabled(zapcore.InfoLevel))

		zapctx.SetLogLevel(zapcore.DebugLevel)

		assert.True(t, logger.Core().Enabled(zapcore.DebugLevel))
		assert.True(t, logger.Core().Enabled(zapcore.InfoLevel))
	})
}

func TestAddHooks(t *testing.T) {
	t.Run("should add and execute hook", func(t *testing.T) {
		ctx := context.Background()

		i, e, w := 0, 0, 0
		zapctx.AddHooks(func(entry zapcore.Entry) error {
			switch entry.Level {
			case zapcore.ErrorLevel:
				e++
			case zapcore.InfoLevel:
				i++
			case zapcore.WarnLevel:
				w++
			default:
				assert.Fail(t, fmt.Sprintf("wrong level found: %s", entry.Level))
			}

			return nil
		})

		zapctx.Error(ctx, "error")
		zapctx.Info(ctx, "info")
		zapctx.Warn(ctx, "warning")

		// assert that the hook was called 3 times with the correct level
		assert.Equal(t, 1, i)
		assert.Equal(t, 1, e)
		assert.Equal(t, 1, w)
	})
}
