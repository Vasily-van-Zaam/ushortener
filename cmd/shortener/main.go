package main

import (
	"database/sql"
	"log"

	_ "github.com/mattn/go-sqlite3"

	"github.com/Vasily-van-Zaam/ushortener/docs"
	"github.com/Vasily-van-Zaam/ushortener/internal/service"
	litestore "github.com/Vasily-van-Zaam/ushortener/internal/storage/sqllite"
	"github.com/Vasily-van-Zaam/ushortener/internal/transport/rest"
	"github.com/Vasily-van-Zaam/ushortener/internal/transport/rest/handler"
)

func main() {
	docs.SwaggerInfo.Title = "Shortener API"
	docs.SwaggerInfo.Description = "This is a link shortener server."
	docs.SwaggerInfo.Version = "1.1"

	db, err := sql.Open("sqlite3", "store_shortener.db")
	if err != nil {
		panic(err)
	}
	defer db.Close()
	storage := litestore.New(db)

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
