package goretry

import (
	"errors"
	"testing"
)

var (
	ErrOdd  = errors.New("Odd number")
	ErrEven = errors.New("Even number")
)

func TestReDo(t *testing.T) {
	var tests = []struct {
		option RetryOption
		seed   int
		want   error
	}{
		// the option is set to make sure so that
		// retries will be executed the maximum number of times.
		{RetryOption{
			attempt:  6,
			backoff:  0,
			judgment: func(err ...interface{}) bool { return err[0] != ErrEven },
		}, 8, ErrMaxApt},
		{RetryOption{
			attempt:  6,
			backoff:  0,
			judgment: func(err ...interface{}) bool { return err[0] == ErrOdd },
		}, 9, ErrOdd},
	}
	for _, test := range tests {
		got := test.option.ReDo(
			func() error {
				if mod := test.seed % 2; mod == 0 {
					return ErrEven
				}
				return ErrOdd
			})
		if got != test.want {
			t.Errorf("Result is %v, want %v", got, test.want)
		}
	}
}

func TestReTry(t *testing.T) {
	var tests = []struct {
		option RetryOption
		seed   int
		want   interface{}
	}{
		// If use data to participate in conditional judgments,
		// it must use type assertions.
		{RetryOption{
			attempt:  6,
			backoff:  0,
			judgment: func(arg ...interface{}) bool { return arg[0].(int) >= 3 },
		}, 0, 3},
		{RetryOption{
			attempt:  10,
			backoff:  0,
			judgment: func(arg ...interface{}) bool { return arg[0].(int) >= 10 },
		}, 0, 10},
	}
	for _, test := range tests {
		got, err := test.option.ReTry(
			func() (interface{}, error) {
				test.seed++
				return test.seed, nil
			},
		)
		if err != nil {
			t.Errorf("Retry error: %v", err)
		}
		if got != test.want {
			t.Errorf("Result is %v, want %v", got, test.want)
		}
	}
}

func TestSetAttempt(t *testing.T) {
	var tests = []struct {
		apt int
		err error
	}{
		{1, nil},
		{2, nil},
		{100, nil},
		{600, nil},
		{-1, ErrInput},
	}
	for _, test := range tests {
		op := NewRetryOption()
		_, err := op.SetAttempt(test.apt)
		if err != test.err {
			t.Errorf("Set attempt ended with failure: %v", err)
		}
		if err == nil && op.attempt != test.apt {
			t.Errorf("Setting result is %v, want %v", op.attempt, test.apt)
		}
	}
}

func TestSetBackoff(t *testing.T) {
	var tests = []struct {
		apt int
		err error
	}{
		{1, nil},
		{2, nil},
		{100, nil},
		{600, nil},
		{-1, ErrInput},
	}
	for _, test := range tests {
		op := NewRetryOption()
		_, err := op.SetAttempt(test.apt)
		if err != test.err {
			t.Errorf("Set attempt ended with failure: %v", err)
		}
		if err == nil && op.attempt != test.apt {
			t.Errorf("Setting result is %v, want %v", op.attempt, test.apt)
		}
	}
}
