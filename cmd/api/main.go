package main

import (
	"context"
	"errors"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"

	"goflagsmith/internal/handlers"
	"goflagsmith/internal/service/flags"
	"goflagsmith/internal/service/router"
	"goflagsmith/internal/state"
	"goflagsmith/internal/util/hash"
)

func main() {

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	appState := state.NewState()

	api_key := os.Getenv("FLAGSMITH_API_KEY")
	flagsSvc, err := flags.NewClient(ctx, api_key)
	if err != nil {
		log.Fatalf("FATAL: failed to initialize Flagsmith client: %v", err)
	}
	appState.SetClientsReady()

	interval := 2 * time.Second
	flagsSvc.MonitorFlagsReady(ctx, appState, interval)
	flagsSvc.StartRulesSync(ctx, interval)

	h := handlers.NewAppHandler(appState, flagsSvc)

	var bc router.BucketCalculator = hash.NormalizedHash
	var now router.Now = time.Now
	eng := router.NewEngine(flagsSvc, appState, bc, now)
	rh := handlers.NewRouteHandler(eng)

	r := gin.Default()

	r.GET("/healthz", func(c *gin.Context) {
		c.String(http.StatusOK, "server is running")
	})

	r.GET("/readyz", h.Readyz)

	r.POST("/decide", rh.Handle)

	srv := &http.Server{
		Addr:    ":8080",
		Handler: r,
	}

	go func() {
		log.Printf("INFO: HTTP server listening on %s", srv.Addr)
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("FATAL: HTTP server error: %v", err)
		}
	}()

	<-ctx.Done()
	log.Println("WARN: Shutdown signal received, shutting down gracefully...")

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(shutdownCtx); err != nil {
		log.Fatalf("ERROR: Server forced to shutdown: %v", err)
	}

	log.Println("INFO: Server exited cleanly")
}
