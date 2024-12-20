package options_test

import (
	"context"

	"github.com/amie-go/adk/options"
)

// --------------------------------------------
// Package creator: Define in a file your private configuration and public configuration functions

// Define a private configuration structure
type settingsDef struct {
	Name   string
	Cities []string
}

// Define a private option functions to set default values
func setDefault(o *settingsDef) {
	o.Name = "foo"
	o.Cities = []string{"Paris", "London"}
}

// Define some public option functions
func WithCity(value string) options.WithFn[settingsDef] {
	return func(o *settingsDef) {
		o.Cities = append(o.Cities, value)
	}
}

// --------------------------------------------
// Package creator: Define your constructor function

func NewMyClass2(opts ...options.With[settingsDef]) *MyClass2 {
	//var config = options.New(context.Background(), append([]options.With[settingsDef]{options.SetDefaults(setDefault)}, opts...)...)
	var config = options.NewWithDefaults(context.Background(), setDefault, opts...)

	return &MyClass2{config: config}
}

type MyClass2 struct {
	config *settingsDef
}

// --------------------------------------------
// Package user: Use the package

func ExampleNewWithDefaults() {
	_ = NewMyClass2(WithCity("Mexico"))
}
