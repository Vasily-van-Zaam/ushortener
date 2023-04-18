// Memory store
package memorystore

import (
	"context"
	"errors"
	"fmt"
	"strconv"

	"github.com/Vasily-van-Zaam/ushortener/internal/core"
)

// Main structure.
type Store struct {
	Config *core.Config
	Data   []*core.Link
}

// Creeate new store.
func New(conf *core.Config) (*Store, error) {
	return &Store{
		Config: conf,
		Data:   make([]*core.Link, 0),
	}, nil
}

// Get url.
func (s *Store) GetURL(ctx context.Context, id string) (string, error) {
	for _, l := range s.Data {
		idInt, _ := strconv.Atoi(id)
		if l.ID == idInt {
			if l.Deleted {
				return "", errors.New("deleted")
			}
			return l.Link, nil
		}
	}
	if id == "" {
		return "", errors.New("not Found")
	}
	return "", nil
}

// Sset url.
func (s *Store) SetURL(ctx context.Context, link *core.Link) (string, error) {
	for _, l := range s.Data {
		if l.Link == link.Link {
			url := fmt.Sprint(l.ID)
			return url, core.NewErrConflict()
		}
	}
	dataLength := len(s.Data) + 1
	d := core.Link{
		ID:      dataLength,
		Link:    link.Link,
		UUID:    link.UUID,
		Deleted: false,
	}
	s.Data = append(s.Data, &d)
	url := fmt.Sprint(d.ID)
	return url, nil
}

// Get list user urls.
func (s *Store) GetUserURLS(ctx context.Context, userID string) ([]*core.Link, error) {
	links := make([]*core.Link, 0, 10)
	for _, l := range s.Data {
		if l.UUID == userID {
			links = append(links, l)
		}
	}
	if userID == "" {
		return nil, errors.New("not Found")
	}
	return links, nil
}

// Set list user urls.
func (s *Store) SetURLSBatch(ctx context.Context, links []*core.Link) ([]*core.Link, error) {
	result := make([]*core.Link, 0, 10)
	var errConflict *core.ErrConflict
	for _, l := range links {
		id := len(s.Data) + 1
		exists := false
		for _, ls := range s.Data {
			if ls.Link == l.Link {
				l.ID = ls.ID
				exists = true
				errConflict = core.NewErrConflict()
				break
			}
		}
		if !exists {
			l.ID = id
			s.Data = append(s.Data, l)
		}
		result = append(result, l)
	}

	return result, errConflict
}
func (s *Store) Close() error {
	return nil
}
func (s *Store) Ping(ctx context.Context) error {
	return nil
}

// Delete list urls by list id.
func (s *Store) DeleteURLSBatch(ctx context.Context, ids []*string, userID string) error {
	for _, l := range s.Data {
		for _, idStr := range ids {
			id, _ := strconv.Atoi(*idStr)
			if l.ID == id && l.UUID == userID {
				l.Deleted = true
			}
		}
	}
	return nil
}

func (s *Store) Update() {}
