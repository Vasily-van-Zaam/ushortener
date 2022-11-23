package main

import (
	"database/sql"
	"log"

	_ "github.com/mattn/go-sqlite3"

	"github.com/Vasily-van-Zaam/ushortener/internal/service"
	sqllite_storage "github.com/Vasily-van-Zaam/ushortener/internal/storage/sqllite"
	"github.com/Vasily-van-Zaam/ushortener/internal/transsport/rest"
	"github.com/Vasily-van-Zaam/ushortener/internal/transsport/rest/handler"
)

func main() {
	db, err := sql.Open("sqlite3", "store_shortener.db")
	if err != nil {
		panic(err)
	}
	defer db.Close()
	storage := sqllite_storage.New(db)

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
