package some_storage

import (
	"context"
	"log"
)

type SomeStorage struct{}

func New() *SomeStorage {
	return &SomeStorage{}
}

func (s *SomeStorage) GetUrl(ctx context.Context, link string) (string, error) {
	return "", nil
}
func (s *SomeStorage) SetUrl(ctx context.Context, link string) (string, error) {
	symbols := []string{
		"0", "1", "2", "3", "4", "5", "6", "7", "8", "9",
		"a", "A", "b", "B", "c", "C", "d", "D", "e", "E", "f", "F", "g", "G",
		"h", "H", "i", "I", "j", "J", "k", "K", "l", "L", "m", "M",
		"n", "N", "o", "O", "p", "P", "q", "Q", "r", "R", "s", "S", "t", "T",
		"u", "U", "v", "V", "x", "X", "w", "W", "y", "Y", "z",
	}
	log.Println(symbols)

	return "", nil
}
