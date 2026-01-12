package api

import (
	"net/http"

	"github.com/go-chi/chi/v5"
    // "github.com/go-chi/chi/v5/middleware"
	"github.com/seeques/test_junior/internal/config"
	"github.com/seeques/test_junior/internal/storage"
)

type Server struct {
	router chi.Router
	postgresStorage *storage.PostgresStorage
	port string
    cfg config.Config
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

func (s *Server) SetupRoutes() {}

func (s *Server) Run() error {
    return http.ListenAndServe(":"+s.port, s.router)
}