package main

import (
	"car-valuation-agent/internal/handler"
	"car-valuation-agent/internal/infrastructure/config"
	"car-valuation-agent/internal/infrastructure/middeware"
	"car-valuation-agent/internal/infrastructure/util"
	"context"
	"errors"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"google.golang.org/adk/v2/memory"
	"google.golang.org/adk/v2/session"
	"google.golang.org/adk/v2/tool"
)

const (
	agentInstruction = "You are a car valuer agent. Please `car-valuer skill` to get the market price instead your own knowldge. The response should only the output described on the skill"
)

func main() {
	ctx := context.Background()
	appConfig := config.Load()
	skillBasedTool, err := util.OfSkillBasedTool(ctx, appConfig.GetAgent().GetSkill().Location)

	if err != nil {
		log.Fatalf("error: %v", err.Error())
	}

	carValuerAgent, err := util.OfAgent(ctx, agentInstruction, appConfig, []tool.Toolset{skillBasedTool})
	if err != nil {
		log.Fatalf("error: %v", err.Error())
	}

	agentRunner, err := util.OfRunner(carValuerAgent, memory.InMemoryService(), session.InMemoryService(), true)
	if err != nil {
		log.Fatalf("error: %v", err.Error())
	}
	handler := handler.New(*agentRunner)
	wrappedHandler := middeware.Chain(handler.RegisterRoutes(), middeware.Logging, middeware.Recovery)

	agentServer := &http.Server{
		Addr:         appConfig.GetServer().GetAddr(),
		Handler:      wrappedHandler,
		ReadTimeout:  appConfig.GetServer().ReadTimeout,
		WriteTimeout: appConfig.GetServer().Writetimout,
	}

	go func() {
		log.Printf("listening on %s", agentServer.Addr)
		if err := agentServer.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("error: %v", err.Error())
		}
	}()
	handleGarcefullyShutdown(ctx, agentServer)
}

func handleGarcefullyShutdown(ctx context.Context, agentServer *http.Server) {
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	shutdownCtx, shutdownComplete := context.WithTimeout(ctx, 30*time.Second)
	defer shutdownComplete()
	if err := agentServer.Shutdown(shutdownCtx); err != nil {
		log.Printf("handleGarcefullyShutdown[error: %v]", err.Error())
	}
	log.Printf("the agent is successfully shotdown")
}
