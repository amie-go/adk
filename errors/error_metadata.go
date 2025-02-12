package exterr

import (
	"errors"
)

func WithMetadata(kvs ...KV) func(error) error {
	return func(err error) error {
		return &errMetadata{
			inner: err,
			kvs:   kvs,
		}
	}
}

func GetMetadata(err error) (result []KV) {
	for err != nil {
		if e, ok := err.(*errMetadata); ok {
			result = append(result, e.kvs...)
		}
		err = errors.Unwrap(err)
	}
	return result
}

type KV struct {
	Key   string
	Value any
}

type errMetadata struct {
	inner error
	kvs   []KV
}

func (e *errMetadata) Error() string {
	return e.inner.Error()
}

func (e *errMetadata) Unwrap() error {
	return e.inner
}
