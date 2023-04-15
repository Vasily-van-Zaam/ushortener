package handler

import (
	"context"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"

	"github.com/Vasily-van-Zaam/ushortener/internal/core"
	"github.com/go-chi/chi/v5"
)

type basicService interface {
	GetURL(ctx context.Context, link string) (string, error)
	SetURL(ctx context.Context, link string) (string, error)
	Ping(ctx context.Context) error
	core.AUTHService
}

type basicHandler struct {
	Service basicService
	Config  *core.Config
}

func newBasic(conf *core.Config, s basicService) *basicHandler {
	return &basicHandler{
		Service: s,
		Config:  conf,
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
func (h *basicHandler) GetURL(w http.ResponseWriter, r *http.Request) {
	service := h.Service
	ctx := r.Context()

	link := chi.URLParam(r, "id")

	res, err := service.GetURL(ctx, strings.TrimSpace(link))

	if err != nil {
		// http.Error(w, err.Error(), http.StatusBadRequest)
		if err.Error() == "deleted" {
			w.WriteHeader(http.StatusGone)
			h.Config.LogResponse(w, r, err.Error(), http.StatusGone)
			_, errW := w.Write(nil)
			if errW != nil {
				log.Println(errW)
			}
			return
		}
		h.Config.LogResponse(w, r, err.Error(), http.StatusBadRequest)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Location", res)
	w.WriteHeader(http.StatusTemporaryRedirect)
	_, errW := w.Write([]byte(res))
	if errW != nil {
		log.Println(errW)
	}
	h.Config.LogResponse(w, r, res, http.StatusTemporaryRedirect)
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
func (h *basicHandler) SetURL(w http.ResponseWriter, r *http.Request) {
	service := h.Service
	ctx := r.Context()
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, fmt.Sprintf("error get body: %s", err.Error()), http.StatusBadRequest)
		h.Config.LogResponse(w, r, err.Error(), http.StatusBadRequest)
		return
	}

	if len(body) == 0 {
		w.WriteHeader(http.StatusCreated)
		_, errW := w.Write([]byte(h.Config.BaseURL))
		if errW != nil {
			log.Println(errW)
		}
		h.Config.LogResponse(w, r, h.Config.BaseURL, http.StatusCreated)
		return
	}

	res, err := service.SetURL(ctx, strings.TrimSpace(string(body)))
	if err != nil && !errors.Is(err, core.NewErrConflict()) {
		http.Error(w, fmt.Sprintf("error: %s", err.Error()), http.StatusBadRequest)
		h.Config.LogResponse(w, r, err.Error(), http.StatusBadRequest)
		return
	}
	w.Header().Set("Content-Type", "text/plain")
	if errors.Is(err, core.NewErrConflict()) {
		w.WriteHeader(http.StatusConflict)
	} else {
		w.WriteHeader(http.StatusCreated)
	}

	_, errW := w.Write([]byte(res))
	if errW != nil {
		log.Println(errW)
	}
	h.Config.LogResponse(w, r, res, http.StatusCreated)
}

func (h *basicHandler) ping(w http.ResponseWriter, r *http.Request) {
	if h.Service.Ping(r.Context()) != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
	w.WriteHeader(http.StatusOK)
}
