// Serveice API
package service

import (
	"context"
	"net/http"
	"net/url"
	"reflect"
	"testing"

	"github.com/Vasily-van-Zaam/ushortener/internal/core"
)

type mockAPIStore struct {
}

func newTestStore() *mockAPIStore {
	return &mockAPIStore{}
}

func (s *mockAPIStore) GetStats(ctx context.Context) (*core.Stats, error) {
	return &core.Stats{
		Urls:  100011,
		Users: 10,
	}, nil
}
func (s *mockAPIStore) GetURL(ctx context.Context, id string) (string, error) {
	return "", nil
}
func (s *mockAPIStore) Close() error {
	return nil
}
func (s *mockAPIStore) SetURL(ctx context.Context, link *core.Link) (string, error) {
	return "", nil
}
func (s *mockAPIStore) GetUserURLS(ctx context.Context, userID string) ([]*core.Link, error) {
	return nil, nil
}
func (s *mockAPIStore) SetURLSBatch(ctx context.Context, links []*core.Link) ([]*core.Link, error) {
	return nil, nil
}
func (s *mockAPIStore) DeleteURLSBatch(ctx context.Context, ids []*string, userID string) error {
	return nil
}
func (s *mockAPIStore) Ping(ctx context.Context) error {
	return nil
}
func (s *mockAPIStore) Update() {
}

func TestAPI_APIGetStats(t *testing.T) {
	type fields struct {
		storage     Storage
		config      *core.Config
		AUTHService *AUTHService
	}
	type args struct {
		r *http.Request
	}
	haderXReal := http.Header{}
	haderXReal.Add("X-Real-IP", "203.0.113.195,70.41.3.18,150.172.238.178")

	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *core.Stats
		wantErr bool
	}{
		{
			fields: fields{
				config: &core.Config{
					TrustedSubnet: "203.0.113.0/24",
				},
				storage: newTestStore(),
			},
			args: args{
				r: &http.Request{
					Method: http.MethodGet,
					URL:    &url.URL{Host: "example.com"},
					Header: haderXReal,
				},
			},
			want: &core.Stats{
				Urls:  100011,
				Users: 10,
			},
		},
		{
			fields: fields{
				config: &core.Config{
					TrustedSubnet: "203.1.113.0/24",
				},
				storage: newTestStore(),
			},
			args: args{
				r: &http.Request{
					Method: http.MethodGet,
					URL:    &url.URL{Host: "example.com"},
					Header: haderXReal,
				},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := NewAPI(tt.fields.config, &tt.fields.storage, tt.fields.AUTHService)

			got, err := a.APIGetStats(tt.args.r)
			if (err != nil) != tt.wantErr {
				t.Errorf("API.APIGetStats() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("API.APIGetStats() = %v, want %v", got, tt.want)
			}
		})
	}
}
