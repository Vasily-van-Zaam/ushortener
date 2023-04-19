// Server rest api.
package rest

import (
	"log"
	"net/http"
	"time"

	"github.com/Vasily-van-Zaam/ushortener/internal/core"
	"github.com/go-chi/chi/v5"
)

// Route.
type router interface {
	Run(string) error
}

type server struct {
	router *chi.Mux
	config *core.Config
}

// Create new route.
func New(conf *core.Config, h *chi.Mux) (router, error) {
	return &server{
		router: h,
		config: conf,
	}, nil
}

// Function server Run.
func (s *server) Run(addresPort string) error {
	log.Println("START SERVER ", addresPort, s.config.ServerTimeout)
	server := &http.Server{
		Addr:              addresPort,
		ReadHeaderTimeout: time.Duration(s.config.ServerTimeout) * time.Second,
		Handler:           s.router,
	}
	return server.ListenAndServe()
}
