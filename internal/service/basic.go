package service

import (
	"context"
)

type ShortenerStorage interface {
	GetURL(ctx context.Context, id string) (string, error)
	SetURL(ctx context.Context, link string) (string, error)
	Close() error
}

type BasicService struct {
	storage ShortenerStorage
}

func NewBasic(s *ShortenerStorage) *BasicService {
	return &BasicService{
		*s,
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
	res, err := s.storage.SetURL(ctx, link)

	if err != nil {
		return "null", err
	}
	return res, nil
}
