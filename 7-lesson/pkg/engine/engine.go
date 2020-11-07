package engine

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"go-core/7-lesson/pkg/index"
	"go-core/7-lesson/pkg/storage"
)

// Пакет engine принимает пользовательский ввод, запускает поиск в индексе и выводит результаты в консоль

// Service это структура для хранения состояния поискового движка
type Service struct {
	Urls    []string
	Storage storage.ReaderWriter
	index   index.FileData
	channel <-chan int
}

// New создает объект Service
func New() *Service {
	s := Service{}
	return &s
}

// Init инициализирует объект Service
func (s *Service) Init(st storage.ReaderWriter) error {
	s.Storage = st
	// Инициализируем индекс
	err := s.initIndex()
	if err != nil {
		return err
	}
	return nil
}

// Channel инициализирует канал для сообщений о необходимости обновить данные индекса
func (s *Service) Channel(ch <-chan int) {
	s.channel = ch
}

// Listen слушает канал и обновляет currentIndex при поступлении сигнала
func (s *Service) Listen() {
	for {
		select {
		case <-s.channel:
			// Обработка события: выводим сообщение, чтоб было понятнее
			err := s.initIndex()
			if err == nil {
				fmt.Println("\n[update] Индекс обновлен")
			}
		}
	}
}

// initIndex обновляет текущий индекс для поиска
func (s *Service) initIndex() error {
	ind, err := index.ReadData(s.Storage)
	if err != nil {
		return err
	}
	s.index = ind
	return nil
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
