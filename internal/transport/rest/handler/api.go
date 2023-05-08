// Handlers API
package handler

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/Vasily-van-Zaam/ushortener/internal/core"
)

// Implements Service.
type apiService interface {
	APISetShorten(ctx context.Context, request *core.RequestAPIShorten) (*core.ResponseAPIShorten, error)
	APISetShortenBatch(ctx context.Context,
		request []*core.RequestAPIShortenBatch,
	) ([]*core.ResponseAPIShortenBatch, error)
	APIGetUserURLS(ctx context.Context) ([]*core.ResponseAPIUserURL, error)
	APIDeleteUserURLS(ctx context.Context, ids []*string) error
	APIGetStats(r *http.Request) (*core.Stats, error)

	core.AUTHService
}

// Api structure.
type apiHandler struct {
	Service apiService
	Config  *core.Config
}

func newAPI(conf *core.Config, s apiService) *apiHandler {
	return &apiHandler{
		Service: s,
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
func (h *apiHandler) apiSetShorten(w http.ResponseWriter, r *http.Request) {
	service := h.Service

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
	var conflictError *core.ErrConflict
	if err != nil {

	} else {
		res, errAPI := service.APISetShorten(ctx, &query)
		if errAPI != nil && !errors.Is(errAPI, core.NewErrConflict()) {
			http.Error(w, fmt.Sprintf("error: %s", errAPI.Error()), http.StatusBadRequest)
			h.Config.LogResponse(w, r, errAPI.Error(), http.StatusBadRequest)
			return
		}
		if errors.Is(errAPI, core.NewErrConflict()) {
			conflictError = core.NewErrConflict()
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
	if conflictError != nil {
		w.WriteHeader(http.StatusConflict)
	} else {
		w.WriteHeader(http.StatusCreated)
	}
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
func (h *apiHandler) apiGetUserURLS(w http.ResponseWriter, r *http.Request) {
	service := h.Service

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

// //////////////////////////////////////////////////////////////
// TODO дописать инструкцию
// @Tags         API
// @Summary      Api Delete user urls
// @Accept       plain
// @Produce      plain
// @Success		 202
// @Failure      400  {string} 	"error"
// @Router       /api/user/urls [delete].
func (h *apiHandler) apiDeleteUserURLS(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, fmt.Sprintf("error get body: %s", err.Error()), http.StatusBadRequest)
		h.Config.LogResponse(w, r, err.Error(), http.StatusBadRequest)
		return
	}
	request := make([]*string, 0)
	err = json.Unmarshal(body, &request)
	if err != nil {
		http.Error(w, fmt.Sprintf("error get body: %s", err.Error()), http.StatusBadRequest)
		h.Config.LogResponse(w, r, err.Error(), http.StatusBadRequest)
		return
	}

	err = h.Service.APIDeleteUserURLS(r.Context(), request)
	if err != nil {
		http.Error(w, fmt.Sprintf("error delte urls: %s", err.Error()), http.StatusBadRequest)
		h.Config.LogResponse(w, r, err.Error(), http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusAccepted)
}

// TODO дописать инструкциюы
// @Router       /api/shorten/batch [post].
func (h *apiHandler) apiSetShortenBatch(w http.ResponseWriter, r *http.Request) {
	service := h.Service

	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, fmt.Sprintf("error get body: %s", err.Error()), http.StatusBadRequest)
		h.Config.LogResponse(w, r, err.Error(), http.StatusBadRequest)
		return
	}
	request := make([]*core.RequestAPIShortenBatch, 0, 10)
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
	if errors.Is(errAPI, core.NewErrConflict()) {
		w.WriteHeader(http.StatusConflict)
	} else {
		w.WriteHeader(http.StatusCreated)
	}

	response, _ := json.Marshal(res)
	_, err = w.Write(response)
	if err != nil {
		log.Println(err)
	}
	h.Config.LogResponse(w, r, string(response), http.StatusCreated)
}

// TODO: add swager
// Gets statistics urls, users.
func (h *apiHandler) apiGetStats(w http.ResponseWriter, r *http.Request) {
	var (
		resp *core.Stats
		err  error
	)

	resp, err = h.Service.APIGetStats(r)
	if err != nil {
		if err.Error() == "403" {
			w.WriteHeader(http.StatusForbidden)
			return
		}
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	response, _ := json.Marshal(resp)
	_, err = w.Write(response)
	if err != nil {
		log.Println(err)
	}
	h.Config.LogResponse(w, r, string(response), http.StatusCreated)
}
