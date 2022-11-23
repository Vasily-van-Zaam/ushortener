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
	res, err := s.storage.GetUrl(ctx, link)
	if err != nil {
		return "null", err
	}
	return res, nil
}
func (s *Service) SetUrl(ctx context.Context, link string) (string, error) {
	res, err := s.storage.SetUrl(ctx, link)

	if err != nil {
		return "null", err
	}
	return res, nil
}
