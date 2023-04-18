package service

import (
	"context"
	"errors"
	"fmt"
	"log"
	"sync"

	"github.com/Vasily-van-Zaam/ushortener/internal/core"
	"github.com/Vasily-van-Zaam/ushortener/pkg/shorter"
)

var BUF chan *core.BuferDeleteURL = make(chan *core.BuferDeleteURL, 1000)

type API struct {
	storage Storage
	config  *core.Config
	shorter iShorter
	core.AUTHService
}

func NewAPI(conf *core.Config, s *Storage, auth *AUTHService) *API {
	return &API{
		*s,
		conf,
		shorter.NewShorter59(),
		auth,
	}
}
func (s *API) BindBuferIds() {
	// defer s.BindBuferIds()
	for b := range BUF {
		wg := &sync.WaitGroup{}
		buf := *b
		wg.Add(1)
		go func() {
			defer wg.Done()
			ids := buf.UnConvertIDS()
			err := s.storage.DeleteURLSBatch(buf.Ctx, ids, buf.User.ID)
			if err != nil {
				log.Println("DELETED ERROR", err)
			}
			log.Println("DELETED OK")
		}()
		wg.Wait()
	}
}

func (s *API) APISetShorten(ctx context.Context, request *core.RequestAPIShorten) (*core.ResponseAPIShorten, error) {
	user := core.User{}

	err := user.SetUserIDFromContext(ctx)
	if err != nil {
		return nil, err
	}
	l := core.Link{
		Link:    request.URL,
		UUID:    user.ID,
		Deleted: false,
	}
	res, err := s.storage.SetURL(ctx, &l)
	if err != nil && !errors.Is(err, core.NewErrConflict()) {
		return nil, err
	}
	return &core.ResponseAPIShorten{
		Result: fmt.Sprint(s.config.BaseURL+"/", s.shorter.Convert(res)),
	}, err
}

func (s *API) APIGetUserURLS(ctx context.Context) ([]*core.ResponseAPIUserURL, error) {
	user := core.User{}
	err := user.SetUserIDFromContext(ctx)
	if err != nil {
		return nil, err
	}
	res, err := s.storage.GetUserURLS(ctx, user.ID)
	if err != nil {
		log.Println("ERROR: GetUserURLS", err)
	}

	resAPI := make([]*core.ResponseAPIUserURL, 0, 10)
	for _, r := range res {
		if r != nil {
			resAPI = append(resAPI, &core.ResponseAPIUserURL{
				ShortURL:    fmt.Sprint(s.config.BaseURL, "/", r.ConverID()),
				OriginalURL: r.Link,
			})
		}
	}

	return resAPI, err
}

func (s *API) APISetShortenBatch(ctx context.Context, request []*core.RequestAPIShortenBatch) ([]*core.ResponseAPIShortenBatch, error) {
	user := core.User{}
	err := user.SetUserIDFromContext(ctx)
	if err != nil {
		return nil, err
	}
	links := make([]*core.Link, 0, 10)
	for _, r := range request {
		links = append(links, &core.Link{
			UUID:    user.ID,
			Link:    r.OriginalURL,
			Deleted: false,
		})
	}
	res := make([]*core.ResponseAPIShortenBatch, 0, 10)
	resDB, err := s.storage.SetURLSBatch(ctx, links)

	for i, r := range resDB {
		res = append(res, &core.ResponseAPIShortenBatch{
			CorrelationID: request[i].CorrelationID,
			ShortURL:      fmt.Sprint(s.config.BaseURL, "/", r.ConverID()),
		})
	}

	return res, err
}

func (s *API) APIDeleteUserURLS(ctx context.Context, ids []*string) error {
	var user core.User
	err := user.SetUserIDFromContext(ctx)
	if err != nil {
		return err
	}
	BUF <- &core.BuferDeleteURL{
		User: &user,
		Ctx:  ctx,
		IDS:  ids,
	}
	return nil
}
