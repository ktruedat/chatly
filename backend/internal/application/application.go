package application

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/ktruedat/chatly/internal/application/config"
	"github.com/ktruedat/chatly/internal/application/handlers/websocket"
	"github.com/ktruedat/chatly/internal/application/services/openrouter"
)

type Application struct {
	config *config.Config
	server *http.Server
	router chi.Router
}

func New(cfg *config.Config) *Application {
	router := chi.NewRouter()

	router.Use(middleware.Logger)
	router.Use(middleware.Recoverer)
	router.Use(middleware.RequestID)
	router.Use(middleware.RealIP)
	router.Use(middleware.Timeout(60 * time.Second))

	// CORS for development
	router.Use(
		func(next http.Handler) http.Handler {
			return http.HandlerFunc(
				func(w http.ResponseWriter, r *http.Request) {
					w.Header().Set("Access-Control-Allow-Origin", "*")
					w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
					w.Header().Set("Access-Control-Allow-Headers", "Accept, Authorization, Content-Type, X-CSRF-Token")

					if r.Method == "OPTIONS" {
						w.WriteHeader(http.StatusOK)
						return
					}

					next.ServeHTTP(w, r)
				},
			)
		},
	)

	return &Application{
		config: cfg,
		router: router,
	}
}

func (a *Application) Initialize() error {
	// Create services
	openRouterSvc := openrouter.NewService(a.config)

	wsHandler := websocket.NewHandlers(a.router, a.config, openRouterSvc)
	if err := wsHandler.Register(); err != nil {
		return fmt.Errorf("failed to register websocket handlers: %w", err)
	}

	a.router.Get(
		"/health", func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte("OK"))
		},
	)

	// Create HTTP server
	addr := fmt.Sprintf("%s:%d", a.config.Server.Host, a.config.Server.Port)
	a.server = &http.Server{
		Addr:         addr,
		Handler:      a.router,
		ReadTimeout:  a.config.Server.ReadTimeout,
		WriteTimeout: a.config.Server.WriteTimeout,
		IdleTimeout:  a.config.Server.IdleTimeout,
	}

	return nil
}

func (a *Application) Run() error {
	serverErrors := make(chan error, 1)
	go func() {
		log.Printf("Starting server on %s", a.server.Addr)
		serverErrors <- a.server.ListenAndServe()
	}()

	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, os.Interrupt, syscall.SIGTERM)

	select {
	case err := <-serverErrors:
		return fmt.Errorf("server error: %w", err)
	case sig := <-shutdown:
		log.Printf("Received shutdown signal: %v", sig)

		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()

		if err := a.server.Shutdown(ctx); err != nil {
			if err := a.server.Close(); err != nil {
				return fmt.Errorf("failed to close server: %w", err)
			}
			return fmt.Errorf("failed to shutdown server gracefully: %w", err)
		}

		log.Println("Server stopped gracefully")
	}

	return nil
}
