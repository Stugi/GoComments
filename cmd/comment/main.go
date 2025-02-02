package main

import (
	"log"
	"net/http"
	"time"

	cache "stugi/go-comment/pkg/cache"
	"stugi/go-comment/pkg/service"
	"stugi/go-comment/pkg/storage"
	api "stugi/go-comment/pkg/transport"
)

// конфигурация приложения.
type config struct {
	URLS   []string `json:"rss"`
	Period int      `json:"request_period"`
}

func main() {
	// инициализация зависимостей приложения
	cache := cache.New(time.Hour * 24)

	// инициализация БД
	db, err := storage.New(cache)
	if err != nil {
		log.Fatal(err)
	}
	// инициализация сервиса комментариев
	comment := service.New(db)
	// инициализация API
	api := api.New(comment)

	// запуск веб-сервера с API и приложением
	err = http.ListenAndServe(":80", api.Router())
	if err != nil {
		log.Fatal(err)
	}
}
