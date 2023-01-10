package service

import (
	"context"
	"errors"
	"log"

	"github.com/Vasily-van-Zaam/ushortener/internal/core"
)

type BasicService struct {
	storage Storage
	config  *core.Config
	core.AUTHService
}

func NewBasic(conf *core.Config, s *Storage, auth *AUTHService) *BasicService {
	return &BasicService{
		*s,
		conf,
		auth,
	}
}

func (s *BasicService) GetURL(ctx context.Context, id string) (string, error) {
	res, err := s.storage.GetURL(ctx, id)
	if err != nil {
		return "null", err
	}
	return res, nil
}
func (s *BasicService) SetURL(ctx context.Context, link string) (string, error) {
	user := core.User{}
	err := user.SetUserIDFromContext(ctx)
	if err != nil {
		return "", err
	}
	l := core.Link{
		Link: link,
		UUID: user.ID,
	}
	res, err := s.storage.SetURL(ctx, &l)
	log.Println("====USER ID", user.ID)

	if err != nil && !errors.Is(err, core.NewErrConflict()) {
		return "null", err
	}
	return s.config.BaseURL + "/" + res, err
}

func (s *BasicService) Ping(ctx context.Context) error {
	return s.storage.Ping(ctx)
}
