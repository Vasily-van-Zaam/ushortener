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
	mu     sync.RWMutex
}

func New(conf *core.Config) (*Store, error) {
	return &Store{
		Config: conf,
	}, nil
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

func (s *Store) GetURL(ctx context.Context, id string) (string, error) {
	resultOk := make(chan *string)
	resultError := make(chan *error)
	if id == "" {
		return "", errors.New("not Found")
	}
	go func() {
		s.mu.RLock()
		defer s.mu.RUnlock()

		data, err := s.newOpenFile()
		if err != nil {
			resultError <- &err
			return
		}
		defer data.file.Close()
		line := 0
		for data.scanner.Scan() {
			d := strings.Split(data.scanner.Text(), ",")
			if len(d) >= 1 {
				if d[0] == id {
					if d[3] == "true" {
						err := errors.New("deleted")
						resultError <- &err
						return
					}
					resultOk <- &d[1]
					return
				}
			}
			line++
		}
	}()

	select {
	case res := <-resultOk:
		return *res, nil
	case err := <-resultError:
		return "", *err
	}
}

type OkErr struct {
	res string
	err error
}

func (s *Store) SetURL(ctx context.Context, link *core.Link) (string, error) {
	resultOk := make(chan *string)
	resultError := make(chan *error)
	resConflick := make(chan *OkErr)

	go func() {
		s.mu.Lock()
		defer s.mu.Unlock()

		data, err := s.newOpenFile()
		if err != nil {
			resultError <- &err
			return
		}
		line := 0
		lastElementID := 0
		for data.scanner.Scan() {
			d := strings.Split(data.scanner.Text(), ",")
			if len(d) >= 1 {
				if d[1] == link.Link {
					url := fmt.Sprint(d[0])
					resConflick <- &OkErr{url, core.NewErrConflict()}
					return
				}
			}
			lastElementID, _ = strconv.Atoi(d[0])
			line++
		}
		_, errWriteData := data.writer.Write([]byte(fmt.Sprint(lastElementID+1, ",", link.Link, ",", link.UUID, ",", link.Deleted)))
		if errWriteData != nil {
			resultError <- &errWriteData
			return
		}
		errWriteByte := data.writer.WriteByte('\n')
		if errWriteByte != nil {
			resultError <- &errWriteByte
			return
		}
		err = data.writer.Flush()
		if err != nil {
			resultError <- &err
			return
		}
		defer data.file.Close()
		url := fmt.Sprint(lastElementID + 1)
		resultOk <- &url
		return
	}()

	select {
	case res := <-resultOk:
		return *res, nil
	case err := <-resultError:
		return "", *err
	case conflict := <-resConflick:
		return conflict.res, conflict.err
	}
}

func (s *Store) GetUserURLS(ctx context.Context, userID string) ([]*core.Link, error) {
	resultOk := make(chan []*core.Link)
	resultError := make(chan *error)

	go func() {
		s.mu.RLock()
		defer s.mu.RUnlock()

		data, err := s.newOpenFile()
		if err != nil {
			resultError <- &err
			return
		}
		defer data.file.Close()
		links := make([]*core.Link, 0, 10)
		line := 0
		for data.scanner.Scan() {
			d := strings.Split(data.scanner.Text(), ",")
			if len(d) != 0 {
				if d[2] == userID {
					id, _ := strconv.Atoi(d[0])
					deleted := false
					if d[3] == "true" {
						deleted = true
					}
					links = append(links, &core.Link{
						ID:      id,
						Link:    d[1],
						UUID:    d[2],
						Deleted: deleted,
					})
				}
			}
			line++
		}

		if userID == "" {
			err := errors.New("not Found")
			resultError <- &err
			return
		}
		resultOk <- links
		return
	}()

	select {
	case res := <-resultOk:
		return res, nil
	case err := <-resultError:
		return nil, *err
	}
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

type OkConflictLin struct {
	links []*core.Link
	err   error
}

func (s *Store) SetURLSBatch(ctx context.Context, links []*core.Link) ([]*core.Link, error) {
	resultError := make(chan *error)
	resOkConflict := make(chan *OkConflictLin)

	go func() {
		s.mu.Lock()
		defer s.mu.Unlock()
		data, err := s.newOpenFile()
		if err != nil {
			resultError <- &err
			return
		}
		defer data.file.Close()

		dataList := scan(data)
		result := make([]*core.Link, 0)
		count := 0
		var errConflict *core.ErrConflict
		for _, l := range links {
			exists := false
			lastElementID := 0
			for _, ls := range dataList {
				if ls.Link == l.Link {
					exists = true
					result = append(result, ls)
				}
				lastElementID = ls.ID
			}

			if exists {
				errConflict = core.NewErrConflict()
				continue
			}
			count++
			_, errWriteData :=
				data.writer.Write([]byte(fmt.Sprint(lastElementID+count, ",", l.Link, ",", l.UUID, ",", l.Deleted)))
			if errWriteData != nil {
				resultError <- &errWriteData
				return
			}
			errWriteByte := data.writer.WriteByte('\n')
			if errWriteByte != nil {
				resultError <- &errWriteByte
				return
			}
			err = data.writer.Flush()
			if err != nil {
				resultError <- &err
				return
			}
			result = append(result, &core.Link{
				ID:      lastElementID + count,
				Link:    l.Link,
				UUID:    l.UUID,
				Deleted: false,
			})
		}
		resOkConflict <- &OkConflictLin{
			links: result,
			err:   errConflict,
		}
		return
	}()

	select {
	case ok := <-resOkConflict:
		return ok.links, ok.err
	case err := <-resultError:
		return nil, *err
	}
}

func (s *Store) DeleteURLSBatch(ctx context.Context, ids []*string, userID string) error {
	resultOk := make(chan any)
	resultError := make(chan *error)

	go func() {
		s.mu.Lock()
		defer s.mu.Unlock()
		data, err := s.newOpenFile()
		if err != nil {
			resultError <- &err
			return
		}

		dataLines := ""
		for data.scanner.Scan() {
			d := strings.Split(data.scanner.Text(), ",")
			exists := false
			for _, id := range ids {
				if d[0] == *id && d[2] == userID {
					exists = true
				}
			}
			dataLines += fmt.Sprint(d[0], ",", d[1], ",", d[2], ",", exists, "\n")
		}
		data.file.Close()
		err = os.WriteFile(s.Config.Filestore, []byte(dataLines), 0644)
		if err != nil {
			resultError <- &err
			return
		}
		resultOk <- nil
		return
	}()

	select {
	case ok := <-resultOk:
		log.Println(ok)
		return nil
	case err := <-resultError:
		return *err		
	}
}
func (s *Store) Close() error {
	return nil
}

func (s *Store) Ping(ctx context.Context) error {
	return nil
}
