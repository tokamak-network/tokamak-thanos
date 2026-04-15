package retry

import (
	"context"
	"errors"
	"fmt"
	"time"
)

// ErrUnrecoverable wraps an error to signal that Do should stop retrying immediately.
// Use Unrecoverable(err) to wrap an error before returning it from the operation.
type ErrUnrecoverable struct {
	Err error
}

func (e *ErrUnrecoverable) Error() string { return e.Err.Error() }
func (e *ErrUnrecoverable) Unwrap() error { return e.Err }

// Unrecoverable wraps err so that Do stops retrying and returns immediately.
func Unrecoverable(err error) error {
	return &ErrUnrecoverable{Err: err}
}

// ErrFailedPermanently is an error raised by Do when the
// underlying Operation has been retried maxAttempts times.
type ErrFailedPermanently struct {
	attempts int
	LastErr  error
}

func (e *ErrFailedPermanently) Error() string {
	return fmt.Sprintf("operation failed permanently after %d attempts: %v", e.attempts, e.LastErr)
}

func (e *ErrFailedPermanently) Unwrap() error {
	return e.LastErr
}

type pair[T, U any] struct {
	a T
	b U
}

func Do2[T, U any](ctx context.Context, maxAttempts int, strategy Strategy, op func() (T, U, error)) (T, U, error) {
	f := func() (pair[T, U], error) {
		a, b, err := op()
		return pair[T, U]{a, b}, err
	}
	res, err := Do(ctx, maxAttempts, strategy, f)
	return res.a, res.b, err
}

// Do performs the provided Operation up to maxAttempts times
// with delays in between each retry according to the provided
// Strategy.
func Do[T any](ctx context.Context, maxAttempts int, strategy Strategy, op func() (T, error)) (T, error) {
	var empty, ret T
	var err error
	if maxAttempts < 1 {
		return empty, fmt.Errorf("need at least 1 attempt to run op, but have %d max attempts", maxAttempts)
	}

	for i := 0; i < maxAttempts; i++ {
		if ctx.Err() != nil {
			return empty, ctx.Err()
		}
		ret, err = op()
		if err == nil {
			return ret, nil
		}
		// Stop immediately for unrecoverable errors (no retries, no sleep)
		var unrecoverable *ErrUnrecoverable
		if errors.As(err, &unrecoverable) {
			return empty, unrecoverable.Err
		}
		// Don't sleep when we are about to exit the loop & return ErrFailedPermanently
		if i != maxAttempts-1 {
			time.Sleep(strategy.Duration(i))
		}
	}
	return empty, &ErrFailedPermanently{
		attempts: maxAttempts,
		LastErr:  err,
	}
}

// Do0 is like Do but for operations that return no value.
func Do0(ctx context.Context, maxAttempts int, strategy Strategy, op func() error) error {
	_, err := Do(ctx, maxAttempts, strategy, func() (struct{}, error) {
		return struct{}{}, op()
	})
	return err
}
