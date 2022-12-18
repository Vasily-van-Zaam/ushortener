package main

import (
	"fmt"
	"log"
	"strings"

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
	handlers := handler.NewHandlers(handler.NewShortenerHandler(service))
	server, routerError := rest.NewServer(handlers)
	if routerError != nil {
		log.Println("routerError:", routerError)
	}
	serverAddress := strings.Split(cfg.ServerAddress, ":")
	port := ":8080"
	if len(serverAddress) >= 3 {
		port = fmt.Sprint(":", serverAddress[2])
	}
	errorServer := server.Run(port)
	if errorServer != nil {
		log.Println("errorServer:", errorServer)
	}

}
