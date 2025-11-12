package errors

import (
	apperrors "github.com/ktruedat/chatly/pkg/errors"
)

type ErrorOpt func(*apperrors.GenericError)

func WithCauseError(cause error) ErrorOpt {
	return func(e *apperrors.GenericError) {
		e.Cause = cause
	}
}

func NewInternalError(err error) apperrors.ApplicationError {
	return &apperrors.GenericError{
		Code:    ErrorCodeInternal,
		Message: "An internal error occurred",
		Cause:   err,
	}
}

func NewBadRequestError(msg string, opts ...ErrorOpt) apperrors.ApplicationError {
	ge := &apperrors.GenericError{
		Code:    ErrorCodeBadRequest,
		Message: msg,
	}
	for _, opt := range opts {
		opt(ge)
	}
	return ge
}

func NewModelError(msg string, cause error) apperrors.ApplicationError {
	return &apperrors.GenericError{
		Code:    ErrorCodeModelError,
		Message: msg,
		Cause:   cause,
	}
}

func NewNotFoundError(msg string) apperrors.ApplicationError {
	return &apperrors.GenericError{
		Code:    ErrorCodeNotFound,
		Message: msg,
	}
}
