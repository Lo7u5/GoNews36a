package main

import (
	"GoNews36a/pkg/api"
	"GoNews36a/pkg/dbase/postgresql"
	"GoNews36a/pkg/filler"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"time"
)

type config struct {
	URLS   []string `json:"rss"`
	Period int      `json:"request_period"`
}

func main() {
	//подключение базы данных postgresql
	constr := os.Getenv("NEWSDB")
	db, err := postgresql.New(constr)
	if err != nil {
		log.Fatal(err)
	}

	api := api.New(db)
	//чтение файла конфигурации
	b, err := os.ReadFile("/Users/kv/GolandProjects/GoNews36a/cmd/gonews/config.json")
	if err != nil {
		log.Fatal(err)
	}
	var config config
	err = json.Unmarshal(b, &config)
	if err != nil {
		log.Fatal(err)
	}
	//каналы для обработки новостей и ошибок, запуск обработки rss
	chPosts := make(chan []postgresql.Post)
	chErrors := make(chan error)
	for _, url := range config.URLS {
		go parseURL(url, chPosts, chErrors, config.Period)
	}
	//горутина для обработки новостей
	go func() {
		for posts := range chPosts {
			err := db.AddPost(posts)
			if err != nil {
				log.Println(err)
			}
		}
	}()
	//горутина для обработки ошибок
	go func() {
		for err := range chErrors {
			log.Println(err)
		}
	}()
	//запуск сервера
	err = http.ListenAndServe(":80", api.Router())
	if err != nil {
		log.Fatal(err)
	}
}

// чтение новостей и запись новостей и ошибок в каналы
func parseURL(url string, posts chan<- []postgresql.Post, errors chan<- error, period int) {
	for {
		news, err := filler.Parse(url)
		if err != nil {
			errors <- err
			continue
		}
		posts <- news
		time.Sleep(time.Minute * time.Duration(period))
	}
}
