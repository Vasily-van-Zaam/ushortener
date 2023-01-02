package service

import (
	"context"
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
	if err != nil {
		return nil, err
	}
	return &core.ResponseAPIShorten{
		Result: s.config.BaseURL + "/" + res,
	}, nil
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
