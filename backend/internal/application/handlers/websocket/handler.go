package websocket

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"github.com/ktruedat/chatly/internal/application/common/errors"
	"github.com/ktruedat/chatly/internal/application/config"
	"github.com/ktruedat/chatly/internal/application/handlers"
	"github.com/ktruedat/chatly/internal/application/services"
	"github.com/ktruedat/chatly/internal/application/services/chat"
	wsService "github.com/ktruedat/chatly/internal/application/services/websocket"
	"github.com/ktruedat/chatly/internal/domain"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

type handler struct {
	router            chi.Router
	cfg               *config.Config
	openRouterService services.OpenRouterService
}

func NewHandlers(
	router chi.Router,
	cfg *config.Config,
	openRouterService services.OpenRouterService,
) handlers.Handlers {
	return &handler{
		router:            router,
		cfg:               cfg,
		openRouterService: openRouterService,
	}
}

func (h *handler) Register() error {
	h.router.Get("/ws", h.handleWebSocket)
	return nil
}

func (h *handler) handleWebSocket(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("Failed to upgrade connection: %v", err)
		return
	}
	defer conn.Close()

	log.Printf("New WebSocket connection from %s", r.RemoteAddr)

	wsSvc := wsService.NewService(conn)
	chatSvc := chat.NewService(h.cfg, h.openRouterService)

	for {
		_, message, err := conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("WebSocket error: %v", err)
			}
			break
		}

		var baseMsg struct {
			Type domain.MessageType `json:"type"`
		}
		if err := json.Unmarshal(message, &baseMsg); err != nil {
			log.Printf("Failed to parse message type: %v", err)
			continue
		}

		switch baseMsg.Type {
		case domain.MessageTypeSubmitRequest:
			h.handleSubmitRequest(r.Context(), message, wsSvc, chatSvc)
		case domain.MessageTypeCancelRequest:
			h.handleCancelRequest(r.Context(), message, wsSvc)
		default:
			log.Printf("Unknown message type: %s", baseMsg.Type)
		}
	}

	log.Printf("WebSocket connection closed from %s", r.RemoteAddr)
}

type ClientSubmitRequest struct {
	Type        string                 `json:"type"`
	RequestID   string                 `json:"request_id"`
	InputType   string                 `json:"input_type"`
	Content     *string                `json:"content"`
	ImageBase64 *string                `json:"image_base64"`
	Metadata    map[string]interface{} `json:"metadata"`
}

func (h *handler) handleSubmitRequest(
	ctx context.Context,
	message []byte,
	wsSvc services.WebSocketService,
	chatSvc services.ChatService,
) {
	var clientReq ClientSubmitRequest
	if err := json.Unmarshal(message, &clientReq); err != nil {
		log.Printf("Failed to parse submit_request: %v", err)
		return
	}

	// Parse request ID
	requestID, err := uuid.Parse(clientReq.RequestID)
	if err != nil {
		log.Printf("Invalid request ID: %v", err)
		if sendErr := wsSvc.SendError(
			uuid.Nil,
			domain.ErrorCodeBadRequest,
			"Invalid request ID format",
		); sendErr != nil {
			log.Printf("Failed to send error response: %v", sendErr)
		}
		return
	}

	// Send acknowledgment
	if err := wsSvc.SendAck(requestID, domain.AckStatusAccepted); err != nil {
		log.Printf("Failed to send ack for request %s: %v", requestID, err)
		return
	}

	// Create domain request
	inputType := domain.InputType(clientReq.InputType)
	chatRequest, err := domain.NewChatRequest(
		requestID,
		inputType,
		clientReq.Content,
		clientReq.ImageBase64,
		clientReq.Metadata,
	)
	if err != nil {
		if sendErr := wsSvc.SendError(requestID, domain.ErrorCodeBadRequest, err.Error()); sendErr != nil {
			log.Printf("Failed to send error response for request %s: %v", requestID, sendErr)
		}
		return
	}

	// Process in background
	go h.processRequest(ctx, chatRequest, wsSvc, chatSvc)
}

