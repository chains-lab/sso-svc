package ape

import (
	"errors"

	"google.golang.org/grpc/codes"
	"google.golang.org/protobuf/protoadapt"
)

type Error struct {
	code    codes.Code
	reason  string
	message string
	details []protoadapt.MessageV1
	cause   error
}

func (e *Error) Error() string {
	return e.message
}

func (e *Error) Unwrap() error {
	return e.cause
}

func (e *Error) Is(target error) bool {
	var be *Error
	if errors.As(target, &be) {
		return e.reason == be.reason
	}
	return false
}

func (e *Error) Reason() string {
	return e.reason
}

func (e *Error) Details() []protoadapt.MessageV1 {
	if e.details == nil {
		return nil
	}

	return e.details
}

func (e *Error) Code() codes.Code {
	return e.code
}
