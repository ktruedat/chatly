package errors

import "net/http"

type ErrorCode string

func (e ErrorCode) String() string {
	return string(e)
}

func (e ErrorCode) ToHTTPCode() int {
	if code, exists := serviceErrorsToHTTPCode[e]; exists {
		return code
	}
	return http.StatusInternalServerError
}

const (
	ErrorCodeInternal   ErrorCode = "internal_error"
	ErrorCodeBadRequest ErrorCode = "bad_request"
	ErrorCodeModelError ErrorCode = "model_error"
	ErrorCodeNotFound   ErrorCode = "not_found"
)

var serviceErrorsToHTTPCode = map[ErrorCode]int{
	ErrorCodeInternal:   http.StatusInternalServerError,
	ErrorCodeBadRequest: http.StatusBadRequest,
	ErrorCodeModelError: http.StatusServiceUnavailable,
	ErrorCodeNotFound:   http.StatusNotFound,
}
