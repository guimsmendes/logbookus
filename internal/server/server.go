package server

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"os/signal"
	"syscall"
	"time"

	"github.com/danielgtaylor/huma/v2"
	"github.com/danielgtaylor/huma/v2/adapters/humago"
	"github.com/danielgtaylor/huma/v2/autopatch"
	"github.com/guimsmendes/logbookus/config"
	"github.com/guimsmendes/logbookus/internal/db"
	"gorm.io/gorm"
)

var gracefulShutdownTimeout = 10 * time.Second

type Server struct {
	port int
	db   *gorm.DB
}

// New creates a new server instance by booting up the database and setting up
// the repositories.
func New(env config.Environment, port int) (*Server, error) {
	cfg, err := config.Load(env)
	if err != nil {
		return nil, fmt.Errorf("load config: %w", err)
	}

	postgresDB, err := db.Connect(cfg)
	if err != nil {
		return nil, fmt.Errorf("connect to database: %w", err)
	}

	return &Server{
		port: port,
		db:   postgresDB,
	}, nil
}

func (s *Server) Start(ctx context.Context) error {
	// Create a new router & API
	router := http.NewServeMux()

	// Setup huma API
	config := huma.DefaultConfig("Gophernet", "1.0.0")
	humaAPI := humago.New(router, config)

	// Set up our API request handlers. We only pass the repositories to the API, not
	// the database itself to ensure strict separation.
	//h := a
	//registerRoutes(humaAPI, h)

	// Autopatch adds PATCH support if the resource has both a GET and PUT handler.
	autopatch.AutoPatch(humaAPI)

	// Graceful shutdown
	ctx, stop := signal.NotifyContext(ctx, syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	// Setup server
	srv := &http.Server{
		Addr:    fmt.Sprintf(":%d", s.port),
		Handler: router,
	}

	go func() {
		// Start the server!
		slog.Info("server started",
			"server", fmt.Sprintf("http://localhost:%d", s.port),
			"docs", fmt.Sprintf("http://localhost:%d/docs", s.port),
		)

		err := srv.ListenAndServe()
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			slog.Error(fmt.Sprintf("listen and serve: %v", err))
		}
	}()

	// Wait for the context to be done, which will be triggered by a signal
	<-ctx.Done()

	// Create a new context with a deadline for the shutdown process
	shutdownCtx, cancel := context.WithTimeout(context.Background(), gracefulShutdownTimeout)
	defer cancel()

	// Attempt to gracefully shut down the server
	if err := srv.Shutdown(shutdownCtx); err != nil {
		slog.Error(fmt.Sprintf("graceful server shutdown failed: %v", err))
		return err
	}

	slog.Info("server stopped gracefully")
	return nil
}
