package options_test

import (
	"context"
	"testing"

	"github.com/amie-go/adk/options"
	"github.com/stretchr/testify/assert"
)

// --------------------------------------------
// Test materials

type config struct {
	Verbose     bool
	Values      []string
	FromContext string
}

// --------------------------------------------

func WithVerbose(value bool) options.WithFn[config] {
	return func(o *config) { o.Verbose = value }
}

// --------------------------------------------

type dataCtxKeyType string

var dataCtxKey = dataCtxKeyType("data")

func WithDataFromContext() options.WithCtxFn[config] {
	return func(ctx context.Context, o *config) { o.FromContext = ctx.Value(dataCtxKey).(string) }
}

// --------------------------------------------

func WithStructOpt(someValues ...string) options.With[config] {
	return structOpt(someValues)
}

type structOpt []string

func (obj structOpt) Apply(ctx context.Context, c *config) { c.Values = append(c.Values, obj...) }

// --------------------------------------------

func WithStructPtrOpt(someValues ...string) options.With[config] {
	return &structPtrOpt{values: someValues}
}

type structPtrOpt struct {
	values []string
}

func (obj *structPtrOpt) Apply(ctx context.Context, c *config) {
	c.Values = append(c.Values, obj.values...)
}

// --------------------------------------------

func TestApply(t *testing.T) {
	var ctx = context.Background()

	t.Run("Test nil config", func(t *testing.T) {
		var nilConfig *config
		assert.NotPanics(t, func() { options.Apply(ctx, nilConfig, nil) })
		assert.NotPanics(t, func() { options.Apply(ctx, nilConfig, options.WithFn[config](nil)) })
		assert.NotPanics(t, func() { options.Apply(ctx, nilConfig, WithVerbose(true)) })
	})

	t.Run("Test nil options", func(t *testing.T) {
		var cfg config
		assert.NotPanics(t, func() { options.Apply(ctx, &cfg, nil) })
		assert.NotPanics(t, func() { options.Apply(ctx, &cfg, options.WithFn[config](nil)) })
		assert.NotPanics(t, func() { options.Apply(ctx, &cfg, nil, WithVerbose(true), nil) })
		assert.True(t, cfg.Verbose)
	})

	t.Run("Test config", func(t *testing.T) {
		var cfg config
		options.Apply(ctx, &cfg, WithVerbose(true))
		assert.True(t, cfg.Verbose)
	})

	t.Run("Test config with context", func(t *testing.T) {
		var expected = "some value expected"
		ctx = context.WithValue(ctx, dataCtxKey, expected)
		var cfg config
		assert.NotPanics(t, func() { options.Apply(ctx, &cfg, WithDataFromContext()) })
		assert.Equal(t, expected, cfg.FromContext)
	})

	t.Run("Test config with struct opt", func(t *testing.T) {
		var cfg config
		options.Apply(ctx, &cfg, WithStructOpt("a", "b"))
		assert.Equal(t, []string{"a", "b"}, cfg.Values)
		options.Apply(ctx, &cfg, WithStructPtrOpt("c", "d"))
		assert.Equal(t, []string{"a", "b", "c", "d"}, cfg.Values)
	})
}

func TestNew(t *testing.T) {
	var ctx = context.Background()

	t.Run("Test new any", func(t *testing.T) {
		var result = options.New[any](ctx)
		assert.NotNil(t, result)
	})

	t.Run("Test new config", func(t *testing.T) {
		var result = options.New(ctx, WithVerbose(true))
		assert.NotNil(t, result)
		assert.EqualValues(t, true, result.Verbose)
	})
}
