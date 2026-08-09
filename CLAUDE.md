# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Commands

```bash
go build ./...              # Build all packages
go run ./cmd/api             # Run the API server
go test ./...                 # Run all tests
go vet ./...                 # Static analysis
```

```bash
docker compose build         # Build the image (see Dockerfile)
docker compose up            # Build + run via docker-compose.yaml
```

Config is loaded purely from process environment variables via `os.Getenv` in `config.Load()` (`internal/infrastructure/config/util.go`). For `go run`, export the variables into your shell first; under Docker, `docker-compose.yaml`'s `env_file: - .env` supplies them instead.

| Variable | Default | Notes |
|---|---|---|
| `SERVER_HOST` | *localhost* | Bind address. `localhost` only accepts loopback connections — leave empty or use `0.0.0.0` so Docker's port mapping can reach it. |
| `SERVER_PORT` | `8080` | |
| `SERVER_READ_TIME` | `30s` | Parsed with `time.ParseDuration`; falls back to default on parse error. |
| `SERVER_WRITE_TIME` | `2m` | Same fallback behavior. |
| `AGENT_NAME` | *(empty)* | Agent's name; also becomes the ADK runner's `AppName` via `targetAgent.Name()`. |
| `BASE_URL` | *(empty)* | OpenAI-compatible base URL for the model provider (e.g. OpenRouter). |
| `API_KEY` | *(empty)* | |
| `MODEL_NAME` | *(empty)* | |
| `SKILL_LOCATION` | *(empty)* | Filesystem path (relative to CWD) passed to `os.DirFS` to load skill definitions. |
| `SKILL_NAME` | *(empty)* | Loaded into config but currently unused — only `Skill.Location` is read. |

## Architecture

This is a Go service built on Google's Agent Development Kit (`google.golang.org/adk/v2`) that exposes a car-valuation chat agent over HTTP. The flow, wired up in `cmd/api/server.go`, is:

1. **Config** (`internal/infrastructure/config`) — `Load()` reads env vars into `AppConfig{server Server, agent Agent{model Model, skill Skill}}`.
2. **Skill tool** (`internal/infrastructure/util/llm.go: OfSkillBasedTool`) — builds an ADK `skilltoolset.SkillToolset` from `skill.NewFileSystemSource(os.DirFS(SKILL_LOCATION))`. Skills are Markdown files with YAML frontmatter; the only one currently defined is `skills/car-valuer/SKILL.md`, which instructs the model to respond with a strict JSON price range (brand/model/year/market/quoted_price/quoted_at) and nothing else.
3. **LLM + Agent** (`internal/infrastructure/util/llm.go: OfAgent`) — always builds an `openaimodel.NewModel` (OpenAI-compatible client pointed at `BASE_URL`/`API_KEY`/`MODEL_NAME`; works against OpenRouter). There is no Gemini/provider-dispatch path. The model and skill toolset are wrapped into an `llmagent.New` with a fixed instruction constant (`agentInstruction` in `cmd/api/server.go`) that tells the agent to defer to the `car-valuer` skill instead of its own knowledge.
4. **Runner** (`internal/infrastructure/util/llm.go: OfRunner`) — wires the agent to `runner.New` with in-memory session/memory services (`memory.InMemoryService()`, `session.InMemoryService()`) and `AutoCreateSession: true`. `AppName` comes from `targetAgent.Name()` (i.e. `AGENT_NAME`), not a separate parameter.
5. **HTTP layer** (`internal/handler/handler.go`) — `RegisterRoutes()` mounts `POST /chat` and `GET /health` on a fresh `http.ServeMux`.
   - `askToAgent` decodes a `model.ChatRequest{sessionId, userId, message}` body via `sonic`, generates a session ID (`util.NewSessionId()`, a `uuid.New()`) when none is supplied, runs the agent via `runner.Run`, concatenates all non-thought text parts from the returned event stream, and writes the result back as **plain text** (not JSON) with a 200 status. A run-loop error is returned via `http.Error` with the raw error string — no generic `{"error": ...}` envelope.
   - `healthCheck` returns `model.HealthCheckResponse{status, checked_at}` as JSON.
   - There is currently no auth/authorization on `userId`/`sessionId` beyond what the request body provides.
6. **Middleware** (`internal/infrastructure/middeware` — package name intentionally matches the misspelled directory) — `Chain(handler, mws...)` composes `Logging` (logs status/method/URI/duration per request) and `Recovery` (recovers panics into a 500) around the mux.

All ADK object construction (skill tool, agent, runner) is centralized in `internal/infrastructure/util`, keeping `cmd/api/server.go` a straight-line composition root.

JSON encoding/decoding uses `github.com/bytedance/sonic` rather than `encoding/json`. There is no monetary-decimal dependency in this version — the old flat-stub `valuation` package has been replaced entirely by the skill-based tool described above.

No test files exist in the repository yet (`find . -name '*_test.go'` is empty).
