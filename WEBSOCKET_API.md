# Chatly WebSocket API Documentation

## Overview
The Chatly backend provides a WebSocket API for real-time AI-powered chat functionality. It intelligently routes requests to different AI models based on input type and complexity.

## Connection

### WebSocket Endpoint
```
ws://localhost:3000/ws
```

### Production Endpoint
```
wss://your-domain.com/ws
```

## Message Formats

All messages are JSON formatted. The protocol uses tagged enums for type safety.

---

## Client → Server Messages

### 1. Submit Request

Send a new request to be processed by the AI.

```json
{
  "type": "submit_request",
  "request_id": "550e8400-e29b-41d4-a716-446655440000",
  "input_type": "Text",
  "content": "What is 2 + 2?",
  "image_base64": null,
  "metadata": null
}
```

**Fields:**

- `type`: (string) Always `"submit_request"`
- `request_id`: (string, UUID) Unique identifier for this request. Generate using UUID v4.
- `input_type`: (enum) One of:
  - `"Text"` - Text-only input
  - `"Image"` - Image with optional text (uses standard vision model)
  - `"ImageHard"` - Complex image analysis (uses advanced vision model)
- `content`: (string | null) The text content of the request
  - Required for `Text` input type
  - Optional for `Image` and `ImageHard` (defaults to "Describe this image")
- `image_base64`: (string | null) Base64-encoded image data
  - Required for `Image` and `ImageHard` input types
  - Must be null for `Text` input type
  - Format: Just the base64 string (no data URI prefix needed)
- `metadata`: (object | null) Optional metadata for tracking/analytics

**Example - Text Request:**
```json
{
  "type": "submit_request",
  "request_id": "123e4567-e89b-12d3-a456-426614174000",
  "input_type": "Text",
  "content": "How do I use async/await in Rust?",
  "image_base64": null,
  "metadata": {
    "user_id": "user_123",
    "session_id": "session_456"
  }
}
```

**Example - Image Request:**
```json
{
  "type": "submit_request",
  "request_id": "987e6543-e21b-12d3-a456-426614174001",
  "input_type": "Image",
  "content": "What's in this image?",
  "image_base64": "/9j/4AAQSkZJRgABAQAA...",
  "metadata": null
}
```

---

### 2. Cancel Request

Cancel an in-progress request (not yet fully implemented).

```json
{
  "type": "cancel_request",
  "request_id": "550e8400-e29b-41d4-a716-446655440000"
}
```

**Fields:**

- `type`: (string) Always `"cancel_request"`
- `request_id`: (string, UUID) The request ID to cancel

---

## Server → Client Messages

### 1. Acknowledgment

Immediate confirmation that the request was received and accepted for processing.

```json
{
  "type": "ack",
  "request_id": "550e8400-e29b-41d4-a716-446655440000",
  "status": "accepted"
}
```

**Fields:**

- `type`: (string) Always `"ack"`
- `request_id`: (string, UUID) The request ID being acknowledged
- `status`: (string) Status of acknowledgment, typically `"accepted"`

---

### 2. Progress

Status updates during request processing.

```json
{
  "type": "progress",
  "request_id": "550e8400-e29b-41d4-a716-446655440000",
  "stage": "categorizing",
  "message": "Analyzing request type..."
}
```

**Fields:**

- `type`: (string) Always `"progress"`
- `request_id`: (string, UUID) The request ID
- `stage`: (string) Current processing stage:
  - `"categorizing"` - Determining request complexity (text-only)
  - `"dispatching"` - Selecting appropriate AI model
  - `"processing"` - Waiting for AI model response
  - `"cancelling"` - Processing cancellation
- `message`: (string) Human-readable status message

**Progress Flow for Text Requests:**
1. `categorizing` - "Analyzing request type..."
2. `dispatching` - "Category determined: easy/advanced/coding"
3. `processing` - "Requesting response from [model name]..."

**Progress Flow for Image Requests:**
1. `dispatching` - "Category determined: image"
2. `processing` - "Requesting response from [model name]..."

---

### 3. Token (Streaming Response)

Real-time streaming of the AI's response, token by token.

```json
{
  "type": "token",
  "request_id": "550e8400-e29b-41d4-a716-446655440000",
  "token": "The ",
  "is_final_token": false
}
```

**Fields:**

- `type`: (string) Always `"token"`
- `request_id`: (string, UUID) The request ID
- `token`: (string) A piece of the response text
  - May be a single character, word, or phrase
  - Tokens should be concatenated in order to build the full response
