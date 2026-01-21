package server

import (
	"fmt"
	"log/slog"
	"messenger-auth/internal/config"
	"messenger-auth/internal/server/handlers"
	"messenger-auth/internal/storage"
	"net/http"

	"github.com/gorilla/mux"
)

type Server struct {
	config  *config.Config
	logger  *slog.Logger
	storage storage.Storage
	router  *mux.Router
}

func NewServer(config *config.Config, logger *slog.Logger, storage storage.Storage) *Server {
	return &Server{config: config, logger: logger, storage: storage, router: mux.NewRouter()}
}

func (s *Server) ServeHTTP() error {
	s.router.Handle("/update_token", handlers.UpdateToken(s.logger.With("handler", "update_token"))).Methods("POST")
	s.router.Handle("/register", handlers.Register(s.logger.With("handler", "register"), s.storage)).Methods("POST")
	s.router.Handle("/login", handlers.Login(s.logger.With("handler", "login"), s.storage)).Methods("POST")

	return http.ListenAndServe(fmt.Sprintf("%s:%v", s.config.Hostname, s.config.Port), s.router)
}
