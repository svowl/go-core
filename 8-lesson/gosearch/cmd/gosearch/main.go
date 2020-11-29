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

	var urls = []string{"https://www.google.com", "https://go.dev", "https://golang.org"} //, "https://habr.com"}

	// Создание и инициализация gosearch
	// Для инициализации передается список сайтов и путь к файлу БД
	g := new()
	err := g.init(urls, "./index.json")
	if err != nil {
		log.Fatalf("Ошибка инициализации: %s", err)
	}

	// Запуск генерации поисковых данных в отдельном процессе
	go g.build()

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
	// Сканер
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
	for _, url := range g.urls {
		fmt.Printf("\n[build] Сканируем  %s...", url)
		data, err := g.spider.Scan(url, 2)
		if err != nil {
			// Ошибка при сканировании сайта, игнорируем и идем дальше
			continue
		}
		fmt.Printf("\n[build]  ...найдено %d документов, обрабатываем...", len(data))
		for _, doc := range data {
			// Добавление документа в БД. При этом для документа генерируется ID
			err := g.db.AddDoc(&doc)
			if err != nil {
				// Ошибка при добавлении документа - игнорируем его и идем дальше
				fmt.Printf("\n[warning] %v", err)
				continue
			}
			// Добавление документа в индекс
			g.index.Add(doc)

		}
	}
	cnt := g.db.Count()
	if cnt == 0 {
		log.Fatal("[build] В БД нет документов!")
	}
	// Сохраняем документы БД в файл
	_, err := g.db.Save()
	if err != nil {
		fmt.Printf("[build] %s", err)
	}
	fmt.Printf("\n[build] Сохранено %d документов", g.db.Count())
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
