package main

import (
	"application"
	"application/route"
	"context"
	"errors"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"
)

// Server holds the main HTTP server, logger and DI container.
type Server struct {
	Container *application.Container
	Logger    *slog.Logger
	HTTP      *http.Server
}

// NewServer initializes the Server with necessary configurations, DI container, and routes.
func NewServer() *Server {
	container := application.NewContainer()
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))

	server := &Server{
		Container: container,
		Logger:    logger,
		HTTP: &http.Server{
			Addr:         ":" + strconv.Itoa(container.Config.Get().Port),
			Handler:      route.Register(container),
			ReadTimeout:  10 * time.Second,
			WriteTimeout: 10 * time.Second,
			IdleTimeout:  15 * time.Second,
			ErrorLog:     slog.NewLogLogger(logger.Handler(), slog.LevelError),
		},
	}

	return server
}

// Start runs the HTTP server and listens for shutdown signals for graceful stopping.
func (s *Server) Start() error {
	s.Logger.Info("Starting server", "address", s.HTTP.Addr)

	// Start the server in a separate goroutine to enable graceful shutdown
	go func() {
		if err := s.HTTP.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			s.Logger.Error("Server encountered an error", "error", err)
			os.Exit(1)
		}
	}()

	// Handle graceful shutdown on interrupt signals
	s.gracefulShutdown()
	return nil
}

// gracefulShutdown manages graceful shutdown of the server upon receiving an interrupt signal.
func (s *Server) gracefulShutdown() {
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	<-stop

	s.Logger.Info("Shutting down server...")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := s.HTTP.Shutdown(ctx); err != nil {
		s.Logger.Error("Server forced to shutdown", "error", err)
	} else {
		s.Logger.Info("Server stopped gracefully")
	}
}

func main() {
	server := NewServer()
	if err := server.Start(); err != nil {
		server.Logger.Error("Failed to start server", "error", err)
		os.Exit(1)
	}
}
