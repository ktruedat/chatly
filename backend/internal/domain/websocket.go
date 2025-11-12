package domain

// MessageType represents the type of WebSocket message
type MessageType string

const (
	MessageTypeSubmitRequest MessageType = "submit_request"
	MessageTypeCancelRequest MessageType = "cancel_request"
	MessageTypeAck           MessageType = "ack"
	MessageTypeProgress      MessageType = "progress"
	MessageTypeToken         MessageType = "token"
	MessageTypeFinalResponse MessageType = "final_response"
	MessageTypeError         MessageType = "error"
)

func (m MessageType) String() string {
	return string(m)
}

// AckStatus represents the status of an acknowledgment message
type AckStatus string

const (
	AckStatusAccepted AckStatus = "accepted"
	AckStatusRejected AckStatus = "rejected"
)

func (a AckStatus) String() string {
	return string(a)
}

// ErrorCode represents standardized error codes for WebSocket responses
type ErrorCode string

const (
	ErrorCodeBadRequest    ErrorCode = "bad_request"
	ErrorCodeInternalError ErrorCode = "internal_error"
	ErrorCodeModelError    ErrorCode = "model_error"
	ErrorCodeRateLimited   ErrorCode = "rate_limited"
	ErrorCodeUnauthorized  ErrorCode = "unauthorized"
	ErrorCodeNotFound      ErrorCode = "not_found"
	ErrorCodeCancelled     ErrorCode = "cancelled"
)

func (e ErrorCode) String() string {
	return string(e)
}
