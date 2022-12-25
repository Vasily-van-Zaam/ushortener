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

	"github.com/Vasily-van-Zaam/ushortener/internal/core"
)

type Event struct {
	file    *os.File
	scanner *bufio.Scanner
	writer  *bufio.Writer
}

type Filestore struct {
	Config *core.Config
}

func New(conf *core.Config) (*Filestore, error) {
	return &Filestore{
		Config: conf,
	}, nil
}

func (s *Filestore) newOpenFile() (*Event, error) {
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

func (s *Filestore) GetURL(ctx context.Context, id string) (string, error) {
	data, err := s.newOpenFile()
	if err != nil {
		return "", err
	}
	line := 0
	for data.scanner.Scan() {
		d := strings.Split(data.scanner.Text(), ",")
		log.Println(d, line)
		if len(d) >= 1 {
			if d[0] == id {
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
func (s *Filestore) SetURL(ctx context.Context, link *core.Link) (string, error) {
	data, err := s.newOpenFile()
	if err != nil {
		return "", err
	}
	line := 0
	lastElementID := 0
	for data.scanner.Scan() {
		d := strings.Split(data.scanner.Text(), ",")
		log.Println(d, line)
		if len(d) >= 1 {
			if d[1] == link.Link {
				return fmt.Sprint(d[0]), nil
			}
		}
		lastElementID, _ = strconv.Atoi(d[0])
		line++
	}

	if _, errWrite := data.writer.Write([]byte(fmt.Sprint(lastElementID+1, ",", link.Link, ",", link.UserID))); err != nil {
		return "", errWrite
	}
	if errW := data.writer.WriteByte('\n'); errW != nil {
		return "", errW
	}
	err = data.writer.Flush()
	if err != nil {
		return "", err
	}
	defer data.file.Close()
	return fmt.Sprint(lastElementID + 1), nil
}

func (s *Filestore) GetUserURLS(ctx context.Context, userID string) ([]*core.Link, error) {
	data, err := s.newOpenFile()
	if err != nil {
		return nil, err
	}
	links := []*core.Link{}
	line := 0
	for data.scanner.Scan() {
		d := strings.Split(data.scanner.Text(), ",")
		if len(d) != 0 {
			if d[2] == userID {
				id, _ := strconv.Atoi(d[0])
				links = append(links, &core.Link{
					ID:     id,
					Link:   d[1],
					UserID: d[2],
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

func (s *Filestore) Close() error {
	return nil
}
