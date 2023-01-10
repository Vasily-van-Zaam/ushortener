package service

import (
	"context"
	"log"

	"github.com/Vasily-van-Zaam/ushortener/internal/core"
)

type Storage interface {
	GetURL(ctx context.Context, id string) (string, error)
	SetURL(ctx context.Context, link *core.Link) (string, error)
	GetUserURLS(ctx context.Context, userID string) ([]*core.Link, error)
	SetURLSBatch(ctx context.Context, links []*core.Link) ([]*core.Link, error)
	Ping(ctx context.Context) error
	Close() error
}

type AUTHService struct {
	storage Storage
	config  *core.Config
}

func NewAuth(conf *core.Config, s *Storage) *AUTHService {
	return &AUTHService{
		storage: *s,
		config:  conf,
	}
}
func (a *AUTHService) CreateUser() {
	log.Println("HELLO USER CREATE")
}
