package server

import (
	"app-server/internal/config"
	rt "app-server/internal/server/router"
	"fmt"
	"net/http"

	"github.com/gorilla/handlers"
	"go.uber.org/zap"
)

type Server struct {
	Config config.Config
	logger *zap.Logger
}

func New(cfg config.Config, logger *zap.Logger) *Server {
	return &Server{Config: cfg, logger: logger}
}

func (s *Server) Run() error {
	router := rt.Create(s.Config.SecretKey)

	headers := handlers.AllowedHeaders([]string{"Content-Type"})
	methods := handlers.AllowedMethods([]string{"GET", "POST", "PUT", "DELETE"})
	origins := handlers.AllowedOrigins([]string{"*"})
	logs := handlers.CORS(headers, methods, origins)(router)

	s.logger.Info("Server started", zap.Int("port", s.Config.Port))

	return http.ListenAndServe(fmt.Sprintf("0.0.0.0:%d", s.Config.Port), logs)
}
