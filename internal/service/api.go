package service

import (
	"context"
	"errors"
	"fmt"
	"log"

	"github.com/Vasily-van-Zaam/ushortener/internal/core"
)

type API struct {
	storage Storage
	config  *core.Config
	core.AUTHService
}

func NewAPI(conf *core.Config, s *Storage, auth *AUTHService) *API {
	return &API{
		*s,
		conf,
		auth,
	}
}

func (s *API) APISetShorten(ctx context.Context, request *core.RequestAPIShorten) (*core.ResponseAPIShorten, error) {
	user := core.User{}

	err := user.SetUserIDFromContext(ctx)
	if err != nil {
		return nil, err
	}
	l := core.Link{
		Link: request.URL,
		UUID: user.ID,
	}
	res, err := s.storage.SetURL(ctx, &l)
	if err != nil && !errors.Is(err, core.NewErrConflict()) {
		return nil, err
	}
	return &core.ResponseAPIShorten{
		Result: s.config.BaseURL + "/" + res,
	}, err
}

func (s *API) APIGetUserURLS(ctx context.Context) ([]*core.ResponseAPIUserURL, error) {
	user := core.User{}
	err := user.SetUserIDFromContext(ctx)
	if err != nil {
		return nil, err
	}
	res, err := s.storage.GetUserURLS(ctx, user.ID)
	if err != nil {
		log.Println("ERROR: GetUserURLS", err)
	}

	resAPI := make([]*core.ResponseAPIUserURL, 0, 10)
	for _, r := range res {
		log.Println(r)
		if r != nil {
			resAPI = append(resAPI, &core.ResponseAPIUserURL{
				ShortURL:    fmt.Sprint(s.config.BaseURL, "/", r.ID),
				OriginalURL: r.Link,
			})
		}
	}

	return resAPI, err
}

func (s *API) APISetShortenBatch(ctx context.Context, request []*core.RequestAPIShortenBatch) ([]*core.ResponseAPIShortenBatch, error) {
	user := core.User{}
	err := user.SetUserIDFromContext(ctx)
	if err != nil {
		return nil, err
	}
	links := make([]*core.Link, 0, 10)
	for _, r := range request {
		links = append(links, &core.Link{
			UUID: user.ID,
			Link: r.OriginalURL,
		})
	}
	res := make([]*core.ResponseAPIShortenBatch, 0, 10)
	resDB, err := s.storage.SetURLSBatch(ctx, links)

	for i, r := range resDB {
		res = append(res, &core.ResponseAPIShortenBatch{
			CorrelationID: request[i].CorrelationID,
			ShortURL:      fmt.Sprint(s.config.BaseURL, "/", r.ID),
		})
	}

	return res, err
}

func (s *API) APIDeleteUserURLS(ctx context.Context, urls []*string) error {
	user := core.User{}
	err := user.SetUserIDFromContext(ctx)
	if err != nil {
		return err
	}

	// go func() {
	// log.Println(user.ID)
	err = s.storage.DeleteURLSBatch(ctx, urls, user.ID)
	if err != nil {
		return err
	}
	// }()

	return nil
}
