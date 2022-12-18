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
	return s.Config.BaseURL, nil
}
func (s *Memorystore) SetURL(ctx context.Context, link string) (string, error) {
	for _, l := range s.Data {
		if l.Link == link {
			return fmt.Sprint(s.Config.BaseURL, "/", l.ID), nil
		}
	}
	dataLength := len(s.Data) + 1
	d := core.Link{
		ID:   dataLength,
		Link: link,
	}
	s.Data = append(s.Data, &d)
	return fmt.Sprint(s.Config.BaseURL, "/", d.ID), nil
}

func (s *Memorystore) Close() error {
	return nil
}
