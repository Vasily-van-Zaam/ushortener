package rest

import (
	"log"
	"net/http"

	"github.com/Vasily-van-Zaam/ushortener/internal/transsport/rest/handler"
)

type Router interface {
	Run(string) error
}

type Server struct {
	router *http.ServeMux
}

func NewServer(h *handler.Handlers) (Router, error) {
	r := http.NewServeMux()
	h.InitAPI(r)
	return &Server{
		router: r,
	}, nil
}

func (s *Server) Run(port string) error {
	log.Print("START: http://localhost", port)
	return http.ListenAndServe(port, s.router)
}
