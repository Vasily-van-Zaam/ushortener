// Сервис кортких ссылок
package main

import (
	"context"
	"log"
	"os/signal"
	"syscall"

	"github.com/caarlos0/env/v6"

	"github.com/Vasily-van-Zaam/ushortener/docs"
	"github.com/Vasily-van-Zaam/ushortener/internal/core"
	"github.com/Vasily-van-Zaam/ushortener/internal/service"

	filestore "github.com/Vasily-van-Zaam/ushortener/internal/storage/file"
	memorystore "github.com/Vasily-van-Zaam/ushortener/internal/storage/memory"
	"github.com/Vasily-van-Zaam/ushortener/internal/storage/psql"
	gshort "github.com/Vasily-van-Zaam/ushortener/internal/transport/grpc"
	"github.com/Vasily-van-Zaam/ushortener/internal/transport/rest"
	"github.com/Vasily-van-Zaam/ushortener/internal/transport/rest/handler"
	"github.com/Vasily-van-Zaam/ushortener/internal/transport/rest/middleware"
)

// These variables are configured using ldflags.
// For example:
// go run -ldflags "-X main.buildVersion=v2.1.1  -X 'main.buildDate=2023-04-26'" ./cmd/shortener/main.go.
var (
	buildVersion = "N/A"
	buildDate    = "N/A"
	buildCommit  = "N/A"
)

func main() {
	docs.SwaggerInfo.Title = "Shortener API"
	docs.SwaggerInfo.Description = "This is a link shortener server."
	docs.SwaggerInfo.Version = "1.1"

	var cfg core.Config
	var storage service.Storage
	err := env.Parse(&cfg)
	if err != nil {
		log.Fatal(err)
	}
	cfg.SetDefault()

	switch {
	case cfg.DataBaseDNS != "":
		storage, err = psql.New(&cfg)
		log.Println("POSTGRES STORE")
		if err != nil {
			log.Panicln(err)
		}
	case cfg.Filestore != "":
		storage, err = filestore.New(&cfg)

		log.Println("FILE STORE")
		if err != nil {
			log.Panicln(err)
		}
		go storage.Update()
	default:
		storage, err = memorystore.New(&cfg)
		log.Println("MEMORY STORE")
		if err != nil {
			log.Panicln(err)
		}
	}

	defer storage.Close()
	authService := service.NewAuth(&cfg, &storage)
	basicService := service.NewBasic(&cfg, &storage, authService)

	apiService := service.NewAPI(&cfg, &storage, authService)
	go apiService.BindBuferIds()

	handlers := handler.New(
		&cfg,
		basicService,
		apiService,
		middleware.NewAuth(&cfg, authService).Handle,
		middleware.NewGzip(&cfg).Handle,
	)

	server, routerError := rest.New(
		&cfg,
		handlers,
	)

	grpcServer := gshort.New(&cfg, apiService, basicService)
	if routerError != nil {
		log.Println("routerError:", routerError)
	}

	log.Printf("\nBuild version: %s\nBuild date: %s\nBuild commit: %s\n", buildVersion, buildDate, buildCommit)
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	go func() {
		errorServer := server.Run(cfg.ServerAddress)
		if errorServer != nil {
			log.Println("errorServer:", errorServer)
		}
	}()

	go func() {
		const grpcAddress = ":3200"
		grpcErr := grpcServer.Run(grpcAddress)
		if grpcErr != nil {
			log.Println("grpc:", grpcErr)
		}
	}()

	<-ctx.Done()
}
