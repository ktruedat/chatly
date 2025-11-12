# Chatly Backend

A WebSocket-based AI chat backend that intelligently routes requests to different AI models based on input type and complexity.

## Features

- **WebSocket Communication**: Real-time bidirectional communication
- **Smart Routing**: Automatic categorization of requests (easy/advanced/coding)
- **Multi-Model Support**: Routes to optimal AI models via OpenRouter
- **Streaming Responses**: Token-by-token response streaming
- **Image Support**: Vision models for image analysis
- **Clean Architecture**: Domain-driven design with layered architecture

## Architecture

The project follows a clean architecture pattern:

```
backend/
├── cmd/app/                    # Application entry point
├── internal/
│   ├── domain/                # Domain entities and business rules
│   │   ├── chat.go           # Chat domain models
│   │   └── errors.go         # Domain errors
│   └── application/          # Application layer
│       ├── config/           # Configuration management
│       ├── handlers/         # Request handlers
│       │   └── websocket/    # WebSocket handler
│       └── services/         # Business logic services
│           ├── chat/         # Chat service
│           ├── openrouter/   # OpenRouter API client
│           └── websocket/    # WebSocket message service
└── pkg/                      # Shared packages
    ├── errors/              # Error handling utilities
    └── optional/            # Optional type implementation
```

## Setup

### Prerequisites

- Go 1.21 or higher
- OpenRouter API key ([Get one here](https://openrouter.ai/))

### Installation

1. Clone the repository:
```bash
cd /Users/ktruedat/Projects/work/chatly/backend
```

2. Install dependencies:
```bash
go mod download
```

3. Create `.env` file:
```bash
cp .env.example .env
```

4. Edit `.env` and add your OpenRouter API key:
```env
OPENROUTER_API_KEY=sk-or-v1-your-api-key-here
```

### Running the Server

```bash
go run cmd/app/main.go
```

The server will start on `http://localhost:3000` by default.

## Configuration

The server is configured via `config.yaml` and environment variables:

### Environment Variables

- `OPENROUTER_API_KEY` (required): Your OpenRouter API key
- `OPENROUTER_BASE_URL` (optional): Override OpenRouter base URL
- `SERVER_PORT` (optional): Override server port
- `CONFIG_PATH` (optional): Path to config file

### Config File

See `config.yaml` for full configuration options:

- Server settings (port, timeouts)
- OpenRouter configuration
- Model selection
- Categorization prompt

## API Documentation

See [WEBSOCKET_API.md](../WEBSOCKET_API.md) for complete WebSocket API documentation.

### Quick Example

Connect to WebSocket:
```javascript
const ws = new WebSocket('ws://localhost:3000/ws');

ws.onopen = () => {
  // Send a text request
  ws.send(JSON.stringify({
    type: 'submit_request',
    request_id: crypto.randomUUID(),
    input_type: 'Text',
    content: 'What is 2 + 2?',
    image_base64: null,
    metadata: null
  }));
};

ws.onmessage = (event) => {
  const message = JSON.parse(event.data);
  console.log('Received:', message);
};
```

## Development

### Project Structure

The project follows clean architecture principles:

- **Domain Layer**: Business entities and rules (`internal/domain/`)
- **Application Layer**: Use cases and services (`internal/application/`)
- **Handlers Layer**: Input/output handling (`internal/application/handlers/`)
- **Infrastructure**: External services (`pkg/`)

### Adding a New Model

1. Add model ID to `config.yaml`:
```yaml
models:
  new_category: "provider/model-name:free"
```

2. Add category to `internal/domain/chat.go`:
```go
const (
    CategoryNewCategory Category = "new_category"
)
```

3. Update routing logic in `internal/application/services/chat/service.go`

### Error Handling

The application uses a structured error system:

```go
// Domain errors (internal/domain/errors.go)
var ErrContentRequired = errors.New("content is required")

// Application errors (internal/application/common/errors/)
NewBadRequestError("Invalid input")
NewModelError("API failed", cause)
NewInternalError(err)
```

## Testing

### Manual Testing

Use the WebSocket test client:

```bash
# Install wscat
npm install -g wscat

# Connect and test
wscat -c ws://localhost:3000/ws

# Send a request
{"type":"submit_request","request_id":"550e8400-e29b-41d4-a716-446655440000","input_type":"Text","content":"What is 2+2?","image_base64":null}
```

## Deployment

### Docker (Coming Soon)

```bash
docker build -t chatly-backend .
docker run -p 3000:3000 -e OPENROUTER_API_KEY=your-key chatly-backend
```

### Production Considerations

1. **WebSocket Security**: Use WSS (WebSocket Secure) in production
2. **CORS**: Update CORS settings in `application.go`
3. **Rate Limiting**: Implement rate limiting for API protection
4. **Authentication**: Add authentication/authorization
5. **Monitoring**: Add logging and metrics
6. **Origin Checking**: Update `upgrader.CheckOrigin` in websocket handler

## Troubleshooting

### Common Issues

**WebSocket connection fails**
- Check if server is running: `curl http://localhost:3000/health`
- Verify port is not in use: `lsof -i :3000`

**API errors**
- Verify OPENROUTER_API_KEY is set correctly
- Check OpenRouter quota and rate limits
- Review server logs for detailed error messages

**Categorization not working**
- Ensure categorizer model is accessible
- Check categorization prompt in config.yaml
- Model may return unexpected format

## Contributing

This is a university project. For questions or issues, please contact the development team.

## License

University Project - Not for commercial use

## Credits

- Architecture inspired by admin-api project
- AI models provided by [OpenRouter](https://openrouter.ai/)
- Built with Go, Chi, and Gorilla WebSocket
