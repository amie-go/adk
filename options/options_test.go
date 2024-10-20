package options_test

import (
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

func (obj structOpt) Apply(c *config) { c.Values = append(c.Values, obj...) }

// --------------------------------------------

func WithStructPtrOpt(someValues ...string) options.With[config] {
	return &structPtrOpt{values: someValues}
}

type structPtrOpt struct {
	values []string
}

func (obj *structPtrOpt) Apply(c *config) { c.Values = append(c.Values, obj.values...) }

// --------------------------------------------

func TestApply(t *testing.T) {
	t.Run("Test nil config", func(t *testing.T) {
		var nilConfig *config
		assert.NotPanics(t, func() { options.Apply(nilConfig, nil) })
		assert.NotPanics(t, func() { options.Apply(nilConfig, WithVerbose(true)) })
	})

	t.Run("Test config", func(t *testing.T) {
		var config config
		assert.NotPanics(t, func() { options.Apply(&config, nil) })
		options.Apply(&config, WithVerbose(true))
		assert.True(t, config.Verbose)
		assert.NotPanics(t, func() { options.Apply(&config, nil, WithVerbose(false), nil) })
		assert.False(t, config.Verbose)
	})

	t.Run("Test config with struct opt", func(t *testing.T) {
		var config config
		options.Apply(&config, WithStructOpt("a", "b"))
		assert.Equal(t, []string{"a", "b"}, config.Values)
		options.Apply(&config, WithStructPtrOpt("c", "d"))
		assert.Equal(t, []string{"a", "b", "c", "d"}, config.Values)
	})
}

// --------------------------------------------

type configWithDefault struct {
	Values []string
}

func (c *configWithDefault) SetDefault() {
	c.Values = []string{"default"}
}

// --------------------------------------------

type configWithValidate struct {
	err error
}

func (c *configWithValidate) Validate() error {
	return c.err
}

func TestNew(t *testing.T) {
	t.Run("Test new any", func(t *testing.T) {
		result, err := options.New[any]()
		assert.NotNil(t, result)
		assert.NoError(t, err)
	})

	t.Run("Test new config", func(t *testing.T) {
		result, err := options.New(WithVerbose(true))
		assert.NotNil(t, result)
		assert.EqualValues(t, true, result.Verbose)
		assert.NoError(t, err)
	})

	t.Run("Test new config with default", func(t *testing.T) {
		var appender = func(args ...string) options.WithFn[configWithDefault] {
			return func(c *configWithDefault) { c.Values = append(c.Values, args...) }
		}
		result, err := options.New(appender("foo", "bar"))
		assert.NotNil(t, result)
		assert.EqualValues(t, []string{"default", "foo", "bar"}, result.Values)
		assert.NoError(t, err)
	})

	var setValidateErr = func(err error) options.WithFn[configWithValidate] {
		return func(c *configWithValidate) { c.err = err }
	}
	t.Run("Test new config with validate no error", func(t *testing.T) {
		result, err := options.New(setValidateErr(nil))
		assert.NotNil(t, result)
		assert.NoError(t, err)
	})
	t.Run("Test new config with validate with error", func(t *testing.T) {
		result, err := options.New(setValidateErr(assert.AnError))
		assert.Nil(t, result)
		assert.EqualError(t, err, assert.AnError.Error())
	})
}
