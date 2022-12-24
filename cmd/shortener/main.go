package main

import (
	"log"

	"github.com/caarlos0/env/v6"

	_ "github.com/mattn/go-sqlite3"

	"github.com/Vasily-van-Zaam/ushortener/docs"
	"github.com/Vasily-van-Zaam/ushortener/internal/core"
	"github.com/Vasily-van-Zaam/ushortener/internal/service"
	filestore "github.com/Vasily-van-Zaam/ushortener/internal/storage/file"
	memorystore "github.com/Vasily-van-Zaam/ushortener/internal/storage/memory"
	sqlite "github.com/Vasily-van-Zaam/ushortener/internal/storage/sqlite"
	"github.com/Vasily-van-Zaam/ushortener/internal/transport/rest"
	"github.com/Vasily-van-Zaam/ushortener/internal/transport/rest/handler"
	"github.com/Vasily-van-Zaam/ushortener/internal/transport/rest/middleware"
)

func main() {
	docs.SwaggerInfo.Title = "Shortener API"
	docs.SwaggerInfo.Description = "This is a link shortener server."
	docs.SwaggerInfo.Version = "1.1"

	var cfg core.Config
	var storage service.ShortenerStorage
	err := env.Parse(&cfg)
	if err != nil {
		log.Fatal(err)
	}
	cfg.SetDefault()

	middlewares := []rest.Middleware{
		middleware.NewGzip(&cfg),
	}
	switch {
	case cfg.SqliteDB != "":
		storage, err = sqlite.New(&cfg)
		if err != nil {
			log.Panicln(err)
		}
	case cfg.Filestore != "":
		storage, err = filestore.New(&cfg)
		if err != nil {
			log.Panicln(err)
		}
	default:
		storage, err = memorystore.New(&cfg)
		if err != nil {
			log.Panicln(err)
		}
	}

	defer storage.Close()

	service := service.NewService(storage)
	handlers := handler.NewHandlers(handler.NewShortenerHandler(service, &cfg))
	server, routerError := rest.NewServer(handlers, &cfg, middlewares)
	if routerError != nil {
		log.Println("routerError:", routerError)
	}

	errorServer := server.Run(cfg.ServerAddress)
	if errorServer != nil {
		log.Println("errorServer:", errorServer)
	}
}
