// Test store memory
package memorystore_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/Vasily-van-Zaam/ushortener/internal/core"
	memorystore "github.com/Vasily-van-Zaam/ushortener/internal/storage/memory"
	"github.com/google/uuid"
)

func TestStore_GetURL(t *testing.T) {
	type fields struct {
		Config *core.Config
	}

	type args struct {
		ctx  context.Context
		id   string
		uuid string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    string
		want1   int
		wantErr bool
	}{
		{
			name: "test",
			args: args{
				ctx:  context.Background(),
				id:   "1",
				uuid: "user_id",
			},
			want:  "http://link/1",
			want1: 1,
		},
		{
			name: "test",
			args: args{
				ctx:  context.Background(),
				id:   "2",
				uuid: "user_id",
			},
			want:  "http://link/2",
			want1: 2,
		},
		{
			name: "test",
			args: args{
				ctx:  context.Background(),
				id:   "4",
				uuid: "user_id",
			},
			want:  "",
			want1: 3,
		},
	}
	s, err := memorystore.New(&core.Config{})
	if err != nil {
		t.Error(err)
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err = s.SetURL(tt.args.ctx, &core.Link{
				Link: "http://link/" + tt.args.id,
				UUID: tt.args.uuid,
			})
			if err != nil {
				t.Error(err)
			}

			got, errGet := s.GetURL(tt.args.ctx, tt.args.id)
			if (errGet != nil) != tt.wantErr {
				t.Errorf("Store.GetURL() error = %v, wantErr %v", errGet, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("Store.GetURL() = %v, want %v", got, tt.want)
			}
			got1, errGetUrls := s.GetUserURLS(tt.args.ctx, tt.args.uuid)
			if errGetUrls != nil {
				t.Errorf("Store.GetURL() error = %v, wantErr %v", errGetUrls, tt.wantErr)
				return
			}
			if len(got1) != tt.want1 {
				t.Errorf("Store.GetUserURL() got1 = %v, want %v", len(got1), tt.want1)
			}
		})
	}
}

func BenchmarkSetUrl(b *testing.B) {
	userLink := "https://www.google.com/"
	store, err := memorystore.New(&core.Config{})
	if err != nil {
		b.Error(err)
	}

	b.ResetTimer()
	for i := range make([]int, 1000) {
		uuid := uuid.New().String()
		ctx := context.Background()
		link := fmt.Sprint(userLink, i)
		_, err1 := store.SetURL(ctx, &core.Link{
			UUID: uuid,
			Link: fmt.Sprint(userLink, link),
		})
		if err1 != nil {
			b.Error(err)
			return
		}
	}
}
