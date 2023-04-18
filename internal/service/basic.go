package service

import (
	"context"
	"errors"
	"fmt"

	"github.com/Vasily-van-Zaam/ushortener/internal/core"
	"github.com/Vasily-van-Zaam/ushortener/pkg/shorter"
)

type iShorter interface {
	Convert(id string) string
	UnConnvert(id string) string
}
type BasicService struct {
	storage Storage
	config  *core.Config
	shorter iShorter
	core.AUTHService
}

func NewBasic(conf *core.Config, s *Storage, auth *AUTHService) *BasicService {
	return &BasicService{
		*s,
		conf,
		shorter.NewShorter59(),
		auth,
	}
}

// Get URL, response user url.
func (s *BasicService) GetURL(ctx context.Context, id string) (string, error) {
	res, err := s.storage.GetURL(ctx, fmt.Sprint(s.shorter.UnConnvert(id)))
	if err != nil {
		return "", err
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
		Link:    link,
		UUID:    user.ID,
		Deleted: false,
	}
	res, err := s.storage.SetURL(ctx, &l)

	if err != nil && !errors.Is(err, core.NewErrConflict()) {
		return "", err
	}

	url := s.config.BaseURL + "/" + s.shorter.Convert(res)
	return url, err
}

func (s *BasicService) Ping(ctx context.Context) error {
	return s.storage.Ping(ctx)
}
