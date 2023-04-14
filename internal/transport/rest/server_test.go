package rest_test

import (
	"testing"
	"time"

	"github.com/Vasily-van-Zaam/ushortener/internal/core"
	"github.com/Vasily-van-Zaam/ushortener/internal/transport/rest"
	"github.com/go-chi/chi/v5"
)

func Test_server_Run(t *testing.T) {
	type fields struct {
		router *chi.Mux
		config *core.Config
	}
	type args struct {
		addresPort string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "RUN SERVER",
			fields: fields{
				router: chi.NewMux(),
				config: &core.Config{
					ServerAddress: "localhost:8080",
					BaseURL:       "http://localhost:8080",
				},
			},
			args: args{
				addresPort: "localhost:8080",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s, _ := rest.New(tt.fields.config, tt.fields.router)

			go func() {
				if err := s.Run(tt.args.addresPort); (err != nil) != tt.wantErr {
					t.Errorf("server.Run() error = %v, wantErr %v", err, tt.wantErr)
				}
			}()
			time.Sleep(time.Second)
		})
	}
}
