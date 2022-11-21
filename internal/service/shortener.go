package service

import "context"

type ShortenerStorage interface {
	GetUrl(ctx context.Context, link string) (string, error)
	SetUrl(ctx context.Context, link string) (string, error)
}

type Service struct {
	storage ShortenerStorage
}

func NewService(s ShortenerStorage) *Service {
	return &Service{
		s,
	}
}

func (s *Service) GetUrl(ctx context.Context, link string) (string, error) {
	return link, nil
}
func (s *Service) SetUrl(ctx context.Context, link string) (string, error) {
	return "1234", nil
}
