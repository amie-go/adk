# Description

This package allows you to implement the functional options pattern in your packages.

# Usage

## Package creator side

### 1 - Define your project configuration

A good practice is to define the configuration of your project as **private**.

```go
// Example of private configuration defined in any file of your project
type settings struct {
	Name 	string
	Cities	[]string
    //...
}

// Example of private configuration defined in a file internal/settings.go of your project
type Settings struct {
	Name 	string
	Cities	[]string
    //...
}
```

### 2 - Provide public configuration function

You can provide configuration function based on:

- simple function.<br>
  ```go
  func WithSuffix(suffix string) options.WithFn[settings] {
  	return func(dst *settings) { dst.Name += suffix }
  }
  ```

- structure that implements interface `options.With`.<br>
  ```go
  func WithCities(values ...string) options.With[settings] {
  	return cities(values)
  }

  type cities []string
  
  func (obj cities) Apply(ctx context.Context, dst *settings) { dst.Cities = append(dst.Cities, obj...) }
  ```

> [!TIP]
> For simple function with `context.Context` argument, you can use `options.WithCtxFn`.

### 3 - Generate the configuration

A good practice is to create a `New` instance creation function in your package and generate the configuration using one of the methods above.

```go
func New(opts ...options.With[settings]) *MyClass {
	return &MyClass{config: options.New(opts...)}
}

type MyClass struct {
	config *settings
}
```

You can get your configuration applied by using function:

- `Apply`.<br>
  ```go
  	// Create your configuration with some default values
  	var config = &settings{
  		Name: "foo"
  	}
  	// Apply options on your configuration
  	options.Apply(context.Background(), config, opts...)
  ```

- `New`.<br>
  ```go
  	// Generate your configuration with provided options applied to it
  	var config = options.New(context.Background(), opts...)
  ```

- `NewWithDefaults`.<br>
  ```go
  	// Define a private option functions to set default values
  	func setDefault(o *settingsDef) {
  		o.Name = "foo"
  		o.Cities = []string{"Paris", "London"}
  	}
  
  	// Generate your configuration with provided options applied to it
  	var config = options.NewWithDefaults(context.Background(), setDefault, opts...)
  ```

> [!TIP]
> Calling `options.NewWithDefaults(context.Background(), nil, opts...)` will call (if it exists) the `SetDefaults()` function of your configuration structure.

## User side

```go
func baz() {
	var ctx = context.Background()
	var client = foo.New(ctx, foo.WithSuffix("bar"), foo.WithCities("London", "Paris"))
}
```