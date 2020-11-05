package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"

	"go-core/7-lesson/pkg/index"
	"go-core/7-lesson/pkg/spider"
	"go-core/7-lesson/pkg/storage"
	"go-core/7-lesson/pkg/storage/file"
)

// Engine это структура для хранения состояния поискового движка
type Engine struct {
	Urls         []string
	Spider       spider.Scanner
	Storage      storage.ReaderWriter
	index        *index.Index
	currentIndex index.FileData
}

func main() {

	// Список URL для сканирования
	var urls = []string{"https://www.google.com", "https://go.dev", "https://golang.org"} //, "https://habr.com"}

	// Создаем паука
	var sp spider.Spider

	// Создаем файловый storage.
	// Изначально файл index.json хранит данные после сканирования https://www.google.com
	st := file.Storage{FileName: "./index.json"}

	// Создаем поисковый движок, инициализиуем его urls, spider и storage
	e := new(urls, &sp, &st)

	// Engine запускает два независимых процесса - [сканирование + построение индекса] и [поиск].
	// Каждая часть работает со своим экземпляром индекса.
	// Поиск обновляет свой индекс из файла после каждой записи туда новых данных.

	// Канал для оповещений поисковой части об изменениях в индексе
	ch := make(chan int, 1)

	// Запускаем сканер в отдельном потоке, передаем паука, индекс и канал
	go e.build(ch)

	// Запускаем слушателя в отдельном потоке, чтобы максимально развести crawler и поиск
	go e.listen(ch)

	// Запускаем интерфейсную часть: ввод с клавиатуры и поиск в текущем состоянии индекса
	e.search()
}

// Run запускает поисковый движок
func new(urls []string, sp spider.Scanner, st storage.ReaderWriter) *Engine {

	var e Engine
	e.Urls = urls
	e.Spider = sp
	e.Storage = st

	// Инициализируем текущий индекс (для поиска)
	e.updateCurrentIndex()

	var err error
	// Создаем индекс для сканирования
	e.index, err = index.New(e.Storage)
	if err != nil {
		log.Fatalf("[main] Ошибка инициализации индекса: %s", err)
	}

	return &e
}

// build сканирует страницы, строит индекс
func (e *Engine) build(ch chan<- int) {
	for _, url := range e.Urls {
		fmt.Printf("\n[build] Сканируем  %s...", url)
		data, err := e.Spider.Scan(url, 2)
		if err != nil {
			fmt.Printf("\n[build] ошибка при сканировании сайта %s: %v", url, err)
			continue
		}
		fmt.Printf("\n[build]  ...найдено %d документов, индексируем...", len(data))
		// Строим индекс по списку просканированных документов
		_, err = e.index.Build(data)
		if err != nil {
			fmt.Printf("[build] %s", err)
			continue
		}
		// Сохраняем индекс в файл
		err = e.index.SaveData()
		if err != nil {
			fmt.Printf("[build] %s", err)
			continue
		}
		// Посылаем сигнал в канал о необходимости обновить текущий поисковый индекс
		ch <- 1
		fmt.Printf("\n[build] Проиндексировано %d страниц", e.index.Records.Count)
	}
}

// listen слушает канал ch и обновляет currentIndex при поступлении сигнала
func (e *Engine) listen(ch <-chan int) {
	for {
		select {
		case <-ch:
			// Обработка события: выводим сообщение, чтоб было понятнее
			fmt.Println("\n[update] Индекс обновлен")
			e.updateCurrentIndex()
		}
	}
}

// updateCurrentIndex обновляет текущий индекс для поиска
func (e *Engine) updateCurrentIndex() {
	i, err := index.ReadData(e.Storage)
	if err != nil {
		// При ошибке не выходим, продолжаем искать в старой структуре
		fmt.Printf("\n[update] Ошибка чтения данных из файла: %v", err)
	}
	e.currentIndex = i
}

// Реализуем ввод фразы с клавиатуры и поиск в индексе
func (e *Engine) search() {
	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Print("\n[search] Введите поисковую фразу: ")
		phrase, _ := reader.ReadString('\n')
		phrase = strings.Replace(phrase, "\r\n", "", -1)
		phrase = strings.Replace(phrase, "\n", "", -1)

		if phrase != "" {
			fmt.Printf("\n[search] Поиск по строке \"%s\"", phrase)
			found := false
			for _, document := range e.searchPhrase(phrase) {
				fmt.Printf("\n[search]  %s: %s", document.URL, document.Title)
				found = true
			}
			if !found {
				fmt.Println("\n[search] Ничего не найдено")
			}
		}
	}
}

// Search ищет проиндексированные записи по фразе,
func (e *Engine) searchPhrase(phrase string) []index.Record {
	var res []index.Record
	if e.currentIndex.Hash == nil {
		return res
	}
	if ids, found := e.currentIndex.Hash[strings.ToLower(phrase)]; found {
		// Фраза найдена в хеше, ids содержит индексы документов (Record.ID) в массиве Records
		for _, id := range ids {
			// Поиск записей в Records по id (Record.ID)
			record := e.currentIndex.Records.Search(id)
			if record == nil {
				continue
			}
			res = append(res, record.Value.(index.Record))
		}
	}
	return res
}