func (h *handler) processRequest(
	ctx context.Context,
	request *domain.ChatRequest,
	wsSvc services.WebSocketService,
	chatSvc services.ChatService,
) {
	requestID := request.RequestID

	var category domain.Category
	var err error
	if request.InputType == domain.InputTypeText {
		if err := wsSvc.SendProgress(requestID, domain.StageCategorizing, "Analyzing request type..."); err != nil {
			log.Printf("Failed to send progress for request %s: %v", requestID, err)
			return
		}

		category, err = chatSvc.CategorizeRequest(ctx, request.GetContentOrDefault())
		if err != nil {
			log.Printf("Categorization failed for request %s: %v", requestID, err)
			if sendErr := wsSvc.SendError(
				requestID,
				domain.ErrorCodeInternalError,
				"Failed to categorize request",
			); sendErr != nil {
				log.Printf("Failed to send error response for request %s: %v", requestID, sendErr)
			}
			return
		}

		if err := wsSvc.SendProgress(
			requestID,
			domain.StageDispatching,
			fmt.Sprintf("Category determined: %s", category),
		); err != nil {
			log.Printf("Failed to send progress for request %s: %v", requestID, err)
			return
		}
	} else {
		// Image requests skip categorization
		category = domain.CategoryImage
		if err := wsSvc.SendProgress(requestID, domain.StageDispatching, "Category determined: image"); err != nil {
			log.Printf("Failed to send progress for request %s: %v", requestID, err)
			return
		}
	}

	// Get model for category
	var model string
	switch category {
	case domain.CategoryEasy:
		model = h.cfg.Models.Easy
	case domain.CategoryAdvanced:
		model = h.cfg.Models.Advanced
	case domain.CategoryCoding:
		model = h.cfg.Models.Coding
	case domain.CategoryImage:
		if request.InputType == domain.InputTypeImageHard {
			model = h.cfg.Models.ImageHard
		} else {
			model = h.cfg.Models.Image
		}
	default:
		if sendErr := wsSvc.SendError(requestID, domain.ErrorCodeInternalError, "Unknown category"); sendErr != nil {
			log.Printf("Failed to send error response for request %s: %v", requestID, sendErr)
		}
		return
	}

	if err := wsSvc.SendProgress(
		requestID,
		domain.StageProcessing,
		fmt.Sprintf("Requesting response from %s...", model),
	); err != nil {
		log.Printf("Failed to send progress for request %s: %v", requestID, err)
		return
	}

	var responseBuilder strings.Builder
	if err := h.openRouterService.StreamRequest(
		ctx, model, request.GetContentOrDefault(), request.ImageBase64, func(token string) error {
			responseBuilder.WriteString(token)
			return wsSvc.SendToken(requestID, token, false)
		},
	); err != nil {
		log.Printf("OpenRouter request failed for request %s: %v", requestID, err)

		errorCode := domain.ErrorCodeModelError
		if domainErr := errors.NewModelError("Model request failed", err); domainErr != nil {
			errorCode = domain.ErrorCode(domainErr.ErrCode().String())
		}

		if sendErr := wsSvc.SendError(requestID, errorCode, fmt.Sprintf("API error: %v", err)); sendErr != nil {
			log.Printf("Failed to send error response for request %s: %v", requestID, sendErr)
		}
		return
	}

	responseText := responseBuilder.String()
	if err := wsSvc.SendFinalResponse(requestID, responseText, model); err != nil {
		log.Printf("Failed to send final response for request %s: %v", requestID, err)
	}
}

type ClientCancelRequest struct {
	Type      string `json:"type"`
	RequestID string `json:"request_id"`
}

func (h *handler) handleCancelRequest(_ context.Context, message []byte, wsSvc services.WebSocketService) {
	var cancelReq ClientCancelRequest
	if err := json.Unmarshal(message, &cancelReq); err != nil {
		log.Printf("Failed to parse cancel_request: %v", err)
		return
	}

	requestID, err := uuid.Parse(cancelReq.RequestID)
	if err != nil {
		log.Printf("Invalid request ID for cancellation: %v", err)
		return
	}

	// TODO: Implement proper cancellation logic
	if err := wsSvc.SendProgress(requestID, domain.StageCancelling, "Cancellation requested"); err != nil {
		log.Printf("Failed to send cancellation progress for request %s: %v", requestID, err)
	}
	log.Printf("Cancellation requested for request %s (not yet implemented)", requestID)
}
