package errors

import "fmt"

type Code interface {
	String() string
	ToHTTPCode() int
}

type ApplicationError interface {
	error
	ErrCode() Code
	ErrMessage() string
	ErrCause() error
}

type GenericError struct {
	Code    Code   `json:"code"`
	Message string `json:"message"`
	Cause   error  `json:"-"`
}

func (e *GenericError) Error() string {
	if e.Cause != nil {
		return fmt.Sprintf("%s: %v", e.Message, e.Cause)
	}
	return e.Message
}

func (e *GenericError) Unwrap() error {
	return e.Cause
}

func (e *GenericError) ErrCode() Code {
	return e.Code
}

func (e *GenericError) ErrMessage() string {
	return e.Message
}

func (e *GenericError) ErrCause() error {
	return e.Cause
}
