package filestore

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"sync"

	"github.com/Vasily-van-Zaam/ushortener/internal/core"
)

type Event struct {
	file    *os.File
	scanner *bufio.Scanner
	writer  *bufio.Writer
}

type Store struct {
	Config *core.Config
	Data   []*core.Link
	saved  chan any
}

func (s *Store) newOpenFile() (*Event, error) {
	file, err := os.OpenFile(s.Config.Filestore, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0777)

	if err != nil {
		return nil, err
	}
	return &Event{
		file:    file,
		scanner: bufio.NewScanner(file),
		writer:  bufio.NewWriter(file),
	}, nil
}

func scan(data *Event) []*core.Link {
	res := make([]*core.Link, 0)
	line := 0
	lastElementID := 0
	for data.scanner.Scan() {
		d := strings.Split(data.scanner.Text(), ",")
		lastElementID, _ = strconv.Atoi(d[0])
		deleted := false
		if d[3] == "true" {
			deleted = true
		}
		if len(d) >= 1 {
			res = append(res, &core.Link{
				ID:      lastElementID,
				UUID:    d[2],
				Link:    d[1],
				Deleted: deleted,
			})
		}
		line++
	}
	return res
}

func New(conf *core.Config) (*Store, error) {
	s := &Store{
		Config: conf,
		Data:   make([]*core.Link, 0),
		saved:  make(chan any),
	}

	e, err := s.newOpenFile()
	if err != nil {
		return nil, err
	}
	defer e.file.Close()
	s.Data = scan(e)
	return s, nil
}

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
	go func() {
		s.saved <- nil
	}()
	return url, nil
}

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
	go func() {
		s.saved <- nil
	}()
	return result, errConflict
}

func (s *Store) Update() {
	for save := range s.saved {
		wg := &sync.WaitGroup{}

		wg.Add(1)
		go func() {
			defer wg.Done()
			log.Println("SAVE", save)
			errorSaved := make(chan *error)
			okSaved := make(chan any)
			dataLines := ""
			for _, l := range s.Data {
				dataLines += fmt.Sprint(l.ID, ",", l.Link, ",", l.UUID, ",", l.Deleted, "\n")
			}
			go func() {
				err := os.WriteFile(s.Config.Filestore, []byte(dataLines), 0600)
				if err != nil {
					errorSaved <- &err
					return
				}
				okSaved <- nil
			}()
			select {
			case err := <-errorSaved:
				log.Println("ERROR SAVED", err)
				return
			case ok := <-okSaved:
				log.Println("SAVED OK", ok)

				return
			}
		}()
		wg.Wait()
	}
}
func (s *Store) Close() error {
	return nil
}
func (s *Store) Ping(ctx context.Context) error {
	return nil
}

func (s *Store) DeleteURLSBatch(ctx context.Context, ids []*string, userID string) error {
	for _, l := range s.Data {
		for _, idStr := range ids {
			id, _ := strconv.Atoi(*idStr)
			if l.ID == id && l.UUID == userID {
				l.Deleted = true
			}
		}
	}
	go func() {
		s.saved <- nil
	}()

	return nil
}
