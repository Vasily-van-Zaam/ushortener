// List Handlers
package handler

import (
	// _ "github.com/Vasily-van-Zaam/ushortener/docs".

	"net/http"

	"github.com/Vasily-van-Zaam/ushortener/internal/core"
	"github.com/go-chi/chi/v5"
	httpSwagger "github.com/swaggo/http-swagger"
)

// Create hadlers.
func New(conf *core.Config, b basicService, a apiService, m ...func(http.Handler) http.Handler) *chi.Mux {
	r := chi.NewMux()
	r.Use(m...)
	basic := newBasic(conf, b)
	api := newAPI(conf, a)

	r.Get("/swagger-docs/*", httpSwagger.Handler())
	// BASIC
	r.Get("/", basic.GetURL)
	r.Get("/ping", basic.ping)
	r.Get("/{id}", basic.GetURL)
	r.Post("/", basic.SetURL)

	// API
	r.Post("/api/shorten", api.apiSetShorten)
	r.Get("/api/user/urls", api.apiGetUserURLS)
	r.Delete("/api/user/urls", api.apiDeleteUserURLS)
	r.Post("/api/shorten/batch", api.apiSetShortenBatch)
	return r
	// return &handlers{
	// 	basic: basic,
	// 	api:   api,
	// }
}
