package search

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"go-core/7-lesson/pkg/index"
	"go-core/7-lesson/pkg/storage"
)

// Service это структура для хранения состояния
type Service struct {
	Urls    []string
	Storage storage.ReaderWriter
	index   index.FileData
	channel <-chan int
}

// New создает объект Service
func New(st storage.ReaderWriter) *Service {
	var s Service
	s.Storage = st
	// Инициализируем индекс
	s.initIndex()
	return &s
}

// Channel инициализирует канал для сообщений о готовности данных
func (s *Service) Channel(ch <-chan int) {
	s.channel = ch
}

// Listen слушает канал ch и обновляет currentIndex при поступлении сигнала
func (s *Service) Listen() {
	for {
		select {
		case <-s.channel:
			// Обработка события: выводим сообщение, чтоб было понятнее
			fmt.Println("\n[update] Индекс обновлен")
			s.initIndex()
		}
	}
}

// initIndex обновляет текущий индекс для поиска
func (s *Service) initIndex() {
	var err error
	s.index, err = index.ReadData(s.Storage)
	if err != nil {
		// При ошибке не выходим, продолжаем искать в старой структуре
		fmt.Printf("\n[update] Ошибка чтения данных из файла: %v", err)
	}
}

// Search реализует ввод фразы с клавиатуры и поиск в индексе
func (s *Service) Search() {
	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Print("\n[search] Введите поисковую фразу: ")
		phrase, _ := reader.ReadString('\n')
		phrase = strings.Replace(phrase, "\r\n", "", -1)
		phrase = strings.Replace(phrase, "\n", "", -1)

		if phrase != "" {
			fmt.Printf("\n[search] Поиск по строке \"%s\"", phrase)
			found := false
			for _, document := range s.searchPhrase(phrase) {
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
func (s *Service) searchPhrase(phrase string) []index.Record {
	var res []index.Record
	if s.index.Hash == nil {
		return res
	}
	if ids, found := s.index.Hash[strings.ToLower(phrase)]; found {
		// Фраза найдена в хеше, ids содержит индексы документов (Record.ID) в массиве Records
		for _, id := range ids {
			// Поиск записей в Records по id (Record.ID)
			record := s.index.Records.Search(id)
			if record == nil {
				continue
			}
			res = append(res, record.Value.(index.Record))
		}
	}
	return res
}
