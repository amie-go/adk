# Description

This package allows you to implement the functional options pattern in your packages.

## Usage

### Package creator side

#### 1 - Define your project configuration

A good practice is to define the settings of your project as **private**.

```go
// Example of private configuration defined in any file of your project
type settings struct {
	Name 	string
	Cities	[]string
    //...
}

// Example of private configuration defined in a file internal/settings.go of your project
type Settings struct {
	Name string
    //...
}
```

#### 2 - Define public configuration function

You can create simple configuration function.

```go
func WithSuffix(suffix string) options.WithFn[settings] {
	return func(dst *settings) { dst.Name += suffix }
}
```

And same more complex configuration function.

```go
func WithCities(values ...string) options.With[settings] {
	return cities(values)
}

type cities []string

func (obj cities) Apply(dst *settings) { dst.Cities = append(dst.Cities, obj...) }
```

#### 3 - Generate the configuration

You can modify your configuration by applying the configuration functions.
```go
	// Create your settings structure with some default values
	var config = &settings{
		Name: "foo"
	}
	// Apply options on your settings structure
	options.Apply(config, opts...)
```

You can also directly generate the configuration with configuration functions applied.

```go
	// Generate your settings structure with provided options applied to it
	config, err := options.New(opts...)
```

> [!TIP]
> `New` calls (if exists) the functions from your configuration structure:
> - `SetDefault()` in first to init default values (before applying configuration functions).
> - `Validate() error` in last to do validation of your configuration (after applying configuration functions).

A good practice is to create a *New* instance creator function in your package, and generates the configuration as above.

```go
func New(opts ...options.With[settings]) *MyClass {
	// Generate your settings structure with provided options applied to it
	config, _ := options.New(opts...)
	return &MyClass{config: config}
}

type MyClass struct {
	config *settings
}
```

### User side

```go
func foo() {
	var client = foo.New(foo.WithSuffix("foo"), foo.WithCities("London", "Paris"))
}
```