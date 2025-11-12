package websocket

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"github.com/ktruedat/chatly/internal/domain"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// mockWebSocketConn is a mock implementation of websocket.Conn for testing
type mockWebSocketConn struct {
	writtenMessages []mockMessage
	writeError      error
}

type mockMessage struct {
	messageType int
	data        []byte
}

func (m *mockWebSocketConn) WriteMessage(messageType int, data []byte) error {
	if m.writeError != nil {
		return m.writeError
	}
	m.writtenMessages = append(m.writtenMessages, mockMessage{
		messageType: messageType,
		data:        data,
	})
	return nil
}

func (m *mockWebSocketConn) getLastMessage() (int, []byte) {
	if len(m.writtenMessages) == 0 {
		return 0, nil
	}
	msg := m.writtenMessages[len(m.writtenMessages)-1]
	return msg.messageType, msg.data
}

func (m *mockWebSocketConn) getAllMessages() []mockMessage {
	return m.writtenMessages
}

// Create a wrapper service for testing that uses our mock
type testService struct {
	mockConn *mockWebSocketConn
	service  *service
}

func newTestService() *testService {
	mockConn := &mockWebSocketConn{
		writtenMessages: []mockMessage{},
	}

	// We need to use reflection or create the service differently
	// For now, we'll create the service with the real struct
	svc := &service{
		conn: nil, // We'll override sendMessage behavior
	}

	return &testService{
		mockConn: mockConn,
		service:  svc,
	}
}

// Override sendMessage for testing
func (ts *testService) sendMessage(msg interface{}) error {
	data, err := json.Marshal(msg)
	if err != nil {
		return err
	}
	return ts.mockConn.WriteMessage(websocket.TextMessage, data)
}

func TestSendAck(t *testing.T) {
	tests := []struct {
		name      string
		requestID uuid.UUID
		status    domain.AckStatus
		wantType  domain.MessageType
	}{
		{
			name:      "send accepted ack",
			requestID: uuid.New(),
			status:    domain.AckStatusAccepted,
			wantType:  domain.MessageTypeAck,
		},
		{
			name:      "send rejected ack",
			requestID: uuid.New(),
			status:    domain.AckStatusRejected,
			wantType:  domain.MessageTypeAck,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockConn := &mockWebSocketConn{}
			svc := &service{conn: mockConn}

			err := svc.SendAck(tt.requestID, tt.status)
			require.NoError(t, err)

			// Verify message was written
			require.Len(t, mockConn.writtenMessages, 1)

			msgType, data := mockConn.getLastMessage()
			assert.Equal(t, websocket.TextMessage, msgType)

			// Parse and verify message content
			var ackMsg AckMessage
			err = json.Unmarshal(data, &ackMsg)
			require.NoError(t, err)

			assert.Equal(t, tt.wantType, ackMsg.Type)
			assert.Equal(t, tt.requestID, ackMsg.RequestID)
			assert.Equal(t, tt.status, ackMsg.Status)
		})
	}
}

func TestSendProgress(t *testing.T) {
	tests := []struct {
		name      string
		requestID uuid.UUID
		stage     domain.ProcessingStage
		message   string
	}{
		{
			name:      "categorizing stage",
			requestID: uuid.New(),
			stage:     domain.StageCategorizing,
			message:   "Analyzing request type...",
		},
		{
			name:      "dispatching stage",
			requestID: uuid.New(),
			stage:     domain.StageDispatching,
			message:   "Category determined: easy",
		},
		{
			name:      "processing stage",
			requestID: uuid.New(),
			stage:     domain.StageProcessing,
			message:   "Requesting response from model...",
		},
		{
			name:      "cancelling stage",
			requestID: uuid.New(),
			stage:     domain.StageCancelling,
			message:   "Cancellation requested",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockConn := &mockWebSocketConn{}
			svc := &service{conn: mockConn}

			err := svc.SendProgress(tt.requestID, tt.stage, tt.message)
			require.NoError(t, err)

			// Verify message was written
			require.Len(t, mockConn.writtenMessages, 1)

			msgType, data := mockConn.getLastMessage()
			assert.Equal(t, websocket.TextMessage, msgType)

			// Parse and verify message content
			var progressMsg ProgressMessage
			err = json.Unmarshal(data, &progressMsg)
			require.NoError(t, err)

			assert.Equal(t, domain.MessageTypeProgress, progressMsg.Type)
			assert.Equal(t, tt.requestID, progressMsg.RequestID)
			assert.Equal(t, tt.stage.String(), progressMsg.Stage)
			assert.Equal(t, tt.message, progressMsg.Message)
		})
	}
}

