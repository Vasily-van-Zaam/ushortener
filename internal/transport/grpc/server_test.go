package gshort

import (
	context "context"
	"errors"
	"log"
	"net/http"
	"testing"

	"github.com/Vasily-van-Zaam/ushortener/internal/core"
	grpc "google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"
)

type mockBasicservice struct {
}

// CreateUser implements basicService.
func (mockBasicservice) CreateUser() {
	panic("unimplemented")
}

// GetURL implements basicService.
func (mockBasicservice) GetURL(ctx context.Context, link string) (string, error) {
	panic("unimplemented")
}

// Ping implements basicService.
func (mockBasicservice) Ping(ctx context.Context) error {
	// isErr := ctx.Value("err")

	isErr, ok := metadata.FromIncomingContext(ctx)
	if ok && isErr["err"] != nil {
		log.Println(isErr)
		return errors.New("test error: ")
	}
	return nil
}

// SetURL implements basicService.
func (mockBasicservice) SetURL(ctx context.Context, link string) (string, error) {
	panic("unimplemented")
}

type mockService struct {
}

// CreateUser implements apiService.
func (mockService) CreateUser() {
	panic("unimplemented")
}

// APIDeleteUserURLS implements apiService.
func (mockService) APIDeleteUserURLS(ctx context.Context, ids []string) error {
	panic("unimplemented")
}

// APIGetStats implements apiService.
func (mockService) APIGetStats(r *http.Request) (*core.Stats, error) {
	return &core.Stats{
		Urls:  100011,
		Users: 10,
	}, nil
}

// APIGetUserURLS implements apiService.
func (mockService) APIGetUserURLS(ctx context.Context) ([]*core.ResponseAPIUserURL, error) {
	panic("unimplemented")
}

// APISetShorten implements apiService.
func (mockService) APISetShorten(ctx context.Context, request *core.RequestAPIShorten) (*core.ResponseAPIShorten, error) {
	panic("unimplemented")
}

// APISetShortenBatch implements apiService.
func (mockService) APISetShortenBatch(ctx context.Context, request []*core.RequestAPIShortenBatch) ([]*core.ResponseAPIShortenBatch, error) {
	panic("unimplemented")
}

type fields struct {
	service      apiService
	basicService basicService
	config       *core.Config
}

const addresPort = ":3200"

func Test_server_Ping(t *testing.T) {
	type args struct {
		pingErr bool
	}
	f := fields{
		service:      &mockService{},
		basicService: &mockBasicservice{},
	}

	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name:    "ping ok",
			fields:  f,
			args:    args{},
			wantErr: false,
		},
		{
			name:   "ping err",
			fields: f,
			args: args{
				pingErr: true,
			},
			wantErr: true,
		},
	}
	srv := New(f.config, f.service, f.basicService)

	go srv.Run(addresPort)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			conn, err := grpc.Dial(addresPort, grpc.WithTransportCredentials(insecure.NewCredentials()))
			if err != nil {
				log.Fatal(err)
			}
			defer conn.Close()
			c := NewGrpcClient(conn)
			ctx := context.Background()
			if tt.args.pingErr {
				md := metadata.New(map[string]string{"err": "true"})
				ctx = metadata.NewOutgoingContext(context.Background(), md)
			}
			_, got1 := c.Ping(ctx, &GetPingRequest{})

			if tt.wantErr && got1 == nil {
				t.Errorf("server.Run() error = %v, wantErr %v", err, tt.wantErr)
			}

			// time.Sleep(time.Second * 5)
			log.Printf("stop server")
		})
	}
	defer func() {
		err := srv.Stop()
		if err != nil {
			log.Println(err)
		}
	}()
}

func Test_server_GetStats(t *testing.T) {
	tests := []struct {
		name    string
		want    *StatsResponse
		wantErr bool
	}{
		{
			name:    "gt stats ok",
			wantErr: false,
			want: &StatsResponse{
				Urls:  100011,
				Users: 10,
			},
		},
	}
	f := fields{
		service:      &mockService{},
		basicService: &mockBasicservice{},
	}
	srv := New(f.config, f.service, f.basicService)

	go srv.Run(addresPort)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			conn, err := grpc.Dial(addresPort, grpc.WithTransportCredentials(insecure.NewCredentials()))
			if err != nil {
				log.Fatal(err)
			}
			defer conn.Close()
			c := NewGrpcClient(conn)
			got, got1 := c.GetStats(context.Background(), &GetStatsRequest{})
			log.Println(got, got1)
			if (err != nil) != tt.wantErr {
				t.Errorf("server.GetStats error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got.Urls != tt.want.Urls || got.Users != tt.want.Users {
				t.Errorf("server.GetStats = %v, want %v", got, tt.want)
			}
		})
	}
	defer func() {
		err := srv.Stop()
		if err != nil {
			log.Println(err)
		}
	}()
}
