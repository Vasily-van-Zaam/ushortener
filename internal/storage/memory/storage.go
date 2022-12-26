package memorystore

import (
	"context"
	"errors"
	"fmt"
	"strconv"

	"github.com/Vasily-van-Zaam/ushortener/internal/core"
)

type Store struct {
	Config *core.Config
	Data   []*core.Link
}

func New(conf *core.Config) (*Store, error) {
	return &Store{
		Config: conf,
		Data:   make([]*core.Link, 0),
	}, nil
}

func (s *Store) GetURL(ctx context.Context, id string) (string, error) {
	for _, l := range s.Data {
		idInt, _ := strconv.Atoi(id)
		if l.ID == idInt {
			return l.Link, nil
		}
	}
	if id == "" {
		return "", errors.New("not Found")
	}
	return "", nil
}
func (s *Store) SetURL(ctx context.Context, link *core.Link) (string, error) {
	for _, l := range s.Data {
		if l.Link == link.Link {
			return fmt.Sprint(s.Config.BaseURL, "/", l.ID), nil
		}
	}
	dataLength := len(s.Data) + 1
	d := core.Link{
		ID:   dataLength,
		Link: link.Link,
	}
	s.Data = append(s.Data, &d)
	return fmt.Sprint(d.ID), nil
}

func (s *Store) GetUserURLS(ctx context.Context, userID string) ([]*core.Link, error) {
	links := []*core.Link{}
	for _, l := range s.Data {
		idInt, _ := strconv.Atoi(userID)
		if l.ID == idInt {
			links = append(links, l)
		}
	}
	if userID == "" {
		return nil, errors.New("not Found")
	}
	return links, nil
}

func (s *Store) Close() error {
	return nil
}
func (s *Store) Ping(ctx context.Context) error {
	return nil
}
