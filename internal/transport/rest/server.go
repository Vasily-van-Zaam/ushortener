package rest

import (
	"log"
	"net/http"
	"time"

	"github.com/Vasily-van-Zaam/ushortener/internal/core"
	"github.com/Vasily-van-Zaam/ushortener/internal/transport/rest/handler"
	"github.com/go-chi/chi/v5"
)

type Middleware interface {
	Handle(next http.Handler) http.Handler
}

type Router interface {
	Run(string) error
}

type Server struct {
	router *chi.Mux
	config *core.Config
}

func NewServer(h *handler.Handlers, conf *core.Config, mws []Middleware) (Router, error) {
	r := chi.NewRouter()
	r.Use(setMiddlewareFuncList(mws)...)
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
		Handler:           s.router,
	}
	return server.ListenAndServe()
}

func setMiddlewareFuncList(m []Middleware) []func(http.Handler) http.Handler {
	middlewares := []func(http.Handler) http.Handler{}
	for _, m := range m {
		middlewares = append(middlewares, m.Handle)
	}
	return middlewares
}
