package filestore

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/Vasily-van-Zaam/ushortener/internal/core"
)

type Event struct {
	file    *os.File
	scanner *bufio.Scanner
	writer  *bufio.Writer
}

type Store struct {
	Config *core.Config
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

func (s *Store) newCreateFile() (*Event, error) {
	file, err := os.Create(s.Config.Filestore)
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
	data, err := s.newOpenFile()
	if err != nil {
		return "", err
	}
	line := 0
	for data.scanner.Scan() {
		d := strings.Split(data.scanner.Text(), ",")
		if len(d) >= 1 {
			if d[0] == id {
				if d[3] == "true" {
					return "", errors.New("deleted")
				}
				return d[1], nil
			}
		}
		line++
	}
	defer data.file.Close()

	if id == "" {
		return "", errors.New("not Found")
	}
	return "", nil
}
func (s *Store) SetURL(ctx context.Context, link *core.Link) (string, error) {
	data, err := s.newOpenFile()
	if err != nil {
		return "", err
	}
	line := 0
	lastElementID := 0
	for data.scanner.Scan() {
		d := strings.Split(data.scanner.Text(), ",")
		if len(d) >= 1 {
			if d[1] == link.Link {
				url := fmt.Sprint(d[0])
				return url, core.NewErrConflict()
			}
		}
		lastElementID, _ = strconv.Atoi(d[0])
		line++
	}
	_, errWriteData := data.writer.Write([]byte(fmt.Sprint(lastElementID+1, ",", link.Link, ",", link.UUID, ",", link.Deleted)))
	if errWriteData != nil {
		return "", errWriteData
	}
	errWriteByte := data.writer.WriteByte('\n')
	if errWriteByte != nil {
		return "", errWriteByte
	}
	err = data.writer.Flush()
	if err != nil {
		return "", err
	}
	defer data.file.Close()
	url := fmt.Sprint(lastElementID + 1)
	return url, nil
}

func (s *Store) GetUserURLS(ctx context.Context, userID string) ([]*core.Link, error) {
	data, err := s.newOpenFile()
	if err != nil {
		return nil, err
	}

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
	defer data.file.Close()

	if userID == "" {
		return nil, errors.New("not Found")
	}
	return links, nil
}

func scan(data *Event) []*core.Link {
	res := make([]*core.Link, 0, 10)
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

func (s *Store) SetURLSBatch(ctx context.Context, links []*core.Link) ([]*core.Link, error) {
	data, err := s.newOpenFile()
	if err != nil {
		return nil, err
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
			return nil, errWriteData
		}
		errWriteByte := data.writer.WriteByte('\n')
		if errWriteByte != nil {
			return nil, errWriteByte
		}
		err = data.writer.Flush()
		if err != nil {
			return nil, err
		}
		result = append(result, &core.Link{
			ID:      lastElementID + count,
			Link:    l.Link,
			UUID:    l.UUID,
			Deleted: false,
		})
	}
	return result, errConflict
}

func (s *Store) DeleteURLSBatch(ctx context.Context, ids []*string, userID string) error {
	data, err := s.newOpenFile()
	if err != nil {
		return err
	}
	defer data.file.Close()
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

	err = os.WriteFile(s.Config.Filestore, []byte(dataLines), 0644)
	if err != nil {
		return err
	}
	return nil
}
func (s *Store) Close() error {
	return nil
}

func (s *Store) Ping(ctx context.Context) error {
	return nil
}