func TestSendToken(t *testing.T) {
	tests := []struct {
		name         string
		requestID    uuid.UUID
		token        string
		isFinalToken bool
	}{
		{
			name:         "regular token",
			requestID:    uuid.New(),
			token:        "Hello",
			isFinalToken: false,
		},
		{
			name:         "final token",
			requestID:    uuid.New(),
			token:        "!",
			isFinalToken: true,
		},
		{
			name:         "empty token",
			requestID:    uuid.New(),
			token:        "",
			isFinalToken: false,
		},
		{
			name:         "token with special characters",
			requestID:    uuid.New(),
			token:        "Hello\n\tWorld!",
			isFinalToken: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockConn := &mockWebSocketConn{}
			svc := &service{conn: mockConn}

			err := svc.SendToken(tt.requestID, tt.token, tt.isFinalToken)
			require.NoError(t, err)

			// Verify message was written
			require.Len(t, mockConn.writtenMessages, 1)

			msgType, data := mockConn.getLastMessage()
			assert.Equal(t, websocket.TextMessage, msgType)

			// Parse and verify message content
			var tokenMsg TokenMessage
			err = json.Unmarshal(data, &tokenMsg)
			require.NoError(t, err)

			assert.Equal(t, domain.MessageTypeToken, tokenMsg.Type)
			assert.Equal(t, tt.requestID, tokenMsg.RequestID)
			assert.Equal(t, tt.token, tokenMsg.Token)
			assert.Equal(t, tt.isFinalToken, tokenMsg.IsFinalToken)
		})
	}
}

func TestSendFinalResponse(t *testing.T) {
	tests := []struct {
		name         string
		requestID    uuid.UUID
		responseText string
		modelUsed    string
	}{
		{
			name:         "simple response",
			requestID:    uuid.New(),
			responseText: "Hello, world!",
			modelUsed:    "gpt-4",
		},
		{
			name:         "long response",
			requestID:    uuid.New(),
			responseText: "This is a much longer response that contains multiple sentences and lots of information.",
			modelUsed:    "claude-3",
		},
		{
			name:         "response with special characters",
			requestID:    uuid.New(),
			responseText: "Response with\nnewlines\tand\ttabs",
			modelUsed:    "llama-2",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockConn := &mockWebSocketConn{}
			svc := &service{conn: mockConn}

			err := svc.SendFinalResponse(tt.requestID, tt.responseText, tt.modelUsed)
			require.NoError(t, err)

			// Verify message was written
			require.Len(t, mockConn.writtenMessages, 1)

			msgType, data := mockConn.getLastMessage()
			assert.Equal(t, websocket.TextMessage, msgType)

			// Parse and verify message content
			var finalMsg FinalResponseMessage
			err = json.Unmarshal(data, &finalMsg)
			require.NoError(t, err)

			assert.Equal(t, domain.MessageTypeFinalResponse, finalMsg.Type)
			assert.Equal(t, tt.requestID, finalMsg.RequestID)
			assert.Equal(t, tt.responseText, finalMsg.ResponseText)
			assert.Equal(t, tt.modelUsed, finalMsg.ModelUsed)
		})
	}
}

func TestSendError(t *testing.T) {
	tests := []struct {
		name      string
		requestID uuid.UUID
		errorCode domain.ErrorCode
		message   string
	}{
		{
			name:      "bad request error",
			requestID: uuid.New(),
			errorCode: domain.ErrorCodeBadRequest,
			message:   "Invalid input",
		},
		{
			name:      "internal error",
			requestID: uuid.New(),
			errorCode: domain.ErrorCodeInternalError,
			message:   "Something went wrong",
		},
		{
			name:      "model error",
			requestID: uuid.New(),
			errorCode: domain.ErrorCodeModelError,
			message:   "Model failed to respond",
		},
		{
			name:      "rate limited error",
			requestID: uuid.New(),
			errorCode: domain.ErrorCodeRateLimited,
			message:   "Too many requests",
		},
		{
			name:      "unauthorized error",
			requestID: uuid.New(),
			errorCode: domain.ErrorCodeUnauthorized,
			message:   "Invalid API key",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockConn := &mockWebSocketConn{}
			svc := &service{conn: mockConn}

			err := svc.SendError(tt.requestID, tt.errorCode, tt.message)
			require.NoError(t, err)

			// Verify message was written
			require.Len(t, mockConn.writtenMessages, 1)

			msgType, data := mockConn.getLastMessage()
			assert.Equal(t, websocket.TextMessage, msgType)

			// Parse and verify message content
			var errorMsg ErrorMessage
			err = json.Unmarshal(data, &errorMsg)
			require.NoError(t, err)

			assert.Equal(t, domain.MessageTypeError, errorMsg.Type)
			assert.Equal(t, tt.requestID, errorMsg.RequestID)
			assert.Equal(t, tt.errorCode, errorMsg.ErrorCode)
			assert.Equal(t, tt.message, errorMsg.Message)
		})
	}
}

