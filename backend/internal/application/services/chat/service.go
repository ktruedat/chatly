package chat

import (
	"context"
	"fmt"
	"strings"

	"github.com/ktruedat/chatly/internal/application/config"
	"github.com/ktruedat/chatly/internal/application/services"
	"github.com/ktruedat/chatly/internal/domain"
)

type service struct {
	cfg               *config.Config
	openRouterService services.OpenRouterService
}

func NewService(cfg *config.Config, openRouterService services.OpenRouterService) services.ChatService {
	return &service{
		cfg:               cfg,
		openRouterService: openRouterService,
	}
}

func (s *service) ProcessRequest(ctx context.Context, request *domain.ChatRequest) (*domain.ChatResponse, error) {
	category, err := s.DetermineCategory(ctx, request)
	if err != nil {
		return nil, fmt.Errorf("failed to determine category: %w", err)
	}

	// Get model for category
	var model string
	switch category {
	case domain.CategoryEasy:
		model = s.cfg.Models.Easy
	case domain.CategoryAdvanced:
		model = s.cfg.Models.Advanced
	case domain.CategoryCoding:
		model = s.cfg.Models.Coding
	case domain.CategoryImage:
		if request.InputType == domain.InputTypeImageHard {
			model = s.cfg.Models.ImageHard
		} else {
			model = s.cfg.Models.Image
		}
	default:
		return nil, domain.ErrModelNotFound
	}

	// Get content
	content := request.GetContentOrDefault()

	// Send request to OpenRouter
	responseText, err := s.openRouterService.SendRequest(ctx, model, content, request.ImageBase64)
	if err != nil {
		return nil, fmt.Errorf("openrouter request failed: %w", err)
	}

	// Create response
	response := domain.NewChatResponse(request.RequestID, responseText, model, category)
	return response, nil
}

func (s *service) CategorizeRequest(ctx context.Context, content string) (domain.Category, error) {
	// Build categorization prompt
	prompt := strings.ReplaceAll(s.cfg.Categorization.Prompt, "{content}", content)

	// Send to categorizer model
	response, err := s.openRouterService.SendRequest(ctx, s.cfg.Models.Categorizer, prompt, nil)
	if err != nil {
		return "", fmt.Errorf("categorization failed: %w", err)
	}

	// Parse response
	category := strings.TrimSpace(strings.ToLower(response))

	switch category {
	case "easy":
		return domain.CategoryEasy, nil
	case "advanced":
		return domain.CategoryAdvanced, nil
	case "coding":
		return domain.CategoryCoding, nil
	default:
		// Default to advanced if we can't determine
		return domain.CategoryAdvanced, nil
	}
}

func (s *service) DetermineCategory(ctx context.Context, request *domain.ChatRequest) (domain.Category, error) {
	// Image requests don't need categorization
	if request.InputType == domain.InputTypeImage || request.InputType == domain.InputTypeImageHard {
		return domain.CategoryImage, nil
	}

	// Text requests need categorization
	if request.Content == nil || *request.Content == "" {
		return "", domain.ErrContentRequired
	}

	return s.CategorizeRequest(ctx, *request.Content)
}
