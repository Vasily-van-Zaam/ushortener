package service

import (
	"context"

	"github.com/Vasily-van-Zaam/ushortener/internal/core"
)

type ShortenerStorage interface {
	GetURL(ctx context.Context, id string) (string, error)
	SetURL(ctx context.Context, link string) (string, error)
	Close() error
}

type Service struct {
	storage ShortenerStorage
}

func NewService(s ShortenerStorage) *Service {
	return &Service{
		s,
	}
}

func (s *Service) GetURL(ctx context.Context, id string) (string, error) {
	res, err := s.storage.GetURL(ctx, id)
	if err != nil {
		return "null", err
	}
	return res, nil
}
func (s *Service) SetURL(ctx context.Context, link string) (string, error) {
	res, err := s.storage.SetURL(ctx, link)

	if err != nil {
		return "null", err
	}
	return res, nil
}
func (s *Service) APISetShorten(ctx context.Context, request *core.RequestAPIShorten) (*core.ResponseAPIShorten, error) {

	res, err := s.storage.SetURL(ctx, request.URL)
	if err != nil {
		return nil, err
	}
	return &core.ResponseAPIShorten{
		Result: res,
	}, nil
}
