package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"

	"go-core/7-lesson/pkg/crawler"
	"go-core/7-lesson/pkg/crawler/spider"
	"go-core/7-lesson/pkg/engine"
	"go-core/7-lesson/pkg/index"
	"go-core/7-lesson/pkg/storage"
	"go-core/7-lesson/pkg/storage/file"
)

// gosearch это структура для хранения состояния поискового движка
type gosearch struct {
	urls    []string
	spider  crawler.Scanner
	storage storage.Interface
	index   *index.Index
	engine  *engine.Service
}

func main() {

	// Список URL для сканирования
	var urls = []string{"https://www.google.com", "https://go.dev", "https://golang.org"} //, "https://habr.com"}

	// Создаем объект gosearch
	g := new()

	// Инициализируем его
	err := g.init(urls, "./index.json")
	if err != nil {
		log.Fatalf("Ошибка инициализации индекса: %s", err)
	}

	// Запускаем генерацию поисковых данных
	go g.build()

	// Запускаем интерфейсную часть: ввод с клавиатуры и поиск в текущем состоянии индекса
	g.search()
}

// new создает объект gosearch
func new() *gosearch {
	var g gosearch
	return &g
}

// init инициализирует объект gosearch
func (g *gosearch) init(urls []string, filename string) error {
	g.urls = urls
	g.spider = spider.New()
	g.storage = file.New(filename)
	g.engine = engine.New()

	g.index = index.New()
	err := g.index.Init(g.storage)
	if err != nil {
		return err
	}

	g.engine.Init(g.index)

	return nil
}

// build запускает сканирование страниц, построение индекса и сохранение данных в хранилище
func (g *gosearch) build() {
	for _, url := range g.urls {
		fmt.Printf("\n[build] Сканируем  %s...", url)
		data, err := g.spider.Scan(url, 2)
		if err != nil {
			// ошибка при сканировании сайта, игнорируем и идем дальше
			continue
		}
		fmt.Printf("\n[build]  ...найдено %d документов, индексируем...", len(data))
		// Строим индекс по списку просканированных документов
		_, err = g.index.Build(data)
		if err != nil {
			fmt.Printf("[build] %s", err)
			continue
		}
		// Сохраняем индекс в файл
		err = g.index.SaveData()
		if err != nil {
			fmt.Printf("[build] %s", err)
			continue
		}
		fmt.Printf("\n[build] Проиндексировано %d страниц", g.index.Records.Count)
	}
}

// search реализует ввод фразы с клавиатуры и поиск в индексе
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
