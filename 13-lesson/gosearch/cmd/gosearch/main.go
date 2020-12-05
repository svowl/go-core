package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
	"sync"

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
	}

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

// scanner это метод, вызывающий процесс сканирования url из канала urls и возвращающий результат в канал results
func (g *gosearch) scanner(id int, urls <-chan string, result chan<- []crawler.Document, wg *sync.WaitGroup) {
	defer wg.Done()
	for url := range urls {
		fmt.Printf("\n[scanner #%d]: Начинает сканировать  %s...", id, url)
		data, err := g.spider.Scan(url, 2)
		if err != nil {
			// Ошибка при сканировании сайта, пишем в канал nil и идем дальше
			fmt.Printf("\n[scanner #%d] Ошибка сканирования сайта: %v", id, err)
			result <- nil
			continue
		}
		fmt.Printf("\n[scanner #%d] Закончил сканировать %s, найдено %d документов", id, url, len(data))
		result <- data
	}
}

// build запускает сканирование сайтов, сохранение данных в БД и построение индекса
func (g *gosearch) build() {
	W := 10
	jobs := make(chan string)
	res := make(chan []crawler.Document)

	// Группа ожидания для воркеров
	var wg sync.WaitGroup
	wg.Add(W)
	// Запускаем сканеры
	for i := 1; i <= W; i++ {
		go g.scanner(i, jobs, res, &wg)
	}

	// Группа ожидания для результатов
	var wgRes sync.WaitGroup
	wgRes.Add(len(g.urls))
	// Запускаем обработчик результатов
	go func(ch chan []crawler.Document) {
		for data := range ch {
			if data != nil {
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
			wgRes.Done()
		}
	}(res)
	// Записываем ссылки в канал jobs
	for _, url := range g.urls {
		jobs <- url
	}
	close(jobs)
	wg.Wait()
	wgRes.Wait()

	// Проверяем есть ли вообще данные в БД и если нет - завершаем работу с ошибкой, т.к. это фатальная ситуация
	cnt := g.db.Count()
	if cnt == 0 {
		log.Fatal("[build] В БД нет документов!")
	}
	// Сохраняем документы БД в файл
	_, err := g.db.Save()
	if err != nil {
		fmt.Printf("[build]  Ошибка сохранения данных в файл: %s", err)
	} else {
		fmt.Printf("\n[build] Сохранено %d документов", g.db.Count())
	}
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
