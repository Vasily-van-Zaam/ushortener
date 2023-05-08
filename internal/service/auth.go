// Business logic
package service

import (
	"context"
	"log"

	"github.com/Vasily-van-Zaam/ushortener/internal/core"
)

// Main function Storage for service.
type Storage interface {
	GetURL(ctx context.Context, id string) (string, error)
	SetURL(ctx context.Context, link *core.Link) (string, error)
	GetUserURLS(ctx context.Context, userID string) ([]*core.Link, error)
	SetURLSBatch(ctx context.Context, links []*core.Link) ([]*core.Link, error)
	DeleteURLSBatch(ctx context.Context, ids []*string, userID string) error
	GetStats(ctx context.Context) (*core.Stats, error)
	Ping(ctx context.Context) error
	Update()
	Close() error
}

// Auth struct.
type AUTHService struct {
	storage Storage
	config  *core.Config
}

// Create new Auth struct.
func NewAuth(conf *core.Config, s *Storage) *AUTHService {
	return &AUTHService{
		storage: *s,
		config:  conf,
	}
}

// Create neew user.
func (a *AUTHService) CreateUser() {
	log.Println("HELLO USER CREATE")
}
