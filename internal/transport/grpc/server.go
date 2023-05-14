package gshort

import (
	context "context"
	"log"
	"net"
	"net/http"

	"github.com/Vasily-van-Zaam/ushortener/internal/core"
	grpc "google.golang.org/grpc"
	"google.golang.org/grpc/peer"
)

// Ranner.
type runner interface {
	Run(string) error
	Stop() error
}

// Implements Service.
type apiService interface {
	APISetShorten(ctx context.Context, request *core.RequestAPIShorten) (*core.ResponseAPIShorten, error)
	APISetShortenBatch(ctx context.Context,
		request []*core.RequestAPIShortenBatch) ([]*core.ResponseAPIShortenBatch, error)
	APIGetUserURLS(ctx context.Context) ([]*core.ResponseAPIUserURL, error)
	APIDeleteUserURLS(ctx context.Context, ids []string) error
	APIGetStats(r *http.Request) (*core.Stats, error)

	core.AUTHService
}

// Implements basic service.
type basicService interface {
	GetURL(ctx context.Context, link string) (string, error)
	SetURL(ctx context.Context, link string) (string, error)
	Ping(ctx context.Context) error
	core.AUTHService
}

type server struct {
	UnimplementedGrpcServer
	service      apiService
	basicService basicService
	config       *core.Config
	listener     net.Listener
}

// DeleteUserURLS implements GRPCServerServer.
func (srv *server) DeleteUserURLS(ctx context.Context, req *DeleteUserURLSRequest) (*DeleteUserURLSResponse, error) {
	err := srv.service.APIDeleteUserURLS(ctx, req.Urls)
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
	url, err := srv.basicService.GetURL(ctx, req.ShortUrl)
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
	p, ok := peer.FromContext(ctx)
	if ok {
		haderXReal.Add("X-Real-IP", p.Addr.Network())
	}
	log.Println(p)

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
func (srv *server) GetUserURLS(ctx context.Context, req *GetUserURLSRequest) (*UserUrlsResponse, error) {
	urls, err := srv.service.APIGetUserURLS(ctx)
	if err != nil {
		return nil, err
	}

	respUrls := make([]*UserUrl, len(urls))

	for i, u := range urls {
		respUrls[i] = &UserUrl{
			OriginalUrl: u.OriginalURL,
			ShortUrl:    u.ShortURL,
		}
	}

	return &UserUrlsResponse{
		Urls: respUrls,
	}, nil
}

// Ping implements GRPCServerServer.
func (srv *server) Ping(ctx context.Context, req *GetPingRequest) (*PingResponse, error) {
	err := srv.basicService.Ping(ctx)
	if err != nil {
		return &PingResponse{
			Error: err.Error(),
		}, err
	}
	return &PingResponse{
		Error: "",
	}, nil
}

// SetURL implements GRPCServerServer.
func (srv *server) SetURL(ctx context.Context, req *UrlRequest) (*UrlResponse, error) {
	setURL, err := srv.service.APISetShorten(ctx, &core.RequestAPIShorten{
		URL: req.Url,
	})
	if err != nil {
		return nil, err
	}
	return &UrlResponse{
		Result: setURL.Result,
	}, nil
}

// SetUrls implements GRPCServerServer.
func (srv *server) SetUrls(ctx context.Context, req *ShortenBatchRequest) (*ShortenBatchRequest, error) {
	request := make([]*core.RequestAPIShortenBatch, len(req.Urls))
	for i, url := range req.Urls {
		request[i] = &core.RequestAPIShortenBatch{
			OriginalURL:   url.OriginalUrl,
			CorrelationID: url.CorrelationId,
		}
	}

	setURLS, err := srv.service.APISetShortenBatch(ctx, request)
	if err != nil {
		return nil, err
	}
	urls := make([]*ShortenBatchReq, len(setURLS))
	for i, url := range urls {
		urls[i] = &ShortenBatchReq{
			CorrelationId: url.CorrelationId,
			OriginalUrl:   url.OriginalUrl,
		}
	}
	return &ShortenBatchRequest{
		Urls: urls,
	}, nil
}

// Run implements runner.
func (srv *server) Run(addresPort string) error {
	listen, err := net.Listen("tcp", addresPort)
	if err != nil {
		log.Fatal(err)
	}
	srv.listener = listen
	s := grpc.NewServer()
	log.Println("Starting grpc server", addresPort)
	RegisterGrpcServer(s, srv)
	return s.Serve(listen)
}

// Stop implements runner.
func (srv *server) Stop() error {
	return srv.listener.Close()
}

// Create server.
func New(conf *core.Config, s apiService, b basicService) runner {
	return &server{
		config:       conf,
		service:      s,
		basicService: b,
	}
}
