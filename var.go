package goretry

import "errors"

var (
	ErrMaxApt = errors.New("Failed after maximum number of attempts to retry")
	ErrInput  = errors.New("The input value is outside the legal range")
)

type RetryOption struct {
	// Number of retries (>=1)
	attempt int

	// Fallback waiting time, unit: second (>=0)
	backoff int

	// Judgment expression,
	// the loop will stop while expression satisfied,
	// e.g. func(arg ...interface{}){return arg[0] == io.EOF}
	judgment func(arg ...interface{}) bool
}
