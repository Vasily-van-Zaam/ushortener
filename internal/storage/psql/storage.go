package psql

import (
	"context"

	"github.com/Vasily-van-Zaam/ushortener/internal/core"
	"github.com/jackc/pgx/v5"
)

type Store struct {
	config *core.Config
	db     *pgx.Conn
}

func New(conf *core.Config) (*Store, error) {
	db, err := pgx.Connect(context.Background(), conf.DataBaseDNS)
	if err != nil {
		panic(err)
	}
	return &Store{
		db:     db,
		config: conf,
	}, nil
}

func (s *Store) GetURL(ctx context.Context, id string) (string, error) {
	return "", nil
}

func (s *Store) SetURL(ctx context.Context, link *core.Link) (string, error) {
	return "", nil
}

func (s *Store) GetUserURLS(ctx context.Context, userID string) ([]*core.Link, error) {
	return []*core.Link{}, nil
}

func (s *Store) Close() error {
	return s.db.Close(context.Background())
}

func (s *Store) Ping(ctx context.Context) error {
	err := s.db.PgConn().CheckConn()
	if err != nil {
		return err
	}
	return nil
}
