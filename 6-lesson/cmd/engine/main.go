package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"

	"go-core/6-lesson/pkg/index"
	"go-core/6-lesson/pkg/spider"
	"go-core/6-lesson/pkg/storage"
	"go-core/6-lesson/pkg/storage/file"
)

// Scanner это интерфейс, используется для тестирования
type Scanner interface {
	Scan(string, int) (map[string]string, error)
}

// Список URL для сканирования
var urls = []string{"https://www.google.com", "https://go.dev"} //, "https://habr.com"}

// currentIndex содержит текущий индекс и документы, по которым происходит поиск
var currentIndex index.FileData

// build сканирует страницы, строит индекс
func build(s Scanner, i *index.Index, ch chan<- int) {
	for _, url := range urls {
		fmt.Printf("\n[build] Сканируем  %s...", url)
		data, err := s.Scan(url, 2)
		if err != nil {
			fmt.Printf("\n[build] ошибка при сканировании сайта %s: %v", url, err)
			continue
		}
		fmt.Printf("\n[build]  ...найдено %d документов, индексируем...", len(data))
		// Строим индекс по списку просканированных документов
		_, err = i.Build(data)
		if err != nil {
			fmt.Printf("[build] %s", err)
			continue
		}
		// Сохраняем индекс в файл
		err = i.SaveData()
		if err != nil {
			fmt.Printf("[build] %s", err)
			continue
		}
		// Посылаем сигнал в канал о необходимости обновить текущий поисковый индекс
		ch <- 1
		fmt.Printf("\n[build] Проиндексировано %d страниц", i.Records.Count)
	}
}

// listen слушает канал ch и обновляет currentIndex при поступлении сигнала
func listen(r storage.ReaderWriter, ch <-chan int) {
	for {
		select {
		case <-ch:
			// Обработка события: выводим сообщение, чтоб было понятнее
			fmt.Println("\n[update] Индекс обновлен")
			updateCurrentIndex(r)
		}
	}
}

// updateCurrentIndex обновляет текущий индекс для поиска
func updateCurrentIndex(r storage.ReaderWriter) {
	i, err := index.ReadData(r)
	if err != nil {
		// При ошибке не выходим, продолжаем искать в старой структуре
		fmt.Printf("\n[update] Ошибка чтения данных из файла: %v", err)
	}
	currentIndex = i
}

// Реализуем ввод фразы с клавиатуры и поиск в индексе
func search() {
	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Print("\n[search] Введите поисковую фразу: ")
		phrase, _ := reader.ReadString('\n')
		phrase = strings.Replace(phrase, "\r\n", "", -1)
		phrase = strings.Replace(phrase, "\n", "", -1)

		if phrase != "" {
			fmt.Printf("\n[search] Поиск по строке \"%s\"", phrase)
			found := false
			for _, document := range searchPhrase(phrase) {
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
func searchPhrase(phrase string) []index.Record {
	var res []index.Record
	if currentIndex.Hash == nil {
		return res
	}
	if ids, found := currentIndex.Hash[strings.ToLower(phrase)]; found {
		// Фраза найдена в хеше, ids содержит индексы документов (Record.ID) в массиве Records
		for _, id := range ids {
			// Поиск записей в Records по id (Record.ID)
			record := currentIndex.Records.Search(id)
			if record == nil {
				continue
			}
			res = append(res, record.Value.(index.Record))
		}
	}
	return res
}

func main() {
	// Создаем паука
	s := new(spider.Spider)

	// Изначально файл index.json хранит данные после сканирования https://www.google.com
	fname := "./index.json"
	storage := file.Storage{FileName: fname}

	// Инициализируем текущий индекс (для поиска)
	updateCurrentIndex(&storage)

	// Создаем индекс для сканирования
	i, err := index.New(&storage)
	if err != nil {
		log.Fatalf("[main] Ошибка инициализации индекса: %s", err)
	}

	// Канал для оповещений об изменениях в индексе
	ch := make(chan int, 1)

	// Запускаем сканер в отдельном потоке, передаем паука, индекс и канал
	go build(s, i, ch)

	// Запускаем слушателя в отдельном потоке, чтобы максимально развести crawler и поиск
	go listen(&storage, ch)

	// Запускаем поиск
	search()
}