- `is_final_token`: (boolean) 
  - `false` - More tokens will follow
  - `true` - This is the last token (rare, usually non-streaming responses)

**Frontend Implementation:**
```javascript
let responseText = "";

// On each token message:
responseText += token;
updateUIWithPartialResponse(responseText);
```

---

### 4. Final Response

The complete response after all tokens have been sent.

```json
{
  "type": "final_response",
  "request_id": "550e8400-e29b-41d4-a716-446655440000",
  "response_text": "The answer is 4.",
  "model_used": "google/gemma-3n-e2b-it:free"
}
```

**Fields:**

- `type`: (string) Always `"final_response"`
- `request_id`: (string, UUID) The request ID
- `response_text`: (string) The complete AI response
  - This is the same as concatenating all tokens
  - Provided for convenience and verification
- `model_used`: (string) The AI model that generated the response
  - Useful for analytics and debugging

**Model Names by Category:**

| Input Type | Category | Model |
|-----------|----------|-------|
| Text | Easy | `google/gemma-3n-e2b-it:free` |
| Text | Advanced | `deepseek/deepseek-chat-v3.1:free` |
| Text | Coding | `qwen/qwen3-coder:free` |
| Image | Standard | `qwen/qwen2.5-vl-32b-instruct:free` |
| ImageHard | Complex | `nvidia/nemotron-nano-12b-v2-vl:free` |

---

### 5. Error

Error messages when something goes wrong.

```json
{
  "type": "error",
  "request_id": "550e8400-e29b-41d4-a716-446655440000",
  "error_code": "model_error",
  "message": "API error 429: Rate limit exceeded"
}
```

**Fields:**

- `type`: (string) Always `"error"`
- `request_id`: (string, UUID) The request ID
- `error_code`: (string) Machine-readable error code
  - `"bad_request"` - Invalid message format
  - `"model_error"` - Error from AI model API
  - `"internal_error"` - Server-side error
- `message`: (string) Human-readable error description

---

## Message Flow Examples

### Example 1: Simple Text Question

**Client sends:**
```json
{
  "type": "submit_request",
  "request_id": "req-001",
  "input_type": "Text",
  "content": "What is 2 + 2?",
  "image_base64": null,
  "metadata": null
}
```

**Server responds (in sequence):**

1. Acknowledgment:
```json
{"type": "ack", "request_id": "req-001", "status": "accepted"}
```

2. Progress - Categorizing:
```json
{"type": "progress", "request_id": "req-001", "stage": "categorizing", "message": "Analyzing request type..."}
```

3. Progress - Dispatching:
```json
{"type": "progress", "request_id": "req-001", "stage": "dispatching", "message": "Category determined: easy"}
```

4. Progress - Processing:
```json
{"type": "progress", "request_id": "req-001", "stage": "processing", "message": "Requesting response from google/gemma-3n-e2b-it:free..."}
```

5-N. Tokens (streaming):
```json
{"type": "token", "request_id": "req-001", "token": "2", "is_final_token": false}
{"type": "token", "request_id": "req-001", "token": " +", "is_final_token": false}
{"type": "token", "request_id": "req-001", "token": " ", "is_final_token": false}
{"type": "token", "request_id": "req-001", "token": "2", "is_final_token": false}
{"type": "token", "request_id": "req-001", "token": " =", "is_final_token": false}
{"type": "token", "request_id": "req-001", "token": " ", "is_final_token": false}
{"type": "token", "request_id": "req-001", "token": "4", "is_final_token": false}
```

N+1. Final Response:
```json
{
  "type": "final_response",
  "request_id": "req-001",
  "response_text": "2 + 2 = 4",
  "model_used": "google/gemma-3n-e2b-it:free"
}
```

---

### Example 2: Coding Question

**Client sends:**
```json
{
  "type": "submit_request",
  "request_id": "req-002",
  "input_type": "Text",
  "content": "How do I iterate over a vector in Rust?",
  "image_base64": null,
  "metadata": null
}
```

**Server responds:**
- Ack → Progress (categorizing) → Progress (dispatching: "coding") → Progress (processing) → Tokens → Final Response
- `model_used` will be `"qwen/qwen3-coder:free"`

---

### Example 3: Image Analysis

**Client sends:**
```json
{
  "type": "submit_request",
  "request_id": "req-003",
  "input_type": "Image",
  "content": "What's in this image?",
  "image_base64": "/9j/4AAQSkZJRg...",
  "metadata": null
}
```

