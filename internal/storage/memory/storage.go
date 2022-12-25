package memorystore

import (
	"context"
	"errors"
	"fmt"
	"strconv"

	"github.com/Vasily-van-Zaam/ushortener/internal/core"
)

type Memorystore struct {
	Config *core.Config
	Data   []*core.Link
}

func New(conf *core.Config) (*Memorystore, error) {
	return &Memorystore{
		Config: conf,
		Data:   []*core.Link{},
	}, nil
}

func (s *Memorystore) GetURL(ctx context.Context, id string) (string, error) {
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
func (s *Memorystore) SetURL(ctx context.Context, link *core.Link) (string, error) {
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

func (s *Memorystore) GetUserURLS(ctx context.Context, userID string) ([]*core.Link, error) {
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

func (s *Memorystore) Close() error {
	return nil
}
