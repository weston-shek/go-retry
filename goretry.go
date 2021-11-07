package goretry

import (
	"time"
)

// For functions return `ERROR` ONLY.
// Returned error will be involved in the conditional judgment.
// e.g. err := func(); op.Judgment(err).
func (op RetryOption) ReDo(f func() error) (err error) {
	// The default judgment to stop retry is `err == nil`
	if op.judgment == nil {
		op.judgment = func(arg ...interface{}) bool { return arg[0] == nil }
	}
	for i := 0; i < op.attempt; i++ {
		if i > 0 {
			time.Sleep(time.Duration(op.backoff) * time.Duration(op.attempt) * time.Second)
		}
		// Get the results and judge
		err = f()
		if op.judgment(err) {
			return
		}
	}
	err = ErrMaxApt
	return
}

// For functions return `DATA, ERROR`.
// All returned data will be involved in the conditional judgment.
// e.g. data, err := func(); op.Judgment(data, err).
func (op RetryOption) ReTry(f func() (interface{}, error)) (data interface{}, err error) {
	// The default judgment to stop retry is `err == nil`
	if op.judgment == nil {
		op.judgment = func(arg ...interface{}) bool { return arg[0] == nil }
	}
	for i := 0; i < op.attempt; i++ {
		if i > 0 {
			time.Sleep(time.Duration(op.backoff) * time.Duration(op.attempt) * time.Second)
		}
		// Get the results and judge
		data, err = f()
		// End retry after data or error meets judgment condition
		if op.judgment(data, err) {
			return
		}
	}
	err = ErrMaxApt
	return
}

// Create a new retry attempt option
func NewRetryOption() *RetryOption {
	return &RetryOption{}
}

// Set the number of retry attempts,
// valid input range is natural numbers.
// The input value will be checked,
// and invalid input will result in a error return
func (op *RetryOption) SetAttempt(a int) (*RetryOption, error) {
	if a < 1 {
		return op, ErrInput
	}
	op.attempt = a
	return op, nil
}

// Set the number of retry attempts,
// valid input range is natural numbers.
// The input value will be checked,
// and invalid input will result in a panic
func (op *RetryOption) MustSetAttempt(a int) *RetryOption {
	if a < 1 {
		panic(ErrInput)
	}
	op.attempt = a
	return op
}

// Set retry attempt intervals,
// valid input range is natural numbers.
// The input value will be checked,
// and invalid input will result in a error return
func (op *RetryOption) SetBackoff(b int) (*RetryOption, error) {
	if b < 0 {
		return op, ErrInput
	}
	op.backoff = b
	return op, nil
}

// Set retry attempt intervals,
// valid input range is natural numbers.
// The input value will be checked,
// and invalid input will result in a panic
func (op *RetryOption) MustSetBackoff(b int) *RetryOption {
	if b < 0 {
		panic(ErrInput)
	}
	op.backoff = b
	return op
}

// Set the retry judgment conditions
func (op *RetryOption) SetJudgment(f func(arg ...interface{}) bool) *RetryOption {
	op.judgment = f
	return op
}
