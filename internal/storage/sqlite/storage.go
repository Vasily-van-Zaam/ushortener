package sqlite

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"

	"github.com/Vasily-van-Zaam/ushortener/internal/core"
	_ "github.com/mattn/go-sqlite3" // golint
)

type Sqlitestore struct {
	db     *sql.DB
	Config *core.Config
}

func New(conf *core.Config) (*Sqlitestore, error) {
	db, err := sql.Open("sqlite3", conf.SqliteDB)
	if err != nil {
		panic(err)
	}
	_, errExec := db.Exec(`
	CREATE TABLE IF NOT EXISTS links(
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		link TEXT UNIQUE,
		short_link char(255) UNIQUE
	);`)
	if errExec != nil {
		log.Println("errExec: ", errExec.Error())
	}

	return &Sqlitestore{
		db:     db,
		Config: conf,
	}, nil
}
func (s *Sqlitestore) GetURL(ctx context.Context, id string) (string, error) {
	res := s.db.QueryRowContext(ctx, `
	SELECT * FROM links WHERE id=$1;
	`, id)
	linkDB := core.Link{}
	err := res.Scan(&linkDB.ID, &linkDB.Link, &linkDB.ShortLink)
	if err != nil {
		log.Println("errorSelectSqlLiteGet", err, linkDB)
	}
	if linkDB.ID != 0 {
		return fmt.Sprint(linkDB.Link), nil

		// return "", errors.New("not Found")
	}
	if id == "" {
		return "", errors.New("not Found")
	}
	return s.Config.BaseURL, nil
}
func (s *Sqlitestore) SetURL(ctx context.Context, link string) (string, error) {
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
		return fmt.Sprint(s.Config.BaseURL, "/", linkDB.ID), nil
	}
	res, err := s.db.ExecContext(ctx, `
	INSERT INTO links (link) VALUES ($link);
	`, link)
	if err != nil {
		log.Println("errorInsertSqlLitePost", err)
	}
	resID, _ = res.LastInsertId()

	return fmt.Sprint(s.Config.BaseURL, "/", resID), nil
}

func (s *Sqlitestore) Close() error {
	return s.db.Close()
}
