package retrier

import (
	"errors"
	"time"
)

type options struct {
	attemps *int
	timeout time.Duration
}

type Option func(options *options) error

func WithAttemps(attemps int) Option {
	return func(options *options) error {
		if attemps <= 0 {
			return errors.New("attemps must to be positive")
		}

		options.attemps = &attemps
		return nil
	}
}

func WithTimeOut(secondsNumber int) Option {
	return func(options *options) error {
		if secondsNumber <= 0 {
			return errors.New("secondsnumber must be positive")
		}

		options.timeout = time.Duration(secondsNumber) * time.Second
		return nil
	}
}

func Connect[T any](connFunc func() (T, error), opts ...Option) (T, error) {
	var (
		err     error = errors.New("__")
		payload T
		options options
	)

	for _, option := range opts {
		if err = option(&options); err != nil {
			return payload, err
		}
	}

	if options.attemps == nil {
		for err != nil {
			payload, err = connFunc()
		}

		return payload, err
	}

	for range *options.attemps {
		payload, err = connFunc()
		if err == nil {
			return payload, err
		}
	}

	return payload, err
}
