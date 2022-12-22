package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"

	"github.com/Vasily-van-Zaam/ushortener/internal/core"
	"github.com/go-chi/chi/v5"
)

type ShortenerService interface {
	GetURL(ctx context.Context, link string) (string, error)
	SetURL(ctx context.Context, link string) (string, error)
	APISetShorten(ctx context.Context, request *core.RequestAPIShorten) (*core.ResponseAPIShorten, error)
}

type ShortenerHandler struct {
	service ShortenerService
}

func NewShortenerHandler(s ShortenerService) *ShortenerHandler {
	return &ShortenerHandler{
		service: s,
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
func (h *ShortenerHandler) GetURL(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()

	link := chi.URLParam(r, "id")

	res, err := h.service.GetURL(ctx, strings.TrimSpace(link))

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Location", res)
	w.WriteHeader(http.StatusTemporaryRedirect)
	_, errW := w.Write([]byte(res))
	if errW != nil {
		log.Println(errW)
	}
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
func (h *ShortenerHandler) SetURL(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, fmt.Sprintf("error get body: %s", err.Error()), http.StatusBadRequest)
		return
	}

	if len(body) == 0 {
		w.WriteHeader(http.StatusCreated)
		_, errW := w.Write([]byte(core.MAINDOMAIN))
		if errW != nil {
			log.Println(errW)
		}
		return
	}

	res, err := h.service.SetURL(ctx, strings.TrimSpace(string(body)))
	if err != nil {
		http.Error(w, fmt.Sprintf("error: %s", err.Error()), http.StatusBadRequest)
		return
	}
	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(http.StatusCreated)
	_, errW := w.Write([]byte(res))
	if errW != nil {
		log.Println(errW)
	}
}

// @Tags         API
// @Summary      Api Set link shortener
// @Description  set shortener link 1
// @Accept       plain
// @Produce      plain
// @Param        body body core.RequestApiShorten true "Body"
// @Success		 200  {object} core.ResponseApiShorten
// @Failure      400  {string} 	"error"
// @Router       /api/shorten [post].
func (h *ShortenerHandler) APISetShorten(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, fmt.Sprintf("error get body: %s", err.Error()), http.StatusBadRequest)
		return
	}
	query := core.RequestAPIShorten{}
	responseAPI := &core.ResponseAPIShorten{}
	err = json.Unmarshal(body, &query)
	if err != nil {

	} else {
		res, errAPI := h.service.APISetShorten(ctx, &query)
		if errAPI != nil {
			http.Error(w, fmt.Sprintf("error: %s", errAPI.Error()), http.StatusBadRequest)
			return
		}
		responseAPI = res
	}

	if err != nil {
		http.Error(w, fmt.Sprintf("error: %s", err.Error()), http.StatusBadRequest)
		return
	}

	response, _ := json.Marshal(responseAPI)

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	_, errW := w.Write(response)
	if errW != nil {
		log.Println(errW)
	}
	log.Println("==", query.URL)
}