**Server responds:**
- Ack → Progress (dispatching: "image") → Progress (processing) → Tokens → Final Response
- `model_used` will be `"qwen/qwen2.5-vl-32b-instruct:free"`
- Note: No categorization stage for image requests

---

## Frontend Implementation Guide

### TypeScript Types

```typescript
// Request types
type InputType = "Text" | "Image" | "ImageHard";

interface SubmitRequest {
  type: "submit_request";
  request_id: string;
  input_type: InputType;
  content: string | null;
  image_base64: string | null;
  metadata?: Record<string, any> | null;
}

interface CancelRequest {
  type: "cancel_request";
  request_id: string;
}

type ClientMessage = SubmitRequest | CancelRequest;

// Response types
interface AckMessage {
  type: "ack";
  request_id: string;
  status: string;
}

interface ProgressMessage {
  type: "progress";
  request_id: string;
  stage: string;
  message: string;
}

interface TokenMessage {
  type: "token";
  request_id: string;
  token: string;
  is_final_token: boolean;
}

interface FinalResponseMessage {
  type: "final_response";
  request_id: string;
  response_text: string;
  model_used: string;
}

interface ErrorMessage {
  type: "error";
  request_id: string;
  error_code: string;
  message: string;
}

type ServerMessage = 
  | AckMessage 
  | ProgressMessage 
  | TokenMessage 
  | FinalResponseMessage 
  | ErrorMessage;
```

### JavaScript Example

```javascript
class ChatlyClient {
  constructor(url = "ws://localhost:3000/ws") {
    this.url = url;
    this.ws = null;
    this.handlers = new Map();
  }

  connect() {
    return new Promise((resolve, reject) => {
      this.ws = new WebSocket(this.url);
      
      this.ws.onopen = () => resolve();
      this.ws.onerror = (error) => reject(error);
      this.ws.onmessage = (event) => this.handleMessage(event);
    });
  }

  handleMessage(event) {
    const message = JSON.parse(event.data);
    const handler = this.handlers.get(message.request_id);
    
    if (handler) {
      handler(message);
      
      // Clean up handler after final response or error
      if (message.type === "final_response" || message.type === "error") {
        this.handlers.delete(message.request_id);
      }
    }
  }

  sendRequest(inputType, content, imageBase64 = null, metadata = null) {
    const requestId = crypto.randomUUID();
    
    const request = {
      type: "submit_request",
      request_id: requestId,
      input_type: inputType,
      content,
      image_base64: imageBase64,
      metadata
    };
    
    this.ws.send(JSON.stringify(request));
    
    return {
      requestId,
      onMessage: (handler) => {
        this.handlers.set(requestId, handler);
      }
    };
  }

  cancelRequest(requestId) {
    const cancel = {
      type: "cancel_request",
      request_id: requestId
    };
    
    this.ws.send(JSON.stringify(cancel));
  }
}

// Usage example
const client = new ChatlyClient();

await client.connect();

const { requestId, onMessage } = client.sendRequest(
  "Text",
  "What is 2 + 2?",
  null,
  { user_id: "user_123" }
);

let responseText = "";

onMessage((message) => {
  switch (message.type) {
    case "ack":
      console.log("Request accepted");
      showLoadingAnimation();
      break;
      
    case "progress":
      console.log(`Progress: ${message.stage} - ${message.message}`);
      updateStatusText(message.message);
      break;
      
    case "token":
      responseText += message.token;
      updateResponseUI(responseText);
      break;
      
    case "final_response":
      console.log("Complete response:", message.response_text);
      console.log("Model used:", message.model_used);
      hideLoadingAnimation();
      markResponseComplete();
      break;
      
    case "error":
      console.error(`Error ${message.error_code}: ${message.message}`);
      showErrorUI(message.message);
      hideLoadingAnimation();
      break;
  }
});
```

### React Example

