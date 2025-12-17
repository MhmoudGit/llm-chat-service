# Go LLM Chat Service

[![Go Report Card](https://goreportcard.com/badge/github.com/MhmoudGit/llm-chat-service)](https://goreportcard.com/report/github.com/MhmoudGit/llm-chat-service)
![Go Version](https://img.shields.io/github/go-mod/go-version/MhmoudGit/llm-chat-service)
[![CI](https://github.com/MhmoudGit/llm-chat-service/actions/workflows/ci.yml/badge.svg)](https://github.com/MhmoudGit/llm-chat-service/actions/workflows/ci.yml)

A lightweight, concurrent Go web service powering a single-conversation chat interface with Groq Cloud LLMs. Supported by Server-Sent Events (SSE) for real-time token streaming.

[This Repo Contains code generated using Antigravity AI with Gemini 3 Pro Model]

## Web Interface

A simple, vanilla HTML/CSS/JS web interface is available to interact with the chat service.

- **URL**: `http://localhost:8080/web`
- **Features**: Real-time chat with streaming responses.

## Setup & Run

### Prerequisites
- Go 1.25+
- Docker & Docker Compose
- Groq API Key ([Get one here](https://console.groq.com/))

### Running Locally
1.  **Set Environment Variables**:
    ```bash
    export GROQ_API_KEY="your_api_key"
    export PORT=8080
    export API_KEY="your_secret_key"
    export RATE_LIMIT_RPS=10
    export RATE_LIMIT_BURST=20
    ```
    Or Just use the '.env.example' file in the root directory.

2.  **Run**:
    ```bash
    go run cmd/server/main.go
    ```

### Running with Docker
- Run after setting the environment variables in the '.env' file:
```bash
docker-compose up --build
```
The service will be available at `http://localhost:8080`.

## Security

### API Key Authentication
The service is protected by API Key authentication.
- **Header**: `Authorization: Bearer <your_api_key>` or `X-API-Key: <your_api_key>`
- **Configuration**: Set `API_KEY` environment variable.

### Rate Limiting
Requests are rate-limited per IP address.
- **Default**: 10 requests per second with a burst of 20.
- **Configuration**: Set `RATE_LIMIT_RPS` and `RATE_LIMIT_BURST`.

## API Contract

### 1. Health Check
- **Endpoint**: `GET /health`
- **Response**: `200 OK`
- **Body**: `OK`

### 2. Chat Completion
- **Endpoint**: `POST /chat`

- **Headers**: 
    - `Content-Type: application/json`
    - `Authorization: Bearer <your_api_key>`
- **Body**:
    ```json
    {
      "messages": [
        {"role": "user", "content": "Hello, world!"}
      ],
      "stream": true
    }
    ```
- **Response**: Server-Sent Events (SSE) stream.
    - Event: `data: {"content":"Hello"}`
    - ...
    - End: `data: [DONE]`

#### Sample cURL
```bash
curl -N -X POST http://localhost:8080/chat \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer your_secret_key" \
  -d '{
    "messages": [{"role": "user", "content": "Explain quantum computing briefly"}],
    "stream": true
  }'
```

### 3. Chat History
- **Endpoint**: `GET /history`
- **Headers**: `Authorization: Bearer <your_api_key>`
- **Response**: JSON array of message objects.

## Continuous Integration

This project uses GitHub Actions for CI.

- **Workflow**: `.github/workflows/ci.yml`
- **Triggers**: Push and Pull Request to `main` branch.
- **Jobs**:
    - **Test**: Runs `go test -v ./...`
    - **Lint**: Runs `golangci-lint` to ensure code quality.
## Design Notes

**Architecture**
The project follows a standard Go project layout to ensure maintainability and separation of concerns:
- **`cmd/server`**: Entry point, wiring dependencies (dependency injection).
- **`internal/api`**: HTTP transport layer. responsible for request parsing, middleware (logging, CORS), and SSE streaming logic.
- **`internal/chat`**: Core business domain. Manages conversation history (`HistoryManager`) and orchestrates the LLM interaction (`Service`).
- **`internal/llm`**: Infrastructure adapter for the external Groq API.

**Trade-offs & Decisions**
1.  **In-Memory Persistence**:
    - *Decision*: History is stored in a thread-safe circular buffer (`sync.RWMutex`).
    - *Trade-off*: Extremely fast and simple, but state is ephemeral (lost on restart) and prevents horizontal scaling (stateful). For a production app, I would use Redis or a DB.
2.  **Server-Sent Events (SSE)**:
    - *Decision*: Used SSE over WebSockets.
    - *Trade-off*: SSE is simpler and ideal for unidirectional text generation using standard HTTP. WebSockets are full-duplex but add complexity (ping/pong, connection upgrades) unnecessary for this specific use case.
3.  **Concurrency handling**:
    - *Decision*: Granular locking on the history slice.
    - *Trade-off*: Ensures safety for concurrent requests (req handled via goroutines).
