package websocket

import (
	"encoding/json"
	"fmt"
	"sync"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"github.com/ktruedat/chatly/internal/application/services"
	"github.com/ktruedat/chatly/internal/domain"
)

// wsWriter is an interface for writing WebSocket messages
type wsWriter interface {
	WriteMessage(messageType int, data []byte) error
}

type service struct {
	conn  wsWriter
	mutex sync.Mutex
}

func NewService(conn *websocket.Conn) services.WebSocketService {
	return &service{
		conn: conn,
	}
}

type BaseMessage struct {
	Type      domain.MessageType `json:"type"`
	RequestID uuid.UUID          `json:"request_id"`
}

type AckMessage struct {
	BaseMessage
	Status domain.AckStatus `json:"status"`
}

type ProgressMessage struct {
	BaseMessage
	Stage   string `json:"stage"`
	Message string `json:"message"`
}

type TokenMessage struct {
	BaseMessage
	Token        string `json:"token"`
	IsFinalToken bool   `json:"is_final_token"`
}

type FinalResponseMessage struct {
	BaseMessage
	ResponseText string `json:"response_text"`
	ModelUsed    string `json:"model_used"`
}

type ErrorMessage struct {
	BaseMessage
	ErrorCode domain.ErrorCode `json:"error_code"`
	Message   string           `json:"message"`
}

func (s *service) SendAck(requestID uuid.UUID, status domain.AckStatus) error {
	msg := AckMessage{
		BaseMessage: BaseMessage{
			Type:      domain.MessageTypeAck,
			RequestID: requestID,
		},
		Status: status,
	}
	return s.sendMessage(msg)
}

func (s *service) SendProgress(requestID uuid.UUID, stage domain.ProcessingStage, message string) error {
	msg := ProgressMessage{
		BaseMessage: BaseMessage{
			Type:      domain.MessageTypeProgress,
			RequestID: requestID,
		},
		Stage:   stage.String(),
		Message: message,
	}
	return s.sendMessage(msg)
}

func (s *service) SendToken(requestID uuid.UUID, token string, isFinalToken bool) error {
	msg := TokenMessage{
		BaseMessage: BaseMessage{
			Type:      domain.MessageTypeToken,
			RequestID: requestID,
		},
		Token:        token,
		IsFinalToken: isFinalToken,
	}
	return s.sendMessage(msg)
}

func (s *service) SendFinalResponse(requestID uuid.UUID, responseText string, modelUsed string) error {
	msg := FinalResponseMessage{
		BaseMessage: BaseMessage{
			Type:      domain.MessageTypeFinalResponse,
			RequestID: requestID,
		},
		ResponseText: responseText,
		ModelUsed:    modelUsed,
	}
	return s.sendMessage(msg)
}

func (s *service) SendError(requestID uuid.UUID, errorCode domain.ErrorCode, message string) error {
	msg := ErrorMessage{
		BaseMessage: BaseMessage{
			Type:      domain.MessageTypeError,
			RequestID: requestID,
		},
		ErrorCode: errorCode,
		Message:   message,
	}
	return s.sendMessage(msg)
}

func (s *service) sendMessage(msg interface{}) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	data, err := json.Marshal(msg)
	if err != nil {
		return fmt.Errorf("failed to marshal message: %w", err)
	}

	if err := s.conn.WriteMessage(websocket.TextMessage, data); err != nil {
		return fmt.Errorf("failed to write message: %w", err)
	}

	return nil
}
