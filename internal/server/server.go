package server

import (
	"context"
	"errors"
	"fmt"
	"guimsmendes/personal/logbookus/internal/model"
	"guimsmendes/personal/logbookus/internal/repository"
	"log"
	"log/slog"
	"net/http"
	"os/signal"
	"syscall"
	"time"

	"github.com/danielgtaylor/huma/v2"
	"github.com/danielgtaylor/huma/v2/adapters/humago"
	"github.com/danielgtaylor/huma/v2/autopatch"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var gracefulShutdownTimeout = 10 * time.Second

type Server struct {
	port int
	repo *repository.Repository
}

// New creates a new server instance by booting up the database and setting up
// the repositories.
func New(port int) (*Server, error) {
	conn := "host=localhost user=logbookus password=1 dbname=logbookus port=5432 sslmode=disable TimeZone=Amsterdam/Netherlands"

	db, err := gorm.Open(postgres.Open(conn), &gorm.Config{})
	if err != nil {
		log.Fatalf("failed to connect database: %v", err)
	}
	if err != nil {
		return nil, fmt.Errorf("starting bolt: %v", err)
	}

	err = db.AutoMigrate(model.GetModels()...)
	if err != nil {
		return nil, fmt.Errorf("auto migrate: %w", err)
	}
	return &Server{
		port: port,
		repo: repository.New(db),
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
	h := a
	registerRoutes(humaAPI, h)

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
