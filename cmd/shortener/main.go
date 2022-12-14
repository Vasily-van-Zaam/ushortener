package main

import (
	"log"

	"github.com/caarlos0/env/v6"

	_ "github.com/mattn/go-sqlite3"

	"github.com/Vasily-van-Zaam/ushortener/docs"
	"github.com/Vasily-van-Zaam/ushortener/internal/core"
	"github.com/Vasily-van-Zaam/ushortener/internal/service"
	memorystore "github.com/Vasily-van-Zaam/ushortener/internal/storage/memory"
	sqlitestore "github.com/Vasily-van-Zaam/ushortener/internal/storage/sqllite"
	"github.com/Vasily-van-Zaam/ushortener/internal/transport/rest"
	"github.com/Vasily-van-Zaam/ushortener/internal/transport/rest/handler"
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

	if cfg.SqliteDB == "" {
		// select store memory
		storage, _ = sqlitestore.New(&cfg)
	} else {
		// selcet
		storage, _ = memorystore.New(&cfg)
	}

	defer storage.Close()

	service := service.NewService(storage)
	handlers := handler.NewHandlers(handler.NewShortenerHandler(service))
	server, routerError := rest.NewServer(handlers)
	if routerError != nil {
		log.Println("routerError:", routerError)
	}
	errorServer := server.Run(":8080")
	if errorServer != nil {
		log.Println("errorServer:", errorServer)
	}

}
