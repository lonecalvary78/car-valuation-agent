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
| `SERVER_HOST` | `0.0.0.0` | Bind address by using `0.0.0.0` so Docker's port mapping can reach it. |
| `SERVER_PORT` | `8080` | |
| `SERVER_READ_TIMEOUT` | `30s` | Parsed with `time.ParseDuration`; falls back to default on parse error. |
| `SERVER_WRITE_TIMEOUT` | `2m` | Same fallback behavior. |
| `AGENT_NAME` | *(empty)* | Agent's name; also becomes the ADK runner's `AppName` via `targetAgent.Name()`. |
| `BASE_URL` | *(empty)* | OpenAI-compatible base URL for the model provider (e.g. OpenRouter). |
| `API_KEY` | *(empty)* | |
| `MODEL_NAME` | *(empty)* | |
| `SKILL_LOCATION` | *(empty)* | Filesystem path (relative to CWD) passed to `os.DirFS` to load skill definitions. |
| `SKILL_NAME` | *(empty)* | Loaded into config but currently unused — only `Skill.Location` is read. |
| `WAIT_TIMEOUT` | Wait for the response from the agent(e.g. `60s` ) | `60s` |
| `KEYCLOAK_BASE_URL` | *(empty)* | Base URL of the Keycloak server (e.g. `http://localhost:8180`); combined with the realm to form the OIDC issuer used for provider discovery. |
| `KEYCLOAK_REALM` | *(empty)* | Realm name; issuer is `{KEYCLOAK_BASE_URL}/realms/{KEYCLOAK_REALM}`. |
| `KEYCLOAK_CLIENT_ID` | *(empty)* | Recorded in config for the client this API represents; token verification does not check `aud` against it since Keycloak access tokens are typically issued with `aud=account`. |
| `REDIS_ADDR` | `localhost:6379` | Backing store for the rate-limit middleware. |
| `REDIS_PASSWORD` | *(empty)* | |
| `REDIS_DB` | `0` | |
| `RATE_LIMIT_REQUESTS` | `60` | Max requests per caller per `RATE_LIMIT_WINDOW`. |
| `RATE_LIMIT_WINDOW` | `1m` | Parsed with `time.ParseDuration`. |

## Architecture

This is a Go service built on Google's Agent Development Kit (`google.golang.org/adk/v2`) that exposes a car-valuation chat agent over HTTP. The flow, wired up in `cmd/api/server.go`, is:

1. **Config** (`internal/infrastructure/config`) — `Load()` reads env vars into `AppConfig{server Server, agent Agent{model Model, skill Skill}, keycloak Keycloak}`.
2. **Skill tool** (`internal/infrastructure/util/llm.go: OfSkillBasedTool`) — builds an ADK `skilltoolset.SkillToolset` from `skill.NewFileSystemSource(os.DirFS(SKILL_LOCATION))`. Skills are Markdown files with YAML frontmatter; the only one currently defined is `skills/car-valuer/SKILL.md`, which instructs the model to respond with a strict JSON price range (brand/model/year/market/quoted_price/quoted_at) and nothing else.
3. **LLM + Agent** (`internal/infrastructure/util/llm.go: OfAgent`) — always builds an `openaimodel.NewModel` (OpenAI-compatible client pointed at `BASE_URL`/`API_KEY`/`MODEL_NAME`; works against OpenRouter). There is no Gemini/provider-dispatch path. The model and skill toolset are wrapped into an `llmagent.New` with a fixed instruction constant (`agentInstruction` in `cmd/api/server.go`) that tells the agent to defer to the `car-valuer` skill instead of its own knowledge.
4. **Runner** (`internal/infrastructure/util/llm.go: OfRunner`) — wires the agent to `runner.New` with in-memory session/memory services (`memory.InMemoryService()`, `session.InMemoryService()`) and `AutoCreateSession: true`. `AppName` comes from `targetAgent.Name()` (i.e. `AGENT_NAME`), not a separate parameter.
5. **HTTP layer** (`internal/handler/handler.go`) — `RegisterRoutes(authMiddleware)` mounts `POST /v1/valuations` (wrapped in `authMiddleware`) and `GET /health` (left open for liveness/readiness probes) on a fresh `http.ServeMux`.
   - `askToAgent` decodes a `model.ChatRequest{sessionId, userId, message}` body via `sonic`, generates a session ID (`util.NewSessionId()`, a `uuid.New()`) when none is supplied, runs the agent via `runner.Run`, concatenates all non-thought text parts from the returned event stream, and writes the result back as **plain text** (not JSON) with a 200 status. A run-loop error is returned via `http.Error` with the raw error string — no generic `{"error": ...}` envelope.
   - `healthCheck` returns `model.HealthCheckResponse{status, checked_at}` as JSON.
6. **Middleware** (`internal/infrastructure/middeware` — package name intentionally matches the misspelled directory) — `Chain(handler, mws...)` composes `Logging` (logs status/method/URI/duration per request) and `Recovery` (recovers panics into a 500) around the mux. `Auth(client *keycloak.Client)` (`auth.go`) is a `Middleware` factory that requires a `Bearer` JWT, verifies it against Keycloak's JWKS via `keycloak.Client.VerifyToken`, and stashes a `middleware.User{ID, Username, Roles}` in the request context under `middleware.UserKey`. `RateLimit(client *redis.Client, limit int, window time.Duration)` (`ratelimit.go`) is a Redis-backed fixed-window limiter (`INCR` + `EXPIRE NX` per key) keyed on `middleware.UserKey` when set (so it must run after `Auth`), else on remote IP; it returns `429` with a `Retry-After` header when exceeded, and fails open (logs and allows the request) if Redis is unreachable. `Combine(mws...)` folds several middlewares into one, applied in order. In `cmd/api/server.go`, `Auth` and `RateLimit` are combined via `middleware.Combine` and passed as the single `authMiddleware` argument to `handler.RegisterRoutes`, so both apply only to `/v1/valuations`, not to the whole mux.
7. **Keycloak client** (`internal/infrastructure/external/keycloak/keycloak.go`) — `NewClient(ctx, baseUrl, realm, clientId)` discovers the realm's OIDC configuration (`{baseUrl}/realms/{realm}`) via `coreos/go-oidc` and builds a JWKS-backed `oidc.IDTokenVerifier` with `SkipClientIDCheck: true` (Keycloak access tokens are usually issued with `aud=account`, not the requesting client's ID). `VerifyToken` validates signature/issuer/expiry and unmarshals `sub`, `preferred_username`, `email`, and `realm_access.roles` into `Claims`.
8. **Redis client** (`internal/infrastructure/external/redisclient/redis.go`) — `NewClient(ctx, addr, password, db)` builds a `*redis.Client` (`github.com/redis/go-redis/v9`) and pings it before returning, so a misconfigured/unreachable Redis fails startup fast rather than surfacing later inside the rate limiter.

All ADK object construction (skill tool, agent, runner) is centralized in `internal/infrastructure/util`, keeping `cmd/api/server.go` a straight-line composition root.

JSON encoding/decoding uses `github.com/bytedance/sonic` rather than `encoding/json`. There is no monetary-decimal dependency in this version — the old flat-stub `valuation` package has been replaced entirely by the skill-based tool described above.

No test files exist in the repository yet (`find . -name '*_test.go'` is empty).
