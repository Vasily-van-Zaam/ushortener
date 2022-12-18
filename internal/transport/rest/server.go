package rest

import (
	"log"
	"net/http"

	"github.com/Vasily-van-Zaam/ushortener/internal/transport/rest/handler"
	"github.com/go-chi/chi/v5"
)

type Router interface {
	Run(string) error
}

type Server struct {
	router *chi.Mux
}

func NewServer(h *handler.Handlers) (Router, error) {
	r := chi.NewRouter()
	h.InitAPI(r)
	return &Server{
		router: r,
	}, nil
}

func (s *Server) Run(addresPort string) error {
	log.Print("START SERVER ", addresPort)
	return http.ListenAndServe(addresPort, s.router)
}
