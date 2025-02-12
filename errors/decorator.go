package exterr

// DecoratorFn is a function that decorates an error.
type DecoratorFn func(error) error

// New creates a new error with the given decorators.
func New(err error, decorators ...DecoratorFn) error {
	if err == nil {
		return nil
	}
	for _, v := range decorators {
		if v != nil {
			err = v(err)
		}
	}
	return err
}
