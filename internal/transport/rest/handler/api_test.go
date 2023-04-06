package handler_test

import (
	"context"
	"net/http"
	"testing"

	"github.com/Vasily-van-Zaam/ushortener/internal/core"
	"github.com/Vasily-van-Zaam/ushortener/internal/transport/rest/handler"
)

// /// mock.
type APIServiceMock struct {
}

func (s *APIServiceMock) APISetShorten(
	ctx context.Context, request *core.RequestAPIShorten) (*core.ResponseAPIShorten, error) {
	return &core.ResponseAPIShorten{}, nil
}

func TestApiHandler_APISetShorten(t *testing.T) {
	type fields struct {
		service handler.APIService
		config  *core.Config
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
			// h := &handler.APIHandler{
			// 	Service: tt.fields.service,
			// 	Config:  tt.fields.config,
			// }
			// h.APISetShorten(tt.args.w, tt.args.r)
		})
	}
}

func TestAPIHandler_APIGetUserURLS(t *testing.T) {
	type fields struct {
		Service handler.APIService
		Config  *core.Config
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
			// h := &handler.APIHandler{
			// 	Service: tt.fields.Service,
			// 	Config:  tt.fields.Config,
			// }
			// h.APIGetUserURLS(tt.args.w, tt.args.r)
		})
	}
}

func TestAPIHandler_APIDeleteUserURLS(t *testing.T) {
	type fields struct {
		Service handler.APIService
		Config  *core.Config
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
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := &handler.APIHandler{
				Service: tt.fields.Service,
				Config:  tt.fields.Config,
			}
			h.APIDeleteUserURLS(tt.args.w, tt.args.r)
		})
	}
}

func TestAPIHandler_APISetShortenBatch(t *testing.T) {
	type fields struct {
		Service handler.APIService
		Config  *core.Config
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
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := &handler.APIHandler{
				Service: tt.fields.Service,
				Config:  tt.fields.Config,
			}
			h.APISetShortenBatch(tt.args.w, tt.args.r)
		})
	}
}
