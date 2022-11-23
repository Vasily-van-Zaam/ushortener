package handler

import "net/http"

type Handlers struct {
	shortener *ShortenerHandler
	/// other handlers
}

func NewHandlers(shortener *ShortenerHandler) *Handlers {
	return &Handlers{
		shortener: shortener,
	}
}

func (h *Handlers) InitAPI(r *http.ServeMux) {
	r.HandleFunc("/", h.shortener.GetSetURL)
	// r.HandleFunc("/", h.shortener.SetURL)

	// r.GET("/v1/search", h.Search.Query)
	// r.GET("/v1/suggestion", h.Search.Suggestion)
	// r.POST("/v1/import-from-url", h.Search.ImportXmlFromUrl)
	// r.GET("/v1/by-ids", h.Search.GetProducts)
	// r.GET("/v1/get-file", h.Search.GetFileXml)

}
