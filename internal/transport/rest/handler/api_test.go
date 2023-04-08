package handler_test

import (
	"context"
	"net/http"
	"testing"

	"github.com/Vasily-van-Zaam/ushortener/internal/core"
)

// /// mock.
type APIServiceMock struct {
}

func (s *APIServiceMock) APISetShorten(
	ctx context.Context, request *core.RequestAPIShorten) (*core.ResponseAPIShorten, error) {
	return &core.ResponseAPIShorten{}, nil
}

func Test_apiHandler_apiSetShorten(t *testing.T) {
	type fields struct {
		Service APIServiceMock
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
			// h := handler.New()
			// h.apiSetShorten(tt.args.w, tt.args.r)
		})
	}
}
