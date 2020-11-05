package main

import (
	"go-core/7-lesson/cmd/builder"
	"go-core/7-lesson/cmd/search"
	"go-core/7-lesson/pkg/spider"
	"go-core/7-lesson/pkg/storage/file"
	"log"
)

func main() {

	// Список URL для сканирования
	var urls = []string{"https://www.google.com", "https://go.dev", "https://golang.org"} //, "https://habr.com"}

	// Создаем паука
	var sp spider.Spider

	// Создаем файловый storage.
	// Изначально файл index.json хранит данные после сканирования https://www.google.com
	st := file.Storage{FileName: "./index.json"}

	// Engine запускает два независимых процесса - [сканирование + построение индекса] и [поиск].
	// Каждая часть работает со своим экземпляром индекса.
	// Поиск обновляет свой индекс из файла после каждой записи туда новых данных.

	// Канал для оповещений поисковой части об изменениях в индексе
	ch := make(chan int, 1)

	// Создаем билдер, инициализиуем его urls, spider и storage
	b, err := builder.New(urls, &sp, &st)
	if err != nil {
		log.Fatal(err)
	}
	// Передаем ему канал
	b.Channel(ch)

	// Создаем поисковик и передаем ему канал
	s := search.New(&st)
	s.Channel(ch)

	// Запускаем сканер в отдельном потоке, передаем паука, индекс и канал
	go b.Build()

	// Запускаем слушателя в отдельном потоке, чтобы максимально развести crawler и поиск
	go s.Listen()

	// Запускаем интерфейсную часть: ввод с клавиатуры и поиск в текущем состоянии индекса
	s.Search()
}
