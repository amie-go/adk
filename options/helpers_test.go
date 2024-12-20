package options_test

import (
	"context"
	"testing"

	"github.com/amie-go/adk/options"
	"github.com/stretchr/testify/assert"
)

// --------------------------------------------

type configNoDefault struct {
	Values []string
}

type configWithDefault struct {
	Values []string
}

func (c *configWithDefault) SetDefaults() {
	c.Values = []string{"default"}
}

// --------------------------------------------

func TestSetDefaults(t *testing.T) {
	var ctx = context.Background()

	t.Run("Test nil config", func(t *testing.T) {
		var nilConfig *configNoDefault
		assert.NotPanics(t, func() { options.Apply(ctx, nilConfig, options.SetDefaults[configNoDefault](nil)) })
	})

	t.Run("Test config", func(t *testing.T) {
		var cfg configNoDefault
		assert.NotPanics(t, func() { options.Apply(ctx, &cfg, options.SetDefaults[configNoDefault](nil)) })
		assert.Nil(t, cfg.Values)
	})

	t.Run("Test config with provided default method", func(t *testing.T) {
		var cfg configWithDefault
		var setDefaultFn = func(c *configWithDefault) { c.Values = []string{"custom"} }
		options.Apply(ctx, &cfg, options.SetDefaults(setDefaultFn))
		assert.Equal(t, []string{"custom"}, cfg.Values)
	})

	t.Run("Test config with structure default method", func(t *testing.T) {
		var cfg configWithDefault
		options.Apply(ctx, &cfg, options.SetDefaults[configWithDefault](nil))
		assert.Equal(t, []string{"default"}, cfg.Values)
	})

}

func TestNewWithDefaults(t *testing.T) {
	var ctx = context.Background()
	var appender = func() options.WithFn[configWithDefault] {
		return func(o *configWithDefault) { o.Values = append(o.Values, "foo", "bar") }
	}

	t.Run("Test new config with provided default method", func(t *testing.T) {
		var setDefaultFn = func(c *configWithDefault) { c.Values = []string{"custom"} }
		var result = options.NewWithDefaults(ctx, setDefaultFn, appender())

		assert.NotNil(t, result)
		assert.EqualValues(t, []string{"custom", "foo", "bar"}, result.Values)
	})

	t.Run("Test new config with structure default method", func(t *testing.T) {
		var result = options.NewWithDefaults(ctx, nil, appender())

		assert.NotNil(t, result)
		assert.EqualValues(t, []string{"default", "foo", "bar"}, result.Values)
	})
}
