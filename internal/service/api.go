package service

import (
	"context"
	"errors"
	"fmt"
	"log"

	"github.com/Vasily-van-Zaam/ushortener/internal/core"
)

type API struct {
	storage *Storage
	config  *core.Config
	core.AUTHService
}

func NewAPI(conf *core.Config, s *Storage, auth *AUTHService) *API {
	return &API{
		s,
		conf,
		auth,
	}
}

func (s *API) APISetShorten(ctx context.Context, request *core.RequestAPIShorten) (*core.ResponseAPIShorten, error) {
	user := core.User{}
	user.FromAny(ctx.Value(core.USERDATA))
	l := core.Link{
		Link: request.URL,
		UUID: user.ID,
	}
	res, err := (*s.storage).SetURL(ctx, &l)
	if err != nil && !errors.Is(err, core.NewErrConflict()) {
		return nil, err
	}
	return &core.ResponseAPIShorten{
		Result: s.config.BaseURL + "/" + res,
	}, err
}

func (s *API) APIGetUserURLS(ctx context.Context) ([]*core.ResponseAPIUserURL, error) {
	user := core.User{}
	user.FromAny(ctx.Value(core.USERDATA))
	res, err := (*s.storage).GetUserURLS(ctx, user.ID)
	if err != nil {
		log.Println("ERROR: GetUserURLS", err)
	}

	resAPI := []*core.ResponseAPIUserURL{}
	for _, r := range res {
		resAPI = append(resAPI, &core.ResponseAPIUserURL{
			ShortURL:    fmt.Sprint(s.config.BaseURL, "/", r.ID),
			OriginalURL: r.Link,
		})
	}

	return resAPI, err
}

func (s *API) APISetShortenBatch(ctx context.Context, request []*core.RequestAPIShortenBatch) ([]*core.ResponseAPIShortenBatch, error) {
	user := core.User{}
	user.FromAny(ctx.Value(core.USERDATA))
	links := []*core.Link{}
	for _, r := range request {
		links = append(links, &core.Link{
			UUID: user.ID,
			Link: r.OriginalURL,
		})
	}
	res := []*core.ResponseAPIShortenBatch{}
	resDB, err := (*s.storage).SetURLSBatch(ctx, links)

	// if err != nil && !errors.Is(err, core.NewErrConflict()) {
	// 	return nil, err
	// }
	for i, r := range resDB {
		res = append(res, &core.ResponseAPIShortenBatch{
			CorrelationID: request[i].CorrelationID,
			ShortURL:      fmt.Sprint(s.config.BaseURL, "/", r.ID),
		})
	}

	return res, err
}
