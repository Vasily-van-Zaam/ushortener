package handler

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/Vasily-van-Zaam/ushortener/internal/core"
)

type ShortenerService interface {
	GetURL(ctx context.Context, link string) (string, error)
	SetURL(ctx context.Context, link string) (string, error)
}

type ShortenerHandler struct {
	service ShortenerService
}

func NewShortenerHandler(s ShortenerService) *ShortenerHandler {
	return &ShortenerHandler{
		service: s,
	}
}

func (h *ShortenerHandler) GetSetURL(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()

	switch r.Method {
	case "POST":
		{
			ct := r.Header.Values("Content-Type")
			if len(ct) > 0 && ct[0] != "text/plain" {
				http.Error(w, "body mast be text/plain", http.StatusBadRequest)
				return
			}

			body, err := io.ReadAll(r.Body)
			if err != nil {
				http.Error(w, fmt.Sprintf("error get body: %s", err.Error()), http.StatusBadRequest)
				return
			}

			if len(body) == 0 {
				// http.Error(w, "body cannot be empty", http.StatusBadRequest)
				w.WriteHeader(http.StatusCreated)
				w.Write([]byte(core.MAINDOMAIN))
				return
			}

			res, err := h.service.SetURL(ctx, strings.TrimSpace(string(body)))
			if err != nil {
				http.Error(w, fmt.Sprintf("error: %s", err.Error()), http.StatusBadRequest)
				return
			}
			w.WriteHeader(http.StatusCreated)
			w.Write([]byte(res))

		}
	case "GET":
		{
			url := strings.Split(r.URL.Path, "/")
			link := "/"
			if len(url) >= 1 {
				link = url[1]
			}
			res, err := h.service.GetURL(ctx, strings.TrimSpace(link))

			if err != nil {
				http.Error(w, fmt.Sprintf("error: %s", err.Error()), http.StatusBadRequest)
				return
			}
			w.WriteHeader(http.StatusTemporaryRedirect)
			w.Write([]byte(res))

		}
	default:
		{
			http.Error(w, fmt.Sprintf("method %s not supported", r.Method), http.StatusBadRequest)
			return
		}
	}

}