func TestConcurrentSends(t *testing.T) {
	mockConn := &mockWebSocketConn{}
	svc := &service{conn: mockConn}

	// Send multiple messages concurrently
	done := make(chan bool)
	numGoroutines := 10

	for i := 0; i < numGoroutines; i++ {
		go func(idx int) {
			requestID := uuid.New()
			_ = svc.SendProgress(requestID, domain.StageProcessing, "Concurrent test")
			done <- true
		}(i)
	}

	// Wait for all goroutines to complete
	for i := 0; i < numGoroutines; i++ {
		select {
		case <-done:
		case <-time.After(5 * time.Second):
			t.Fatal("Timeout waiting for concurrent sends")
		}
	}

	// Verify all messages were written
	assert.Len(t, mockConn.writtenMessages, numGoroutines)
}

func TestWriteError(t *testing.T) {
	mockConn := &mockWebSocketConn{
		writeError: assert.AnError,
	}
	svc := &service{conn: mockConn}

	requestID := uuid.New()

	t.Run("SendAck returns error", func(t *testing.T) {
		err := svc.SendAck(requestID, domain.AckStatusAccepted)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "failed to write message")
	})

	t.Run("SendProgress returns error", func(t *testing.T) {
		err := svc.SendProgress(requestID, domain.StageProcessing, "test")
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "failed to write message")
	})

	t.Run("SendToken returns error", func(t *testing.T) {
		err := svc.SendToken(requestID, "token", false)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "failed to write message")
	})

	t.Run("SendFinalResponse returns error", func(t *testing.T) {
		err := svc.SendFinalResponse(requestID, "response", "model")
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "failed to write message")
	})

	t.Run("SendError returns error", func(t *testing.T) {
		err := svc.SendError(requestID, domain.ErrorCodeInternalError, "error")
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "failed to write message")
	})
}

func TestMessageSequence(t *testing.T) {
	mockConn := &mockWebSocketConn{}
	svc := &service{conn: mockConn}
	requestID := uuid.New()

	// Simulate a typical message sequence
	err := svc.SendAck(requestID, domain.AckStatusAccepted)
	require.NoError(t, err)

	err = svc.SendProgress(requestID, domain.StageCategorizing, "Analyzing...")
	require.NoError(t, err)

	err = svc.SendProgress(requestID, domain.StageDispatching, "Category: easy")
	require.NoError(t, err)

	err = svc.SendProgress(requestID, domain.StageProcessing, "Requesting response...")
	require.NoError(t, err)

	err = svc.SendToken(requestID, "Hello", false)
	require.NoError(t, err)

	err = svc.SendToken(requestID, " world", false)
	require.NoError(t, err)

	err = svc.SendToken(requestID, "!", true)
	require.NoError(t, err)

	err = svc.SendFinalResponse(requestID, "Hello world!", "gpt-4")
	require.NoError(t, err)

	// Verify we have 8 messages in the correct order
	messages := mockConn.getAllMessages()
	require.Len(t, messages, 8)

	// Verify first message is ack
	var ackMsg AckMessage
	err = json.Unmarshal(messages[0].data, &ackMsg)
	require.NoError(t, err)
	assert.Equal(t, domain.MessageTypeAck, ackMsg.Type)

	// Verify last message is final response
	var finalMsg FinalResponseMessage
	err = json.Unmarshal(messages[7].data, &finalMsg)
	require.NoError(t, err)
	assert.Equal(t, domain.MessageTypeFinalResponse, finalMsg.Type)
	assert.Equal(t, "Hello world!", finalMsg.ResponseText)
}

func TestNewService(t *testing.T) {
	// This would require a real websocket connection or more sophisticated mocking
	// For now, we just verify the function exists and returns a non-nil service
	// In a real scenario, you'd use httptest and gorilla/websocket test utilities
	t.Run("NewService creates service", func(t *testing.T) {
		// We can't easily test this without a real websocket connection
		// Just verify the function signature is correct by attempting to call it
		var conn *websocket.Conn // nil is fine for this test
		svc := NewService(conn)
		assert.NotNil(t, svc)
	})
}