```tsx
import { useState, useEffect, useRef } from 'react';

function useChatly(url = "ws://localhost:3000/ws") {
  const [isConnected, setIsConnected] = useState(false);
  const [messages, setMessages] = useState<Map<string, any[]>>(new Map());
  const wsRef = useRef<WebSocket | null>(null);

  useEffect(() => {
    const ws = new WebSocket(url);
    
    ws.onopen = () => setIsConnected(true);
    ws.onclose = () => setIsConnected(false);
    
    ws.onmessage = (event) => {
      const message = JSON.parse(event.data);
      
      setMessages(prev => {
        const newMap = new Map(prev);
        const reqMessages = newMap.get(message.request_id) || [];
        newMap.set(message.request_id, [...reqMessages, message]);
        return newMap;
      });
    };
    
    wsRef.current = ws;
    
    return () => ws.close();
  }, [url]);

  const sendRequest = (inputType: string, content: string, imageBase64?: string) => {
    if (!wsRef.current || !isConnected) {
      throw new Error("WebSocket not connected");
    }

    const requestId = crypto.randomUUID();
    
    wsRef.current.send(JSON.stringify({
      type: "submit_request",
      request_id: requestId,
      input_type: inputType,
      content,
      image_base64: imageBase64 || null,
      metadata: null
    }));
    
    return requestId;
  };

  return { isConnected, messages, sendRequest };
}

// Component example
function ChatInterface() {
  const { isConnected, messages, sendRequest } = useChatly();
  const [input, setInput] = useState("");
  const [currentRequestId, setCurrentRequestId] = useState<string | null>(null);

  const handleSend = () => {
    const requestId = sendRequest("Text", input);
    setCurrentRequestId(requestId);
    setInput("");
  };

  const currentMessages = currentRequestId 
    ? messages.get(currentRequestId) || []
    : [];

  const responseTokens = currentMessages
    .filter(m => m.type === "token")
    .map(m => m.token)
    .join("");

  const isProcessing = currentMessages.some(
    m => m.type === "progress" || m.type === "token"
  );

  return (
    <div>
      <div className="chat-display">
        {responseTokens && <p>{responseTokens}</p>}
        {isProcessing && <LoadingIndicator />}
      </div>
      
      <input
        value={input}
        onChange={(e) => setInput(e.target.value)}
        disabled={!isConnected || isProcessing}
      />
      
      <button 
        onClick={handleSend}
        disabled={!isConnected || isProcessing || !input.trim()}
      >
        Send
      </button>
    </div>
  );
}
```

---

## Best Practices

### 1. Request ID Management
- Always generate unique UUIDs for each request
- Track request IDs to match responses
- Clean up completed request handlers to prevent memory leaks

### 2. Streaming Responses
- Concatenate tokens in order as they arrive
- Update UI incrementally for better UX
- Wait for `final_response` before marking completion

### 3. Error Handling
- Always handle `error` messages
- Show user-friendly error messages
- Log error codes for debugging
- Implement retry logic for transient errors

### 4. Connection Management
- Implement reconnection logic
- Handle connection lost scenarios
- Queue requests during disconnection
- Notify users of connection status

### 5. Image Handling
- Convert images to base64 before sending
- Remove data URI prefix if present
- Consider image size limits (OpenRouter typically allows up to 20MB)
- Show upload progress for large images

### 6. Performance
- Don't block UI during token streaming
- Use requestAnimationFrame for smooth updates
- Debounce rapid token updates if needed
- Clean up old messages to prevent memory issues

---

## Testing

The backend includes comprehensive integration tests. To run them:

```bash
# Start the server
OPENROUTER_API_KEY=your_key cargo run

# In another terminal, run tests
cargo test --test integration_test
```

---

## Rate Limits & Quotas

OpenRouter free tier has rate limits:
- **Requests per minute**: Varies by model
- **Tokens per request**: Model-dependent
- **Error code 429**: Rate limit exceeded

Handle rate limits gracefully:
- Implement exponential backoff
- Show friendly messages to users
- Consider caching common responses

---

## Security Considerations

### WebSocket Security
- Use WSS (WebSocket Secure) in production
- Implement authentication/authorization
- Validate all input data
- Sanitize user content before display

### Image Upload
- Validate image format and size
- Scan for malicious content
- Consider implementing virus scanning

---

## Environment Variables

Required environment variables:

```bash
OPENROUTER_API_KEY=sk-or-v1-...
OPENROUTER_BASE_URL=https://api.openrouter.ai  # Optional, defaults to this
```

---

## Support & Questions

For issues or questions:
1. Check the integration tests for working examples
2. Review the AI_models.md file for model specifications
3. Test with the provided examples first
4. Check server logs for detailed error messages

---

## Changelog

### v0.1.0 (Initial Release)
- WebSocket-based real-time communication
- Automatic request categorization (easy/advanced/coding)
- Multi-model AI routing
- Token-by-token streaming responses
- Image analysis support
- Comprehensive error handling
