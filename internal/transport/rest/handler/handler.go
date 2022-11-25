package handler

import (
	"github.com/go-chi/chi/v5"
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

	r.Get("/", h.shortener.GetURL)
	r.Get("/{id}", h.shortener.GetURL)
	r.Post("/", h.shortener.SetURL)

}
