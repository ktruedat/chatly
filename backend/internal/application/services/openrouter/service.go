package openrouter

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/ktruedat/chatly/internal/application/config"
	"github.com/ktruedat/chatly/internal/application/services"
	"github.com/ktruedat/chatly/internal/domain"
)

type service struct {
	cfg    *config.Config
	client *http.Client
}

func NewService(cfg *config.Config) services.OpenRouterService {
	return &service{
		cfg: cfg,
		client: &http.Client{
			Timeout: cfg.OpenRouter.Timeout,
		},
	}
}

type ChatMessage struct {
	Role    string        `json:"role"`
	Content []ContentPart `json:"content"`
}

type ContentPart struct {
	Type     string    `json:"type"`
	Text     string    `json:"text,omitempty"`
	ImageURL *ImageURL `json:"image_url,omitempty"`
}

type ImageURL struct {
	URL string `json:"url"`
}

type ChatRequest struct {
	Model    string        `json:"model"`
	Messages []ChatMessage `json:"messages"`
	Stream   bool          `json:"stream"`
}

type ChatResponse struct {
	Choices []struct {
		Message struct {
			Content string `json:"content"`
		} `json:"message"`
	} `json:"choices"`
	Error *APIError `json:"error,omitempty"`
}

type APIError struct {
	Message string `json:"message"`
	Code    int    `json:"code"`
}

type StreamChunk struct {
	Choices []struct {
		Delta struct {
			Content string `json:"content"`
		} `json:"delta"`
		FinishReason string `json:"finish_reason,omitempty"`
	} `json:"choices"`
	Error *APIError `json:"error,omitempty"`
}

func (s *service) SendRequest(ctx context.Context, model string, content string, imageBase64 *string) (string, error) {
	reqBody := s.buildRequest(model, content, imageBase64, false)

	resp, err := s.executeRequest(ctx, reqBody)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", s.handleErrorResponse(resp)
	}

	var chatResp ChatResponse
	if err := json.NewDecoder(resp.Body).Decode(&chatResp); err != nil {
		return "", fmt.Errorf("failed to decode response: %w", err)
	}

	if chatResp.Error != nil {
		return "", fmt.Errorf("%w: %s", domain.ErrOpenRouterAPIFailed, chatResp.Error.Message)
	}

	if len(chatResp.Choices) == 0 {
		return "", fmt.Errorf("%w: no choices in response", domain.ErrOpenRouterAPIFailed)
	}

	return chatResp.Choices[0].Message.Content, nil
}

func (s *service) StreamRequest(
	ctx context.Context,
	model string,
	content string,
	imageBase64 *string,
	tokenCallback func(token string) error,
) error {
	reqBody := s.buildRequest(model, content, imageBase64, true)

	resp, err := s.executeRequest(ctx, reqBody)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return s.handleErrorResponse(resp)
	}

	// OpenRouter uses Server-Sent Events (SSE) format
	// Each event is prefixed with "data: " and separated by newlines
	reader := resp.Body
	buf := make([]byte, 4096)
	var leftover []byte

	for {
		n, err := reader.Read(buf)
		if n > 0 {
			// Combine leftover from previous read with new data
			data := append(leftover, buf[:n]...)
			lines := strings.Split(string(data), "\n")

			// Save incomplete line for next iteration
			leftover = []byte(lines[len(lines)-1])
			lines = lines[:len(lines)-1]

			for _, line := range lines {
				line = strings.TrimSpace(line)

				// Skip empty lines and SSE comments
				if line == "" || strings.HasPrefix(line, ":") {
					continue
				}

				// SSE data lines start with "data: "
				if strings.HasPrefix(line, "data: ") {
					jsonData := strings.TrimPrefix(line, "data: ")

					// Check for [DONE] signal
					if jsonData == "[DONE]" {
						return nil
					}

					var chunk StreamChunk
					if err := json.Unmarshal([]byte(jsonData), &chunk); err != nil {
						// Skip invalid JSON (might be error messages or comments)
						continue
					}

					// Check for mid-stream errors
					if chunk.Error != nil {
						return fmt.Errorf("%w: %s", domain.ErrOpenRouterAPIFailed, chunk.Error.Message)
					}

					// Process content delta
					if len(chunk.Choices) > 0 && chunk.Choices[0].Delta.Content != "" {
						if err := tokenCallback(chunk.Choices[0].Delta.Content); err != nil {
							return err
						}
					}
				}
			}
		}

		if err != nil {
			if err == io.EOF {
				break
			}
			return fmt.Errorf("failed to read stream: %w", err)
		}
	}

	return nil
}

func (s *service) buildRequest(model string, content string, imageBase64 *string, stream bool) *ChatRequest {
	var contentParts []ContentPart

	// Add text content
	if content != "" {
		contentParts = append(
			contentParts, ContentPart{
				Type: "text",
				Text: content,
			},
		)
	}

	// Add image if present
	if imageBase64 != nil && *imageBase64 != "" {
		imageURL := fmt.Sprintf("data:image/jpeg;base64,%s", *imageBase64)
		contentParts = append(
			contentParts, ContentPart{
				Type: "image_url",
				ImageURL: &ImageURL{
					URL: imageURL,
				},
			},
		)
	}

	return &ChatRequest{
		Model: model,
		Messages: []ChatMessage{
			{
				Role:    "user",
				Content: contentParts,
			},
		},
		Stream: stream,
	}
}

func (s *service) executeRequest(ctx context.Context, reqBody *ChatRequest) (*http.Response, error) {
	jsonData, err := json.Marshal(*reqBody)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	url := fmt.Sprintf("%s/chat/completions", s.cfg.OpenRouter.BaseURL)
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", s.cfg.OpenRouter.APIKey))

	return s.client.Do(req)
}

func (s *service) handleErrorResponse(resp *http.Response) error {
	body, _ := io.ReadAll(resp.Body)

	var apiErr struct {
		Error *APIError `json:"error"`
	}
	if err := json.Unmarshal(body, &apiErr); err == nil && apiErr.Error != nil {
		if resp.StatusCode == http.StatusTooManyRequests {
			return fmt.Errorf("%w: %s", domain.ErrOpenRouterRateLimited, apiErr.Error.Message)
		}
		if resp.StatusCode == http.StatusUnauthorized {
			return fmt.Errorf("%w: %s", domain.ErrOpenRouterInvalidKey, apiErr.Error.Message)
		}
		return fmt.Errorf("%w: %s", domain.ErrOpenRouterAPIFailed, apiErr.Error.Message)
	}

	return fmt.Errorf(
		"%w: HTTP %d - %s",
		domain.ErrOpenRouterAPIFailed,
		resp.StatusCode,
		strings.TrimSpace(string(body)),
	)
}
