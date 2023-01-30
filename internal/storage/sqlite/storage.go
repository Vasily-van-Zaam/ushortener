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

type Store struct {
	db     *sql.DB
	Config *core.Config
}

func New(conf *core.Config) (*Store, error) {
	db, err := sql.Open("sqlite3", conf.SqliteDB)
	if err != nil {
		panic(err)
	}
	_, errExec := db.Exec(`
	CREATE TABLE IF NOT EXISTS links(
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		uuid CHAR(255),
		link TEXT UNIQUE,
		short_link char(255) UNIQUE
	);`)
	if errExec != nil {
		log.Println("errExec: ", errExec.Error())
	}

	return &Store{
		db:     db,
		Config: conf,
	}, nil
}
func (s *Store) GetURL(ctx context.Context, id string) (string, error) {
	res := s.db.QueryRowContext(ctx, `
	SELECT * FROM links WHERE id=$1;
	`, id)
	linkDB := core.Link{}
	err := res.Scan(&linkDB.ID, &linkDB.UUID, &linkDB.Link, &linkDB.ShortLink)
	if err != nil {
		log.Println("errorSelectSqlLiteGet", err, linkDB)
	}

	if linkDB.ID != 0 {
		return linkDB.Link, nil
	}
	if id == "" {
		return "", errors.New("not Found")
	}
	return "", nil
}
func (s *Store) SetURL(ctx context.Context, link *core.Link) (string, error) {
	var resID any
	searchLink := s.db.QueryRowContext(ctx, `
	SELECT * FROM links WHERE link=$1;
	`, link.Link)

	linkDB := core.Link{}

	err := searchLink.Scan(&linkDB.ID, &linkDB.Link, &linkDB.ShortLink, &linkDB.UUID)
	if err != nil {
		log.Println("errorSelectSqlLitePost", err, linkDB)
	}

	if linkDB.ID != 0 {
		url := fmt.Sprint(linkDB.ID)
		return url, nil
	}
	res, err := s.db.ExecContext(ctx, `
	INSERT INTO links (link, uuid) VALUES ($1, $2);
	`, link.Link, link.UUID)
	if err != nil {
		log.Println("errorInsertSqlLitePost", err)
	}
	resID, _ = res.LastInsertId()
	url := fmt.Sprint(resID)
	return url, nil
}

func (s *Store) GetUserURLS(ctx context.Context, userID string) ([]*core.Link, error) {
	/// TODO Пока не работает допилить запрос
	links := make([]*core.Link, 0)
	return links, nil
}

func (s *Store) SetURLSBatch(ctx context.Context, links []*core.Link) ([]*core.Link, error) {
	return nil, nil
}

func (s *Store) DeleteURLSBatch(ctx context.Context, ids []*string, userID string) error {
	return nil
}

func (s *Store) Close() error {
	return s.db.Close()
}

func (s *Store) Ping(ctx context.Context) error {
	return nil
}
