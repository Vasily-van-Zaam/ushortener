package service

import (
	"context"
	"fmt"
	"log"

	"github.com/Vasily-van-Zaam/ushortener/internal/core"
)

type Api struct {
	storage *Storage
	config  *core.Config
	core.AUTHService
}

func NewApi(conf *core.Config, s *Storage, auth *AUTHService) *Api {
	return &Api{
		s,
		conf,
		auth,
	}
}

func (s *Api) APISetShorten(ctx context.Context, request *core.RequestAPIShorten) (*core.ResponseAPIShorten, error) {
	user := core.User{}
	user.FromAny(ctx.Value(core.USERDATA))
	l := core.Link{
		Link:   request.URL,
		UserID: user.ID,
	}
	res, err := (*s.storage).SetURL(ctx, &l)
	if err != nil {
		return nil, err
	}
	return &core.ResponseAPIShorten{
		Result: s.config.BaseURL + "/" + res,
	}, nil
}

func (s *Api) APIGetUserURLS(ctx context.Context) ([]*core.ResponseAPIUserURL, error) {
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
