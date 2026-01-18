package api

import (
	"net/http"
	"log/slog"
	"context"
	"fmt"
	"time"

	"github.com/go-chi/chi/v5"
    "github.com/go-chi/chi/v5/middleware"
	"github.com/seeques/test_junior/internal/config"
	"github.com/seeques/test_junior/internal/storage"
	"github.com/seeques/test_junior/internal/handler"

	_ "github.com/seeques/test_junior/docs"
	httpSwagger "github.com/swaggo/http-swagger"
)

type Server struct {
	router chi.Router
	postgresStorage *storage.PostgresStorage
	port string
    cfg config.Config
	httpServer *http.Server
}

func NewServer(storage *storage.PostgresStorage, cfg config.Config) *Server {
	s := &Server{
        router: chi.NewRouter(),
        postgresStorage: storage,
        port: cfg.Port,
        cfg: cfg,
    }
	s.SetupRoutes()
    return s
}

func (s *Server) SetupRoutes() {
	// middleware
    s.router.Use(middleware.Logger)
    s.router.Use(middleware.Recoverer) 
    s.router.Use(middleware.RequestID) // generates unique id for request and attaches it to the context

	h := handler.NewHandler(s.postgresStorage, s.cfg)

	s.router.Get("/swagger/*", httpSwagger.WrapHandler)

	s.router.Route("/api/v1", func(r chi.Router){
		r.Post("/subscriptions", h.Create)
		r.Get("/subscriptions", h.List)
		r.Get("/subscriptions/total-cost", h.TotalCost)
		r.Get("/subscriptions/{id}", h.GetById)
		r.Put("/subscriptions/{id}", h.Update)
		r.Delete("/subscriptions/{id}", h.Delete)
	})
}

func (s *Server) Run() error {
	s.httpServer = &http.Server{
        Addr:         fmt.Sprintf(":%s", s.cfg.Port),
        Handler:      s.router,
        ReadTimeout:  10 * time.Second,
        WriteTimeout: 10 * time.Second,
        IdleTimeout:  30 * time.Second,
    }

	slog.Info("starting HTTP server", "port", s.cfg.Port)

    return s.httpServer.ListenAndServe()
}

func (s *Server) Shutdown(ctx context.Context) error {
    slog.Info("shutting down HTTP server")
    return s.httpServer.Shutdown(ctx)
}