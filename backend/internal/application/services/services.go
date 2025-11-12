package services

import (
	"context"

	"github.com/google/uuid"
	"github.com/ktruedat/chatly/internal/domain"
)

// ChatService handles the business logic for chat requests
type ChatService interface {
	// ProcessRequest processes a chat request and returns the category and model to use
	ProcessRequest(ctx context.Context, request *domain.ChatRequest) (*domain.ChatResponse, error)

	// CategorizeRequest categorizes a text request into easy/advanced/coding
	CategorizeRequest(ctx context.Context, content string) (domain.Category, error)

	// DetermineCategory determines the category based on input type
	DetermineCategory(ctx context.Context, request *domain.ChatRequest) (domain.Category, error)
}

// OpenRouterService handles communication with OpenRouter API
type OpenRouterService interface {
	// SendRequest sends a request to OpenRouter and returns the response
	SendRequest(ctx context.Context, model string, content string, imageBase64 *string) (string, error)

	// StreamRequest sends a request and streams the response token by token
	StreamRequest(ctx context.Context, model string, content string, imageBase64 *string, tokenCallback func(token string) error) error
}

// WebSocketService manages WebSocket connections and message routing
type WebSocketService interface {
	// SendAck sends an acknowledgment message
	SendAck(requestID uuid.UUID, status domain.AckStatus) error

	// SendProgress sends a progress update
	SendProgress(requestID uuid.UUID, stage domain.ProcessingStage, message string) error

	// SendToken sends a streaming token
	SendToken(requestID uuid.UUID, token string, isFinalToken bool) error

	// SendFinalResponse sends the final complete response
	SendFinalResponse(requestID uuid.UUID, responseText string, modelUsed string) error

	// SendError sends an error message
	SendError(requestID uuid.UUID, errorCode domain.ErrorCode, message string) error
}
