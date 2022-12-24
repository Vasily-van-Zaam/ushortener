package service

import (
	"context"

	"github.com/Vasily-van-Zaam/ushortener/internal/core"
)

type Api struct {
	storage ShortenerStorage
}

func NewApi(s *ShortenerStorage) *Api {
	return &Api{
		storage: *s,
	}
}

func (s *Api) APISetShorten(ctx context.Context, request *core.RequestAPIShorten) (*core.ResponseAPIShorten, error) {

	res, err := s.storage.SetURL(ctx, request.URL)
	if err != nil {
		return nil, err
	}
	return &core.ResponseAPIShorten{
		Result: res,
	}, nil
}
