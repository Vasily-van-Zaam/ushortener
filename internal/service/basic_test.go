package service_test

import (
	"context"
	"testing"

	"github.com/Vasily-van-Zaam/ushortener/internal/core"
	"github.com/Vasily-van-Zaam/ushortener/internal/service"
)

type MockStore struct {
}

func (s *MockStore) GetURL(ctx context.Context, id string) (string, error) {
	return "", nil
}
func (s *MockStore) SetURL(ctx context.Context, link *core.Link) (string, error) {
	return "", nil
}
func (s *MockStore) GetUserURLS(ctx context.Context, userID string) ([]*core.Link, error) {
	return nil, nil
}
func (s *MockStore) SetURLSBatch(ctx context.Context, links []*core.Link) ([]*core.Link, error) {
	return nil, nil
}
func (s *MockStore) DeleteURLSBatch(ctx context.Context, ids []*string, userID string) error {
	return nil
}
func (s *MockStore) Ping(ctx context.Context) error {
	return nil
}
func (s *MockStore) Update() {}
func (s *MockStore) Close() error {
	return nil
}

func TestBasicService_GetURL(t *testing.T) {
	type fields struct {
		storage     service.Storage
		config      *core.Config
		AUTHService core.AUTHService
	}
	type args struct {
		ctx context.Context
		id  string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    string
		wantErr bool
	}{
		{
			fields: fields{
				storage: &MockStore{},
			},
		},
		{
			fields: fields{
				storage: &MockStore{},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := service.NewBasic(tt.fields.config, tt.fields.storage, nil)
			got, err := s.GetURL(tt.args.ctx, tt.args.id)
			if (err != nil) != tt.wantErr {
				t.Errorf("BasicService.GetURL() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("BasicService.GetURL() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestBasicService_SetURL(t *testing.T) {
	type fields struct {
		storage     service.Storage
		config      *core.Config
		AUTHService core.AUTHService
	}
	type args struct {
		ctx  context.Context
		link string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    string
		wantErr bool
	}{
		{
			fields: fields{
				storage: &MockStore{},
				config: &core.Config{
					BaseURL: "test.loc",
				},
			},
			args: args{
				ctx: context.WithValue(
					context.Background(),
					core.USERDATA,
					core.User{
						ID: "1",
					},
				),
				link: "/",
			},
			want: "test.loc/",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := service.NewBasic(tt.fields.config, tt.fields.storage, nil)

			got, err := s.SetURL(tt.args.ctx, tt.args.link)
			if (err != nil) != tt.wantErr {
				t.Errorf("BasicService.SetURL() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("BasicService.SetURL() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestBasicService_Ping(t *testing.T) {
	type fields struct {
		storage     service.Storage
		config      *core.Config
		AUTHService core.AUTHService
	}
	type args struct {
		ctx context.Context
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			fields: fields{
				storage: &MockStore{},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := service.NewBasic(tt.fields.config, tt.fields.storage, nil)
			if err := s.Ping(tt.args.ctx); (err != nil) != tt.wantErr {
				t.Errorf("BasicService.Ping() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
