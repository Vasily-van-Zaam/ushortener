package rest

import (
	"log"
	"net/http"
	"time"

	"github.com/Vasily-van-Zaam/ushortener/internal/core"
	"github.com/go-chi/chi/v5"
)

type middleware interface {
	Handle(next http.Handler) http.Handler
}

type router interface {
	Run(string) error
}

type server struct {
	router *chi.Mux
	config *core.Config
}

func New(conf *core.Config, h *chi.Mux) (router, error) {
	// h.Use(setMiddlewareFuncList(mws)...)

	// h.InitAPI(r)
	return &server{
		router: h,
		config: conf,
	}, nil
}

func (s *server) Run(addresPort string) error {
	log.Println("START SERVER ", addresPort, s.config.ServerTimeout)
	server := &http.Server{
		Addr:              addresPort,
		ReadHeaderTimeout: time.Duration(s.config.ServerTimeout) * time.Second,
		Handler:           s.router,
	}
	return server.ListenAndServe()
}

// func setMiddlewareFuncList(m []middleware) []func(http.Handler) http.Handler {
// 	middlewares := []func(http.Handler) http.Handler{}
// 	for _, m := range m {
// 		middlewares = append(middlewares, m.Handle)
// 	}
// 	return middlewares
// }
