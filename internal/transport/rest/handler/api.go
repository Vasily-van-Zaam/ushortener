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
	APISetShortenBatch(ctx context.Context,
		request []*core.RequestAPIShortenBatch,
	) ([]*core.ResponseAPIShortenBatch, error)
	APIGetUserURLS(ctx context.Context) ([]*core.ResponseAPIUserURL, error)
	core.AUTHService
}

type APIHandler struct {
	Service *APIService
	Config  *core.Config
}

func NewAPI(s APIService, conf *core.Config) *APIHandler {
	return &APIHandler{
		Service: &s,
		Config:  conf,
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
	service := *h.Service

	ctx := r.Context()
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, fmt.Sprintf("error get body: %s", err.Error()), http.StatusBadRequest)
		h.Config.LogResponse(w, r, err.Error(), http.StatusBadRequest)
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
			h.Config.LogResponse(w, r, errAPI.Error(), http.StatusBadRequest)
			return
		}
		responseAPI = res
	}

	if err != nil {
		http.Error(w, fmt.Sprintf("error: %s", err.Error()), http.StatusBadRequest)
		h.Config.LogResponse(w, r, err.Error(), http.StatusBadRequest)
		return
	}

	response, _ := json.Marshal(responseAPI)

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	_, errW := w.Write(response)

	if errW != nil {
		log.Println(errW)
	}

	h.Config.LogResponse(w, r, string(response), http.StatusCreated)
}

// //////////////////////////////////////////////////////////////
// TODO дописать инструкциюы
// @Tags         API
// @Summary      Api Get user urls
// @Accept       plain
// @Produce      plain
// @Success		 200  {object} core.ResponseAPIUserURL
// @Failure      400  {string} 	"error"
// @Router       /api/user/urls [get].
func (h *APIHandler) APIGetUserURLS(w http.ResponseWriter, r *http.Request) {
	service := *h.Service

	res, errAPI := service.APIGetUserURLS(r.Context())
	if errAPI != nil {
		log.Println("ERROR: ", errAPI)
	}
	if len(res) == 0 {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNoContent)
		h.Config.LogResponse(w, r, string(""), http.StatusNoContent)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	response, _ := json.Marshal(res)
	_, err := w.Write(response)
	if err != nil {
		log.Println(err)
	}

	h.Config.LogResponse(w, r, string(response), http.StatusOK)
}

// TODO дописать инструкциюы
// @Router       /api/shorten/batch [post].
func (h *APIHandler) APISetShortenBatch(w http.ResponseWriter, r *http.Request) {
	service := *h.Service

	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, fmt.Sprintf("error get body: %s", err.Error()), http.StatusBadRequest)
		h.Config.LogResponse(w, r, err.Error(), http.StatusBadRequest)
		return
	}
	request := []*core.RequestAPIShortenBatch{}
	err = json.Unmarshal(body, &request)
	if err != nil {
		http.Error(w, fmt.Sprintf("error get body: %s", err.Error()), http.StatusBadRequest)
		h.Config.LogResponse(w, r, err.Error(), http.StatusBadRequest)
		return
	}

	res, errAPI :=
		service.APISetShortenBatch(
			r.Context(),
			request,
		)
	if errAPI != nil {
		log.Println("ERROR: ", errAPI)
	}
	if len(res) == 0 {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNoContent)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)

	response, _ := json.Marshal(res)
	_, err = w.Write(response)
	if err != nil {
		log.Println(err)
	}
}
