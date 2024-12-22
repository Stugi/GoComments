package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	apipkg "stugi/go-comment/pkg/api"
	cachepkg "stugi/go-comment/pkg/cache"
	"stugi/go-comment/pkg/service"
	"stugi/go-comment/pkg/storage"
)

// конфигурация приложения.
type config struct {
	URLS   []string `json:"rss"`
	Period int      `json:"request_period"`
}

func main() {
	// инициализация зависимостей приложения
	cache := cachepkg.New(time.Hour * 24)

	// инициализация БД
	db, err := storage.New(cache)
	if err != nil {
		log.Fatal(err)
	}
	// инициализация сервиса комментариев
	comment := service.New(db)
	// инициализация API
	api := apipkg.New(comment)

	// Асинхронный внутренний процесс (например, горутина, запущенная из функции main),
	// который будет брать новые комментарии, проверять их на наличие
	// недопустимых слов или словосочетаний и ставить отметку:
	// допускается к показу или заблокирован. Под списком фраз имеются в виду ругательства, но так как мы не можем себе позволить писать это в исходниках,
	// то предлагаю использовать в качестве недопустимых слов qwerty, йцукен, zxvbnm.
	go func() {
		for {
			fmt.Println("Checking comments...")
			err = comment.CheckComments()
			if err != nil {
				log.Println(err)
			}
			time.Sleep(time.Second * 10)
		}
	}()

	// запуск веб-сервера с API и приложением
	err = http.ListenAndServe(":80", api.Router())
	if err != nil {
		log.Fatal(err)
	}
}
