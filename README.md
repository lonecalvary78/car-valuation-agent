# Car Valuation Agent

[![car-valuer-agent-ci](https://github.com/lonecalvary78/car-valuation-agent/actions/workflows/ci.yaml/badge.svg)](https://github.com/lonecalvary78/car-valuation-agent/actions/workflows/ci.yaml)

An HTTP service backed by an LLM agent (Google ADK) with any openweight model backend that estimates the market value of a used car based on:
- **Brand**: BMW, Volvo, Toyota, Honda, etc.
- **Model** 
- **Year**: the year the car was manufactured
- **Market**: US, Germany, Japan, etc.

The agent is instructed to answer using a `car-valuer` skill (see [skills/car-valuer/SKILL.md](skills/car-valuer/SKILL.md)) rather than its own background knowledge, and replies with a structured JSON price range.

## Trade-offs
- Both memory and session are still using in-memory
- No setup with TLS
- No support for Gemini model support
- No user authentication

## Stack
- **Language**: Go
- **Libraries**: `net/http`, [`adk-go`](https://google.golang.org/adk/v2), [`sonic`](https://github.com/bytedance/sonic) (JSON), [`decimal`](https://github.com/shopspring/decimal) (money)

## Configuration

Config is loaded from environment variables:

| Variable | Description | Default |
|---|---|---|
| `SERVER_HOST` | Address the HTTP server binds to | `localhost` or `0.0.0.0` |
| `SERVER_PORT` | Port the HTTP server listens on | `8080` |
| `SERVER_READ_TIME` | Request read timeout (e.g. `30s`) | `30s` |
| `SERVER_WRITE_TIME` | Response write timeout (e.g. `2m`) | `2m` |
| `AGENT_NAME` | Name registered with the ADK runner/session service | — |
| `BASE_URL` | OpenAI-compatible base URL for the model provider (e.g. OpenRouter) | — |
| `API_KEY` | API key for the model provider | — |
| `MODEL_NAME` | Model identifier passed to the provider | — |
| `SKILL_LOCATION` | Filesystem path containing skill definitions | — |
| `WAIT_TIMEOUT` | Wait for the response from the agent(e.g. `60s` ) | `60s` |


> Binding `SERVER_HOST` to `localhost` only accepts connections from inside the container/host itself — leave it empty (or use `0.0.0.0`) if the server needs to be reachable through Docker's port mapping.

## Running locally

```bash
go build ./...              # Build all packages
go run ./cmd/api             # Run the API server
go test ./...                 # Run all tests
go vet ./...                 # Static analysis
```

## API

| Method | Path | Description |
|---|---|---|
| `POST` | `/chat` | Send `{"userId", "sessionId", "message"}`; returns the agent's reply with JSON format. `sessionId` is generated automatically if omitted. |
| `GET` | `/health` | Returns `{"status", "checked_at"}`. |

Example:

```bash
curl -X POST localhost:8080/chat \
  -H 'Content-Type: application/json' \
  -d '{"userId": "u1", "sessionId": "s1", "message": "2019 Toyota Camry in the US"}'
```

## Docker

```bash
docker compose build
docker compose up
```

