package main

import (
	"log"

	"github.com/caarlos0/env/v6"

	"net/http"
	_ "net/http/pprof"

	"github.com/Vasily-van-Zaam/ushortener/docs"
	"github.com/Vasily-van-Zaam/ushortener/internal/core"
	"github.com/Vasily-van-Zaam/ushortener/internal/service"
	filestore "github.com/Vasily-van-Zaam/ushortener/internal/storage/file"
	memorystore "github.com/Vasily-van-Zaam/ushortener/internal/storage/memory"
	"github.com/Vasily-van-Zaam/ushortener/internal/storage/psql"
	"github.com/Vasily-van-Zaam/ushortener/internal/transport/rest"
	"github.com/Vasily-van-Zaam/ushortener/internal/transport/rest/handler"
	"github.com/Vasily-van-Zaam/ushortener/internal/transport/rest/middleware"
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
	basicService := service.NewBasic(&cfg, storage, authService)

	apiService := service.NewAPI(&cfg, &storage, authService)
	go apiService.BindBuferIds()
	middlewares := []rest.Middleware{
		middleware.NewGzip(&cfg),
		middleware.NewAuth(&cfg, authService),
	}
	handlers := handler.NewHandlers(
		handler.NewBasic(basicService, &cfg),
		handler.NewAPI(apiService, &cfg),
	)
	server, routerError := rest.NewServer(handlers, &cfg, middlewares)
	if routerError != nil {
		log.Println("routerError:", routerError)
	}

	go server.Run(cfg.ServerAddress)

	log.Fatal(http.ListenAndServe(":8888", nil))
}
