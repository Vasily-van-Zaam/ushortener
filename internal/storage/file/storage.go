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
		if len(d) >= 1 && d[0] == id {
			return d[1], nil
		}
		line++
	}
	defer data.file.Close()
	return "", errors.New("not Found")

}
func (s *Filestore) SetURL(ctx context.Context, link string) (string, error) {

	data, err := s.newOpenFile()
	if err != nil {
		return "", err
	}
	line := 0
	lastElementId := 0
	for data.scanner.Scan() {
		d := strings.Split(data.scanner.Text(), ",")
		log.Println(d, line)
		if len(d) >= 1 && d[1] == link {
			return fmt.Sprint(s.Config.BaseURL, "/", d[0]), nil
		}
		lastElementId, _ = strconv.Atoi(d[0])
		line++
	}

	if _, err := data.writer.Write([]byte(fmt.Sprint(lastElementId+1, ",", link))); err != nil {
		return "", err
	}
	if err := data.writer.WriteByte('\n'); err != nil {
		return "", err
	}
	err = data.writer.Flush()
	if err != nil {
		return "", err
	}
	defer data.file.Close()
	return fmt.Sprint(s.Config.BaseURL, "/", lastElementId+1), nil
}

func (s *Filestore) Close() error {
	return nil
}
