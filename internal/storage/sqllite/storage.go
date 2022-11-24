package litestore

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"

	"github.com/Vasily-van-Zaam/ushortener/internal/core"
	_ "github.com/mattn/go-sqlite3"
)

type SomeStorage struct {
	db *sql.DB
}

func New(db *sql.DB) *SomeStorage {

	_, errExec := db.Exec(`
	CREATE TABLE IF NOT EXISTS links(
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		link TEXT UNIQUE,
		short_link char(255) UNIQUE
	);`)
	if errExec != nil {
		log.Println("errExec: ", errExec.Error())
	}

	return &SomeStorage{
		db: db,
	}
}
func (s *SomeStorage) GetURL(ctx context.Context, id string) (string, error) {

	res := s.db.QueryRowContext(ctx, `
	SELECT * FROM links WHERE id=$1;
	`, id)
	linkDB := core.Link{}
	err := res.Scan(&linkDB.ID, &linkDB.Link, &linkDB.ShortLink)
	if err != nil {
		log.Println("errorSelectSqlLiteGet", err, linkDB)
	}
	if linkDB.ID == 0 {
		return "", errors.New("not Found")
	}
	return fmt.Sprint(linkDB.Link), nil
}
func (s *SomeStorage) SetURL(ctx context.Context, link string) (string, error) {
	// symbols := []string{
	// 	"0", "1", "2", "3", "4", "5", "6", "7", "8", "9",
	// 	"a", "A", "b", "B", "c", "C", "d", "D", "e", "E", "f", "F", "g", "G",
	// 	"h", "H", "i", "I", "j", "J", "k", "K", "l", "L", "m", "M",
	// 	"n", "N", "o", "O", "p", "P", "q", "Q", "r", "R", "s", "S", "t", "T",
	// 	"u", "U", "v", "V", "x", "X", "w", "W", "y", "Y", "z",
	// }

	var resID any

	searchLink := s.db.QueryRowContext(ctx, `
	SELECT * FROM links WHERE link=$1;
	`, link)

	linkDB := core.Link{}

	err := searchLink.Scan(&linkDB.ID, &linkDB.Link, &linkDB.ShortLink)
	if err != nil {
		log.Println("errorSelectSqlLitePost", err, linkDB)
	}

	if linkDB.ID != 0 {
		return fmt.Sprint(core.MAINDOMAIN, linkDB.ID), nil
	}
	res, err := s.db.ExecContext(ctx, `
	INSERT INTO links (link) VALUES ($link);
	`, link)
	if err != nil {
		log.Println("errorInsertSqlLitePost", err)
	}
	resID, _ = res.LastInsertId()

	return fmt.Sprint(core.MAINDOMAIN, resID), nil
}