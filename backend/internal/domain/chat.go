package domain

import (
	"time"

	"github.com/google/uuid"
)

// InputType represents the type of input from the user
type InputType string

const (
	InputTypeText      InputType = "Text"
	InputTypeImage     InputType = "Image"
	InputTypeImageHard InputType = "ImageHard"
)

func (i InputType) String() string {
	return string(i)
}

func (i InputType) IsValid() bool {
	switch i {
	case InputTypeText, InputTypeImage, InputTypeImageHard:
		return true
	default:
		return false
	}
}

// Category represents the complexity category for text requests
type Category string

const (
	CategoryEasy     Category = "easy"
	CategoryAdvanced Category = "advanced"
	CategoryCoding   Category = "coding"
	CategoryImage    Category = "image"
)

func (c Category) String() string {
	return string(c)
}

// ProcessingStage represents the current stage of request processing
type ProcessingStage string

const (
	StageCategorizing ProcessingStage = "categorizing"
	StageDispatching  ProcessingStage = "dispatching"
	StageProcessing   ProcessingStage = "processing"
	StageCancelling   ProcessingStage = "cancelling"
)

func (s ProcessingStage) String() string {
	return string(s)
}

// ChatRequest represents a user's chat request
type ChatRequest struct {
	RequestID   uuid.UUID              `json:"request_id"`
	InputType   InputType              `json:"input_type"`
	Content     *string                `json:"content"`
	ImageBase64 *string                `json:"image_base64"`
	Metadata    map[string]interface{} `json:"metadata,omitempty"`
	CreatedAt   time.Time              `json:"created_at"`
}

// NewChatRequest creates a new chat request
func NewChatRequest(
	requestID uuid.UUID,
	inputType InputType,
	content *string,
	imageBase64 *string,
	metadata map[string]interface{},
) (*ChatRequest, error) {
	// Validate input type
	if !inputType.IsValid() {
		return nil, ErrInvalidInputType
	}

	// Validate required fields based on input type
	if inputType == InputTypeText {
		if content == nil || *content == "" {
			return nil, ErrContentRequired
		}
		if imageBase64 != nil {
			return nil, ErrImageNotAllowedForText
		}
	}

	if inputType == InputTypeImage || inputType == InputTypeImageHard {
		if imageBase64 == nil || *imageBase64 == "" {
			return nil, ErrImageRequired
		}
	}

	return &ChatRequest{
		RequestID:   requestID,
		InputType:   inputType,
		Content:     content,
		ImageBase64: imageBase64,
		Metadata:    metadata,
		CreatedAt:   time.Now(),
	}, nil
}

// GetContentOrDefault returns content or a default message for image requests
func (r *ChatRequest) GetContentOrDefault() string {
	if r.Content != nil && *r.Content != "" {
		return *r.Content
	}
	if r.InputType == InputTypeImage || r.InputType == InputTypeImageHard {
		return "Describe this image"
	}
	return ""
}

// ChatResponse represents the complete response to a chat request
type ChatResponse struct {
	RequestID    uuid.UUID `json:"request_id"`
	ResponseText string    `json:"response_text"`
	ModelUsed    string    `json:"model_used"`
	Category     Category  `json:"category,omitempty"`
	CompletedAt  time.Time `json:"completed_at"`
}

// NewChatResponse creates a new chat response
func NewChatResponse(requestID uuid.UUID, responseText string, modelUsed string, category Category) *ChatResponse {
	return &ChatResponse{
		RequestID:    requestID,
		ResponseText: responseText,
		ModelUsed:    modelUsed,
		Category:     category,
		CompletedAt:  time.Now(),
	}
}
