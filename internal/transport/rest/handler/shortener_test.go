package handler

import (
	"bytes"
	"context"
	"errors"
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/Vasily-van-Zaam/ushortener/internal/core"
	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
)

// /// mock
type ServiceMock struct {
}

func (s *ServiceMock) GetURL(ctx context.Context, id string) (string, error) {
	switch id {
	case "1":
		{
			return "http://example.com/link1", nil
		}
	case "2":
		{
			return "http://example.com/link2", nil
		}
	default:
		{
			return "", errors.New("not Found")
		}
	}
}
func (s *ServiceMock) SetURL(ctx context.Context, link string) (string, error) {
	switch link {
	case "http://example.com/link1":
		{
			return core.MAINDOMAIN + "1", nil
		}
	case "http://example.com/link2":
		{
			return core.MAINDOMAIN + "2", nil
		}
	default:
		{
			return "", nil
		}
	}
}

func TestShortenerHandler_GetSetURL(t *testing.T) {

	service := ServiceMock{}
	type fields struct {
		service ShortenerService
	}
	type want struct {
		code        int
		response    string
		contentType string
	}
	type args struct {
		w *httptest.ResponseRecorder
		r *http.Request
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   want
	}{
		// TODO: Add test cases.
		{
			name: "set short link 1: ",
			fields: fields{
				service: &service,
			},
			args: args{
				w: httptest.NewRecorder(),
				r: httptest.NewRequest(
					http.MethodPost,
					"/",
					bytes.NewReader([]byte("http://example.com/link1")),
					// strings.NewReader("http://example.com/link1"),
				),
			},
			want: want{
				code:        201,
				response:    `http://localhost:8080/1`,
				contentType: "text/plain",
			},
		},
		{
			name: "set short link 2: ",
			fields: fields{
				service: &service,
			},
			args: args{
				w: httptest.NewRecorder(),
				r: httptest.NewRequest(
					http.MethodPost,
					"/",
					strings.NewReader("http://example.com/link2"),
				),
			},
			want: want{
				code:        201,
				response:    `http://localhost:8080/2`,
				contentType: "text/plain",
			},
		},
		{
			name: "set short link empty post body: ",
			fields: fields{
				service: &service,
			},
			args: args{
				w: httptest.NewRecorder(),
				r: httptest.NewRequest(
					http.MethodPost,
					"/",
					nil,
				),
			},
			want: want{
				code:        201,
				response:    `http://localhost:8080/`,
				contentType: "text/plain",
			},
		},
		{
			name: "get short link 1: ",
			fields: fields{
				service: &service,
			},
			args: args{
				w: httptest.NewRecorder(),
				r: httptest.NewRequest(http.MethodGet, "/1", nil),
			},
			want: want{
				code:        307,
				response:    `http://example.com/link1`,
				contentType: "text/plain",
			},
		},
		{
			name: "get short link 2: ",
			fields: fields{
				service: &service,
			},
			args: args{
				w: httptest.NewRecorder(),
				r: httptest.NewRequest(http.MethodGet, "/2", nil),
			},
			want: want{
				code:        307,
				response:    `http://example.com/link2`,
				contentType: "text/plain",
			},
		},
		{
			name: "get error: not Found: ",
			fields: fields{
				service: &service,
			},
			args: args{
				w: httptest.NewRecorder(),
				r: httptest.NewRequest(http.MethodGet, "/", nil),
			},
			want: want{
				code:        400,
				response:    "not Found\n",
				contentType: "text/plain",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := &ShortenerHandler{
				service: tt.fields.service,
			}
			r := chi.NewRouter()
			hs := NewHandlers(h)
			hs.InitAPI(r)
			r.ServeHTTP(tt.args.w, tt.args.r)

			///////// chech response //////////
			res := tt.args.w.Result()

			defer res.Body.Close()
			resBody, err := io.ReadAll(res.Body)
			if err != nil {
				t.Fatal(err)
			}
			assert.Equal(t, tt.want.code, res.StatusCode)
			assert.Equal(t, tt.want.response, string(resBody))

			log.Println(tt.name, string(resBody), res.StatusCode)

		})
	}
}

func TestShortenerHandler_ApiSetShorten(t *testing.T) {
	type fields struct {
		service ShortenerService
	}
	type args struct {
		w http.ResponseWriter
		r *http.Request
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		{},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := &ShortenerHandler{
				service: tt.fields.service,
			}
			h.ApiSetShorten(tt.args.w, tt.args.r)
		})
	}
}
