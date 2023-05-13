package gshort

import (
	context "context"
	"log"
	"net"
	"net/http"

	"github.com/Vasily-van-Zaam/ushortener/internal/core"
	grpc "google.golang.org/grpc"
)

// Ranner.
type runner interface {
	Run(string) error
}

// Implements Service.
type Service interface {
	APISetShorten(ctx context.Context, request *core.RequestAPIShorten) (*core.ResponseAPIShorten, error)
	APISetShortenBatch(ctx context.Context,
		request []*core.RequestAPIShortenBatch) ([]*core.ResponseAPIShortenBatch, error)
	APIGetUserURLS(ctx context.Context) ([]*core.ResponseAPIUserURL, error)
	APIDeleteUserURLS(ctx context.Context, ids []*string) error
	APIGetStats(r *http.Request) (*core.Stats, error)

	GetURL(ctx context.Context, link string) (string, error)
	SetURL(ctx context.Context, link string) (string, error)
	Ping(ctx context.Context) error
	core.AUTHService
}

type server struct {
	config *core.Config
}

// DeleteUserURLS implements GRPCServerServer.
func (*server) DeleteUserURLS(context.Context, *DeleteUserURLSRequest) (*DeleteUserURLSResponse, error) {
	panic("unimplemented")
}

// GetBaseURL implements GRPCServerServer.
func (*server) GetBaseURL(context.Context, *ShortUrlRequest) (*UrlResponse, error) {
	panic("unimplemented")
}

// GetStats implements GRPCServerServer.
func (*server) GetStats(context.Context, *GetStatsRequest) (*StatsResponse, error) {
	panic("unimplemented")
}

// GetUserURLS implements GRPCServerServer.
func (*server) GetUserURLS(context.Context, *GetUserURLSRequest) (*UserUrlsResponse, error) {
	panic("unimplemented")
}

// Ping implements GRPCServerServer.
func (*server) Ping(context.Context, *GetPingRequest) (*PingResponse, error) {
	panic("unimplemented")
}

// SetURL implements GRPCServerServer.
func (*server) SetURL(context.Context, *UrlRequest) (*UrlResponse, error) {
	panic("unimplemented")
}

// SetUrls implements GRPCServerServer.
func (*server) SetUrls(context.Context, *ShortenBatchRequest) (*ShortenBatchRequest, error) {
	panic("unimplemented")
}

// mustEmbedUnimplementedGRPCServerServer implements GRPCServerServer.
func (*server) mustEmbedUnimplementedGRPCServerServer() {
	panic("unimplemented")
}

// Run implements runner.
func (srv *server) Run(addresPort string) error {
	listen, err := net.Listen("tcp", addresPort)
	if err != nil {
		log.Fatal(err)
	}
	s := grpc.NewServer()

	RegisterGRPCServerServer(s, srv)
	return s.Serve(listen)
}

// Create server.
func New(conf *core.Config) runner {
	return &server{
		config: conf,
	}
}
