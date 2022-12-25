package handler

import (
	"context"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"

	"github.com/Vasily-van-Zaam/ushortener/internal/core"
	"github.com/go-chi/chi/v5"
)

type BasicService interface {
	GetURL(ctx context.Context, link string) (string, error)
	SetURL(ctx context.Context, link string) (string, error)
	core.AUTHService
	// APISetShorten(ctx context.Context, request *core.RequestAPIShorten) (*core.ResponseAPIShorten, error)
}

type BasicHandler struct {
	service *BasicService
	config  *core.Config
}

func NewBasic(s BasicService, conf *core.Config) *BasicHandler {
	return &BasicHandler{
		service: &s,
		config:  conf,
	}
}

// @Tags         Main
// @Summary      Get link shortener
// @Description  get shortener link by ID
// @Accept       plain
// @Produce      plain
// @Param        id   path      string  true  "link ID"
// @Success		 307  {string}  "redirect response"
// @Header		 307 {string}  Location "https://some.com/link"
// @Failure      400  {string} 	"error"
// @Router       /{id} [get].
func (h *BasicHandler) GetURL(w http.ResponseWriter, r *http.Request) {
	service := *h.service
	ctx := r.Context()

	link := chi.URLParam(r, "id")

	res, err := service.GetURL(ctx, strings.TrimSpace(link))

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		h.config.LogResponse(w, r, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Location", res)
	w.WriteHeader(http.StatusTemporaryRedirect)
	_, errW := w.Write([]byte(res))
	if errW != nil {
		log.Println(errW)
	}
	h.config.LogResponse(w, r, res, http.StatusTemporaryRedirect)
}

// @Tags         Main
// @Summary      Set link shortener
// @Description  set shortener link
// @Accept       plain
// @Produce      plain
// @Param        link   body     string  true  "your site link"
// @Success		 201  {string}  "http://localhost:8080/1"
// @Failure      400  {string} 	"error"
// @Router       / [post].
func (h *BasicHandler) SetURL(w http.ResponseWriter, r *http.Request) {
	service := *h.service
	ctx := r.Context()
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, fmt.Sprintf("error get body: %s", err.Error()), http.StatusBadRequest)
		h.config.LogResponse(w, r, err.Error(), http.StatusBadRequest)
		return
	}

	if len(body) == 0 {
		w.WriteHeader(http.StatusCreated)
		_, errW := w.Write([]byte(h.config.BaseURL))
		if errW != nil {
			log.Println(errW)
		}
		h.config.LogResponse(w, r, h.config.BaseURL, http.StatusCreated)
		return
	}

	res, err := service.SetURL(ctx, strings.TrimSpace(string(body)))
	if err != nil {
		http.Error(w, fmt.Sprintf("error: %s", err.Error()), http.StatusBadRequest)
		h.config.LogResponse(w, r, err.Error(), http.StatusBadRequest)
		return
	}
	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(http.StatusCreated)
	_, errW := w.Write([]byte(res))
	if errW != nil {
		log.Println(errW)
	}
	h.config.LogResponse(w, r, res, http.StatusCreated)
}
