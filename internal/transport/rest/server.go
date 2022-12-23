package rest

import (
	"log"
	"net/http"
	"time"

	"github.com/Vasily-van-Zaam/ushortener/internal/core"
	"github.com/Vasily-van-Zaam/ushortener/internal/transport/rest/handler"
	"github.com/Vasily-van-Zaam/ushortener/internal/transport/rest/middleware"
	"github.com/go-chi/chi/v5"
)

type Router interface {
	Run(string) error
}

type Server struct {
	router *chi.Mux
	config *core.Config
}

func NewServer(h *handler.Handlers, conf *core.Config) (Router, error) {
	r := chi.NewRouter()
	h.InitAPI(r)
	return &Server{
		router: r,
		config: conf,
	}, nil
}

func (s *Server) Run(addresPort string) error {
	log.Println("START SERVER ", addresPort, s.config.ServerTimeout)
	server := &http.Server{
		Addr:              addresPort,
		ReadHeaderTimeout: time.Duration(s.config.ServerTimeout) * time.Second,
		Handler:           middleware.GzipHandle(s.router, s.config),
	}
	return server.ListenAndServe()
}
