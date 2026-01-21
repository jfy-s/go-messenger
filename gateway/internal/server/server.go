package server

import (
	"fmt"
	"messenger-gateway/internal/config"
	"net/http"

	"github.com/gorilla/mux"
)

type Server struct {
	config *config.Config
	router *mux.Router
}

func NewServer(cfg *config.Config) *Server {
	return &Server{
		config: cfg,
		router: mux.NewRouter(),
	}
}

func (s *Server) ServeHttp() error {
	return http.ListenAndServe(fmt.Sprintf("%s:%v", s.config.Hostname, s.config.Port), s.router)
}
