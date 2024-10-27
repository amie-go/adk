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
	Verbose bool
	Values  []string
}

// --------------------------------------------

func WithVerbose(value bool) options.WithFn[config] {
	return func(o *config) { o.Verbose = value }
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
	t.Run("Test nil config", func(t *testing.T) {
		var nilConfig *config
		assert.NotPanics(t, func() { options.Apply(nilConfig, nil) })
		assert.NotPanics(t, func() { options.Apply(nilConfig, options.WithFn[config](nil)) })
		assert.NotPanics(t, func() { options.Apply(nilConfig, WithVerbose(true)) })
	})

	t.Run("Test config", func(t *testing.T) {
		var cfg config
		assert.NotPanics(t, func() { options.Apply(&cfg, nil) })
		assert.NotPanics(t, func() { options.Apply(&cfg, options.WithFn[config](nil)) })
		options.Apply(&cfg, WithVerbose(true))
		assert.True(t, cfg.Verbose)
		assert.NotPanics(t, func() { options.Apply(&cfg, nil, WithVerbose(false), nil) })
		assert.False(t, cfg.Verbose)
	})

	t.Run("Test config with struct opt", func(t *testing.T) {
		var cfg config
		options.Apply(&cfg, WithStructOpt("a", "b"))
		assert.Equal(t, []string{"a", "b"}, cfg.Values)
		options.Apply(&cfg, WithStructPtrOpt("c", "d"))
		assert.Equal(t, []string{"a", "b", "c", "d"}, cfg.Values)
	})
}

func TestNew(t *testing.T) {
	t.Run("Test new any", func(t *testing.T) {
		var result = options.New[any]()
		assert.NotNil(t, result)
	})

	t.Run("Test new config", func(t *testing.T) {
		var result = options.New(WithVerbose(true))
		assert.NotNil(t, result)
		assert.EqualValues(t, true, result.Verbose)
	})
}
