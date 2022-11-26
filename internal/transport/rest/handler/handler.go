package handler

import (
	_ "github.com/Vasily-van-Zaam/ushortener/docs"
	"github.com/go-chi/chi/v5"
	httpSwagger "github.com/swaggo/http-swagger"
)

type Handlers struct {
	shortener *ShortenerHandler
	/// other handlers
}

func NewHandlers(shortener *ShortenerHandler) *Handlers {
	return &Handlers{
		shortener: shortener,
	}
}

func (h *Handlers) InitAPI(r *chi.Mux) {

	r.Get("/swagger-docs/*", httpSwagger.Handler())

	r.Get("/", h.shortener.GetURL)
	r.Get("/{id}", h.shortener.GetURL)
	r.Post("/", h.shortener.SetURL)

}
