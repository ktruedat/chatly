package domain

import "errors"

var (
	// Validation errors
	ErrInvalidInputType       = errors.New("invalid input type")
	ErrContentRequired        = errors.New("content is required for text requests")
	ErrImageRequired          = errors.New("image is required for image requests")
	ErrImageNotAllowedForText = errors.New("image not allowed for text-only requests")
	ErrInvalidRequestID       = errors.New("invalid request ID")

	// Processing errors
	ErrCategorizationFailed = errors.New("failed to categorize request")
	ErrModelNotFound        = errors.New("model not found for category")
	ErrRequestCancelled     = errors.New("request was cancelled")

	// External service errors
	ErrOpenRouterAPIFailed   = errors.New("openrouter API request failed")
	ErrOpenRouterRateLimited = errors.New("openrouter rate limit exceeded")
	ErrOpenRouterInvalidKey  = errors.New("invalid openrouter API key")
)
