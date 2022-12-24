package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/Vasily-van-Zaam/ushortener/internal/core"
)

type APIService interface {
	APISetShorten(ctx context.Context, request *core.RequestAPIShorten) (*core.ResponseAPIShorten, error)
}

type APIHandler struct {
	service *APIService
	config  *core.Config
}

func NewAPI(s APIService, conf *core.Config) *APIHandler {
	return &APIHandler{
		service: &s,
		config:  conf,
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
func (h *APIHandler) APISetShorten(w http.ResponseWriter, r *http.Request) {
	service := *h.service

	ctx := context.Background()
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, fmt.Sprintf("error get body: %s", err.Error()), http.StatusBadRequest)
		h.config.LogResponse(w, r, err.Error(), http.StatusBadRequest)
		return
	}
	query := core.RequestAPIShorten{}
	responseAPI := &core.ResponseAPIShorten{}
	err = json.Unmarshal(body, &query)
	if err != nil {

	} else {
		res, errAPI := service.APISetShorten(ctx, &query)
		if errAPI != nil {
			http.Error(w, fmt.Sprintf("error: %s", errAPI.Error()), http.StatusBadRequest)
			h.config.LogResponse(w, r, errAPI.Error(), http.StatusBadRequest)
			return
		}
		responseAPI = res
	}

	if err != nil {
		http.Error(w, fmt.Sprintf("error: %s", err.Error()), http.StatusBadRequest)
		h.config.LogResponse(w, r, err.Error(), http.StatusBadRequest)
		return
	}

	response, _ := json.Marshal(responseAPI)

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	_, errW := w.Write(response)

	if errW != nil {
		log.Println(errW)
	}
	h.config.LogResponse(w, r, string(response), http.StatusCreated)
}
