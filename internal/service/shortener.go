package service

import "context"

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
