package main

import (
	"log"

	"github.com/Vasily-van-Zaam/ushortener/internal/service"
	some_storage "github.com/Vasily-van-Zaam/ushortener/internal/storage/some"
	"github.com/Vasily-van-Zaam/ushortener/internal/transsport/rest"
	"github.com/Vasily-van-Zaam/ushortener/internal/transsport/rest/handler"
)

func test() chan string {
	return make(chan string)
}

func test2(t chan string) chan string {
	return t
}

func main() {
	service := service.NewService(some_storage.New())
	handlers := handler.NewHandlers(handler.NewShortenerHandler(service))
	server, routerError := rest.NewServer(handlers)

	log.Println("routerError:", routerError)
	errorServer := server.Run(":8080")
	log.Println("errorServer:", errorServer)
}
