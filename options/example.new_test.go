package options_test

import (
	"github.com/amie-go/adk/options"
)

// --------------------------------------------
// Package creator: Define in a file your private configuration and public configuration functions

// Define a private configuration structure
type settings struct {
	Name   string
	Cities []string
}

// Define some public option functions
func WithSuffix(suffix string) options.WithFn[settings] {
	return func(o *settings) { o.Name += suffix }
}

// Define another kind of public option functions
func WithCities(values ...string) options.With[settings] {
	return cities(values)
}

type cities []string

func (obj cities) Apply(dst *settings) { dst.Cities = append(dst.Cities, obj...) }

// --------------------------------------------
// Package creator: Define your constructor function

// Define a constructor function
func NewMyClass(opts ...options.With[settings]) *MyClass {
	var result = MyClass{config: settings{
		Name: "foo",
	}}
	options.Apply(&result.config, opts...)
	return &result
}

type MyClass struct {
	config settings
}

// --------------------------------------------
// Package user: Use the package

func ExampleNew() {
	_ = NewMyClass(WithSuffix("bar"))
}
