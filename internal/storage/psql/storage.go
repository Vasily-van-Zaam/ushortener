package psql

import (
	"context"
	"errors"
	"fmt"
	"log"

	"github.com/Vasily-van-Zaam/ushortener/internal/core"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Store struct {
	config *core.Config
	db     *pgxpool.Pool // *pgx.Conn
}

func New(conf *core.Config) (*Store, error) {
	ctx := context.Background()

	config, err := pgxpool.ParseConfig(conf.DataBaseDNS)
	if err != nil {
		panic(err)
	}
	config.AfterConnect = func(ctx context.Context, conn *pgx.Conn) error {
		// do something with every new connection
		return nil
	}

	db, err := pgxpool.NewWithConfig(context.Background(), config)

	if err != nil {
		panic(err)
	}

	_, errExecUser := db.Exec(ctx, `
	CREATE TABLE IF NOT EXISTS users(
		id SERIAL PRIMARY KEY,
		name CHAR(255),
		email CHAR(255) UNIQUE,
		phone CHAR(255) UNIQUE
	);`)
	_, errExec := db.Exec(ctx, `
	CREATE TABLE IF NOT EXISTS links(
		id SERIAL PRIMARY KEY,
		uuid CHAR(255),
		link TEXT UNIQUE,
		short_link char(255) UNIQUE,
		user_id INTEGER,
		deleted boolean,
		FOREIGN KEY (user_id) REFERENCES users (id)
	);`)

	if errExecUser != nil {
		log.Println("errExecUser: ", errExecUser.Error())
	}

	if errExec != nil {
		log.Println("errExec: ", errExec.Error())
	}

	return &Store{
		db:     db,
		config: conf,
	}, nil
}

func (s *Store) GetURL(ctx context.Context, id string) (string, error) {
	res := s.db.QueryRow(ctx, `
	SELECT id,uuid,link,deleted  FROM links WHERE id=$1;
	`, id)
	linkDB := core.Link{}
	err := res.Scan(&linkDB.ID, &linkDB.UUID, &linkDB.Link, &linkDB.Deleted)
	if err != nil {
		log.Println("errorSelectSqlLiteGet", err, linkDB)
	}

	if linkDB.ID != 0 {
		if linkDB.Deleted {
			return "", errors.New("deleted")
		}
		return linkDB.Link, nil
	}
	if id == "" {
		return "", errors.New("not Found")
	}
	return "", nil
}

func (s *Store) SetURL(ctx context.Context, link *core.Link) (string, error) {
	var resID any
	searchLink := s.db.QueryRow(ctx, `
	SELECT id,uuid,link,user_id FROM links WHERE link=$1;
	`, link.Link)

	linkDB := core.Link{}

	err := searchLink.Scan(&linkDB.ID, &linkDB.UUID, &linkDB.Link, &linkDB.UserID)
	if err != nil {
		log.Println("errorSelectSqlLitePost", err, linkDB)
	}

	if linkDB.ID != 0 {
		url := fmt.Sprint(linkDB.ID)
		return url, core.NewErrConflict()
	}

	tx, err := s.db.Begin(ctx)
	if err != nil {
		return "", err
	}
	defer func() {
		errRllback := tx.Rollback(ctx)
		if errRllback != nil {
			log.Print("errRllback:", errRllback)
		}
	}()

	errInsert := tx.QueryRow(ctx, `
	INSERT INTO links (link, uuid) VALUES ($1, $2) RETURNING id;
	`, link.Link, link.UUID,
	).Scan(&linkDB.ID)
	if errInsert != nil {
		log.Println("errInsert:", errInsert)
	}
	resID = linkDB.ID

	errCommit := tx.Commit(ctx)
	if errCommit != nil {
		log.Println("errCommit:", errCommit)
	}
	url := fmt.Sprint(resID)
	return url, nil
}

func (s *Store) GetUserURLS(ctx context.Context, userID string) ([]*core.Link, error) {
	// userID = "4f744217-a3cb-4bad-9c76-6880e41d336f"
	query, err := s.db.Query(ctx, `
	SELECT id, link FROM links WHERE uuid=$1
	`, userID)
	if err != nil {
		log.Println("error query", err)
	}
	defer query.Close()
	res := make([]*core.Link, 0, 10)
	for query.Next() {
		linkDB := &core.Link{}
		errScan := query.Scan(&linkDB.ID, &linkDB.Link)
		if errScan != nil {
			log.Println(errScan)
			// return nil, errScan
		}
		res = append(res, linkDB)
	}
	return res, nil
}

func (s *Store) SetURLSBatch(ctx context.Context, links []*core.Link) ([]*core.Link, error) {
	var errConflict *core.ErrConflict
	response := make([]*core.Link, 0)
	tx, err := s.db.Begin(ctx)
	if err != nil {
		return nil, err
	}
	defer func() {
		errRllback := tx.Rollback(ctx)
		if errRllback != nil {
			log.Print("errRllback:", errRllback)
		}
	}()

	for _, l := range links {
		searchLink := tx.QueryRow(ctx, `
			SELECT id,uuid,link,user_id,deleted FROM links WHERE link=$1;
		`, l.Link)

		linkDB := core.Link{}
		err = searchLink.Scan(&linkDB.ID, &linkDB.UUID, &linkDB.Link, &linkDB.UserID, &linkDB.Deleted)
		if err != nil {
			log.Println("errorSelectSqlLitePost", err, linkDB)
		}
		if linkDB.ID != 0 {
			errConflict = core.NewErrConflict()
			response = append(response, &linkDB)
		} else {
			errInsert := tx.QueryRow(ctx, `
				INSERT INTO links (link,uuid) VALUES ($1,$2) RETURNING id;
			`, l.Link, l.UUID,
			).Scan(&linkDB.ID)
			if errInsert != nil {
				log.Println("errInsert:", errInsert)
			}
			response = append(response, &linkDB)
		}
	}
	errCommit := tx.Commit(ctx)
	if errCommit != nil {
		log.Println("errCommit:", errCommit)
		return nil, errCommit
	}
	return response, errConflict
}

func (s *Store) DeleteURLSBatch(ctx context.Context, ids []*string, userID string) error {
	ctx = context.Background()
	tx, err := s.db.Begin(ctx)
	if err != nil {
		return err
	}

	defer func() {
		errRllback := tx.Rollback(ctx)
		if errRllback != nil {
			log.Print("errRllback:", errRllback)
		}
	}()

	listIds := ""
	for i, id := range ids {
		if i == len(ids)-1 {
			listIds += *id
		} else {
			listIds += *id + ","
		}
	}

	_, err = tx.Exec(ctx, `
		update links
		set deleted = true
		where id in (`+listIds+`) and uuid = $1`,
		userID,
	)

	if err != nil {
		log.Println("errExec:", err)
		return err
	}

	err = tx.Commit(ctx)
	if err != nil {
		log.Println("errCommit:", err)
		return err
	}

	return nil
}

func (s *Store) Close() error {
	s.db.Close()
	return nil // s.db.Close(context.Background())
}

func (s *Store) Ping(ctx context.Context) error {
	err := s.db.Ping(ctx) // PgConn().CheckConn()
	if err != nil {
		return err
	}
	return nil
}

func (s *Store) Update() {}
