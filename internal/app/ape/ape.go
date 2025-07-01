package ape

import "errors"

type Violation struct {
	Field       string
	Description string
}

type BusinessError struct {
	reason     string
	message    string
	violations []Violation
	cause      error
}

func (e *BusinessError) Error() string {
	return e.message
}

func (e *BusinessError) Unwrap() error {
	return e.cause
}

func (e *BusinessError) Is(target error) bool {
	var be *BusinessError
	if errors.As(target, &be) {
		return e.reason == be.reason
	}
	return false
}

func (e *BusinessError) Reason() string {
	return e.reason
}

func (e *BusinessError) Violations() []Violation {
	return e.violations
}
