package handler_test

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
	"github.com/Vasily-van-Zaam/ushortener/internal/transport/rest/handler"
	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
)

// /// mock.
type ServiceMock struct {
	core.AUTHService
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
			return "http://localhost:8080/" + "1", nil
		}
	case "http://example.com/link2":
		{
			return "http://localhost:8080/" + "2", nil
		}
	default:
		{
			return "", nil
		}
	}
}
func (s *ServiceMock) Ping(ctx context.Context) error {
	return nil
}

func (s *ServiceMock) APISetShorten(
	ctx context.Context, request *core.RequestAPIShorten) (*core.ResponseAPIShorten, error) {
	return &core.ResponseAPIShorten{}, nil
}

func TestShortenerHandler_GetSetURL(t *testing.T) {
	service := ServiceMock{}
	type fields struct {
		service handler.BasicService
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
		cfg := core.Config{
			ServerAddress: "127.0.0.1:8080/",
			BaseURL:       "http://localhost:8080/",
		}

		t.Run(tt.name, func(t *testing.T) {
			h := &handler.BasicHandler{
				Service: &tt.fields.service,
				Config:  &cfg,
			}
			r := chi.NewRouter()
			hs := handler.NewHandlers(h, nil)
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
