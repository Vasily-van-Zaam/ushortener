package service

import (
	"context"
	"reflect"
	"testing"

	"github.com/Vasily-van-Zaam/ushortener/internal/core"
)

func TestAPI_BindBuferIds(t *testing.T) {
	type fields struct {
		storage     Storage
		config      *core.Config
		AUTHService core.AUTHService
	}
	tests := []struct {
		name   string
		fields fields
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &API{
				storage:     tt.fields.storage,
				config:      tt.fields.config,
				AUTHService: tt.fields.AUTHService,
			}
			s.BindBuferIds()
		})
	}
}

func TestAPI_APISetShorten(t *testing.T) {
	type fields struct {
		storage     Storage
		config      *core.Config
		AUTHService core.AUTHService
	}
	type args struct {
		ctx     context.Context
		request *core.RequestAPIShorten
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *core.ResponseAPIShorten
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &API{
				storage:     tt.fields.storage,
				config:      tt.fields.config,
				AUTHService: tt.fields.AUTHService,
			}
			got, err := s.APISetShorten(tt.args.ctx, tt.args.request)
			if (err != nil) != tt.wantErr {
				t.Errorf("API.APISetShorten() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("API.APISetShorten() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestAPI_APIGetUserURLS(t *testing.T) {
	type fields struct {
		storage     Storage
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
		want    []*core.ResponseAPIUserURL
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &API{
				storage:     tt.fields.storage,
				config:      tt.fields.config,
				AUTHService: tt.fields.AUTHService,
			}
			got, err := s.APIGetUserURLS(tt.args.ctx)
			if (err != nil) != tt.wantErr {
				t.Errorf("API.APIGetUserURLS() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("API.APIGetUserURLS() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestAPI_APISetShortenBatch(t *testing.T) {
	type fields struct {
		storage     Storage
		config      *core.Config
		AUTHService core.AUTHService
	}
	type args struct {
		ctx     context.Context
		request []*core.RequestAPIShortenBatch
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    []*core.ResponseAPIShortenBatch
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &API{
				storage:     tt.fields.storage,
				config:      tt.fields.config,
				AUTHService: tt.fields.AUTHService,
			}
			got, err := s.APISetShortenBatch(tt.args.ctx, tt.args.request)
			if (err != nil) != tt.wantErr {
				t.Errorf("API.APISetShortenBatch() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("API.APISetShortenBatch() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestAPI_APIDeleteUserURLS(t *testing.T) {
	type fields struct {
		storage     Storage
		config      *core.Config
		AUTHService core.AUTHService
	}
	type args struct {
		ctx context.Context
		ids []*string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &API{
				storage:     tt.fields.storage,
				config:      tt.fields.config,
				AUTHService: tt.fields.AUTHService,
			}
			if err := s.APIDeleteUserURLS(tt.args.ctx, tt.args.ids); (err != nil) != tt.wantErr {
				t.Errorf("API.APIDeleteUserURLS() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
