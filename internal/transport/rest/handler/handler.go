package handler

import (
	// _ "github.com/Vasily-van-Zaam/ushortener/docs".

	"github.com/go-chi/chi/v5"
	httpSwagger "github.com/swaggo/http-swagger"
)

type Handlers struct {
	basic *BasicHandler
	api   *APIHandler
	/// other handlers
}

func NewHandlers(basic *BasicHandler, api *APIHandler) *Handlers {	
	return &Handlers{
		basic: basic,
		api:   api,
	}
}

func (h *Handlers) InitAPI(r *chi.Mux) {
	r.Get("/swagger-docs/*", httpSwagger.Handler())

	r.Get("/", h.basic.GetURL)
	r.Get("/{id}", h.basic.GetURL)
	r.Post("/", h.basic.SetURL)
	r.Post("/api/shorten", h.api.APISetShorten)
}
