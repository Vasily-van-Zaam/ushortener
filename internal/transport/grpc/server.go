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
	Ping(ctx context.Context) error
	core.AUTHService
}

type server struct {
	UnimplementedGRPCServer
	service Service
	config  *core.Config
}

// DeleteUserURLS implements GRPCServerServer.
func (srv *server) DeleteUserURLS(ctx context.Context, req *DeleteUserURLSRequest) (*DeleteUserURLSResponse, error) {
	list := make([]*string, len(req.Urls))
	for i, u := range req.Urls {
		url := u
		list[i] = &url
	}
	err := srv.service.APIDeleteUserURLS(ctx, list)
	if err != nil {
		return &DeleteUserURLSResponse{
			Error: err.Error(),
		}, err
	}
	return &DeleteUserURLSResponse{
		Error: "",
	}, nil
}

// GetBaseURL implements GRPCServerServer.
func (srv *server) GetBaseURL(ctx context.Context, req *ShortUrlRequest) (*UrlResponse, error) {
	url, err := srv.service.GetURL(ctx, req.ShortUrl)
	if err != nil {
		return nil, err
	}
	return &UrlResponse{
		Result: url,
	}, nil
}

// GetStats implements GRPCServerServer.
func (srv *server) GetStats(ctx context.Context, req *GetStatsRequest) (*StatsResponse, error) {
	haderXReal := http.Header{}
	haderXReal.Add("X-Real-IP", req.Ip)
	r := &http.Request{
		Header: haderXReal,
	}

	stats, err := srv.service.APIGetStats(r)
	if err != nil {
		return nil, err
	}
	return &StatsResponse{
		Urls:  int32(stats.Urls),
		Users: int32(stats.Users),
	}, nil
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

// Run implements runner.
func (srv *server) Run(addresPort string) error {
	listen, err := net.Listen("tcp", addresPort)
	if err != nil {
		log.Fatal(err)
	}
	s := grpc.NewServer()

	RegisterGRPCServer(s, srv)
	return s.Serve(listen)
}

// Create server.
func New(conf *core.Config, s Service) runner {
	return &server{
		config:  conf,
		service: s,
	}
}
