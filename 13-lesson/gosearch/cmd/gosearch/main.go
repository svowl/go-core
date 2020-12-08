package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"

	"gosearch/pkg/crawler"
	"gosearch/pkg/crawler/spider"
	"gosearch/pkg/engine"
	"gosearch/pkg/index"
	"gosearch/pkg/storage"
	"gosearch/pkg/storage/btree"
	"gosearch/pkg/storage/file"
)

// gosearch содержит сканер, поисковый движок, БД, индекс, а также необходимые для их работы данные:
// хранилище (storage - file), поисковая структура (searcher - btree), список сайтов для сканирования
type gosearch struct {
	urls     []string
	spider   crawler.Scanner
	storage  storage.IReadWriter
	searcher storage.ISearcher
	db       *storage.Db
	index    *index.Index
	engine   *engine.Service
}

func main() {

	// 20 сайтов для сканирования
	var urls = []string{
		"https://www.google.com",
		"https://go.dev",
		"https://golang.org",
		/*
			"https://www.mskagency.ru/",
			"https://www.mos.ru",
			"https://habr.com",
			"https://www.alean.ru",
			"https://www.moscowbooks.ru/",
			"https://www.museum.ru",
			"https://investmoscow.ru",
			"https://mosmetro.ru/",
			"https://govoritmoskva.ru/",
			"https://www.tourister.ru/",
			"https://www.citymoscow.ru/",
			"https://technomoscow.ru/",
			"https://www.tourprom.ru/",
			"https://mosgorzdrav.ru/",
			"https://101hotels.com/",
			"https://moscowchanges.ru/",
			"https://ginza.ru/",
		*/
	}

	// Создание и инициализация gosearch
	// Для инициализации передается список сайтов и путь к файлу БД
	g := new()
	err := g.init(urls, "./index.json")
	if err != nil {
		log.Fatalf("Ошибка инициализации: %s", err)
	}

	// Запуск генерации поисковых данных в отдельном процессе
	g.build()

	// Запуск интерфейсной части: ввод с клавиатуры и выдача результатов поиска в консоли
	g.search()
}

// new создает объект gosearch
func new() *gosearch {
	var g gosearch
	return &g
}

// init инициализирует объект gosearch
func (g *gosearch) init(urls []string, filename string) error {
	// Список сайтов
	g.urls = urls
	// Создаем объект сканер с максимум 10 потоками
	g.spider = spider.New()
	// Хранилище данных для БД
	g.storage = file.New(filename)
	// Поисковая структура для БД
	g.searcher = btree.New()
	// БД
	g.db = storage.New()
	g.db.Init(g.storage, g.searcher)
	// Индекс
	g.index = index.New()
	// Поисковый движок
	g.engine = engine.New()
	g.engine.Init(g.index, g.db)

	return nil
}

// build запускает сканирование сайтов, сохранение данных в БД и построение индекса
func (g *gosearch) build() {
	// Запускаем многопоточное сканирование
	chRes, chErr := g.spider.BatchScan(g.urls, 2, 10)
	// Принимаем отсканированные документы по одному и добавляем в БД и индекс
	go func() {
		for doc := range chRes {
			err := g.db.AddDoc(&doc)
			if err != nil {
				// Ошибка при добавлении документа - игнорируем и идем дальше
				log.Println("Ошибка добавления документа в БД:", err)
				continue
			}
			g.index.Add(doc)
		}
	}()
	go func() {
		for range chErr {
			// Игнориуем ошибки сканирования
		}
	}()
}

// search реализует ввод фразы с клавиатуры, поиск и выдачу результатов в консоль
func (g *gosearch) search() {
	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Print("\n[search] Введите поисковую фразу: ")
		phrase, _ := reader.ReadString('\n')
		phrase = strings.Replace(phrase, "\r\n", "", -1)
		phrase = strings.Replace(phrase, "\n", "", -1)

		if phrase != "" {
			fmt.Printf("\n[search] Поиск по строке \"%s\"", phrase)
			docs, found := g.engine.Search(phrase)
			if found == false {
				fmt.Println("\n[search] Ничего не найдено")
				continue
			}
			for _, document := range docs {
				fmt.Printf("\n[search]  %s: %s", document.URL, document.Title)
			}
		}
	}
}
