package main

import (
	"car-valuation-agent/internal/handler"
	"car-valuation-agent/internal/infrastructure/config"
	"car-valuation-agent/internal/infrastructure/external/keycloak"
	"car-valuation-agent/internal/infrastructure/external/memoryservice"
	"car-valuation-agent/internal/infrastructure/external/redisclient"
	"car-valuation-agent/internal/infrastructure/middleware"
	"car-valuation-agent/internal/infrastructure/util"
	"context"
	"errors"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"google.golang.org/adk/v2/session"
	"google.golang.org/adk/v2/tool"

	"github.com/redis/go-redis/v9"
)

const (
	agentInstruction = "You are a car valuer agent. Use `car-valuer` skill to get the market price instead your own knowldge. The response should only the output described on the skill"
)

// must returns value, terminating the process if constructing it failed. It exists to keep
// main a straight-line composition root instead of a long chain of repeated error checks.
func must[T any](value T, err error) T {
	if err != nil {
		log.Fatalf("error: %v", err.Error())
	}
	return value
}

func main() {
	ctx := context.Background()
	appConfig := must(config.Load())

	skillBasedTool := must(util.OfSkillBasedTool(ctx, appConfig.GetAgent().GetSkill().Location))
	carValuerAgent := must(util.OfAgent(ctx, agentInstruction, appConfig, []tool.Toolset{skillBasedTool}))
	externalMemoryService := must(memoryservice.OfExternalMemoryService(ctx, "", appConfig.GetAgent().GetModel()))
	agentRunner := must(util.OfRunner(carValuerAgent, externalMemoryService, session.InMemoryService(), true))

	keycloakClient := must(keycloak.NewClient(ctx, appConfig.GetKeycloak().BaseUrl,
		appConfig.GetKeycloak().Realm, appConfig.GetKeycloak().ClientId))

	redisClient := must(redisclient.NewClient(ctx, appConfig.GetRedis().Addr,
		appConfig.GetRedis().Password, appConfig.GetRedis().DB))
	defer closeRedisClient(redisClient)

	rateLimitMiddleware := middleware.RateLimit(redisClient, appConfig.GetRateLimit().Limit, appConfig.GetRateLimit().Window)

	handler := handler.New(*agentRunner, appConfig.GetWaitTimeout())
	protectedRoutes := middleware.Combine(middleware.Auth(keycloakClient), rateLimitMiddleware)
	wrappedHandler := middleware.Chain(handler.RegisterRoutes(protectedRoutes), middleware.Logging, middleware.Recovery)

	agentServer := &http.Server{
		Addr:         appConfig.GetServer().GetAddr(),
		Handler:      wrappedHandler,
		ReadTimeout:  appConfig.GetServer().ReadTimeout,
		WriteTimeout: appConfig.GetServer().WriteTimeout,
	}

	go func() {
		log.Printf("listening on %s", agentServer.Addr)
		err := agentServer.ListenAndServe()
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("error: %v", err.Error())
		}
	}()
	handleGracefullyShutdown(ctx, agentServer)
}

func closeRedisClient(redisClient *redis.Client) {
	err := redisClient.Close()
	if err != nil {
		log.Printf("error closing redis client: %v", err)
	}
}

func handleGracefullyShutdown(ctx context.Context, agentServer *http.Server) {
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	shutdownCtx, shutdownComplete := context.WithTimeout(ctx, 30*time.Second)
	defer shutdownComplete()

	err := agentServer.Shutdown(shutdownCtx)
	if err != nil {
		log.Printf("handleGarcefullyShutdown[error: %v]", err.Error())
	}
	log.Printf("the agent is successfully shotdown")
}
