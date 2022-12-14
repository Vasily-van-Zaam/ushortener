package filestore

import (
	"bufio"
	"context"
	"os"

	"github.com/Vasily-van-Zaam/ushortener/internal/core"
)

type Filestore struct {
	File    *os.File
	Config  *core.Config
	Scanner *bufio.Scanner
	Writer  *bufio.Writer
}

func New(conf *core.Config) (*Filestore, error) {
	file, err := os.OpenFile(conf.Filestore, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0777)
	if err != nil {
		return nil, err
	}
	return &Filestore{
		File:    file,
		Config:  conf,
		Scanner: bufio.NewScanner(file),
		Writer:  bufio.NewWriter(file),
	}, nil

}

func (s *Filestore) GetURL(ctx context.Context, id string) (string, error) {
	
	return "", nil
}
func (s *Filestore) SetURL(ctx context.Context, link string) (string, error) {
	return "", nil
}

func (f *Filestore) Close() error {
	return f.File.Close()
}
