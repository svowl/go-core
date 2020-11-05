package builder

import (
	"fmt"

	"go-core/7-lesson/pkg/index"
	"go-core/7-lesson/pkg/spider"
	"go-core/7-lesson/pkg/storage"
)

// Service это структура для хранения состояния билдера
type Service struct {
	Urls    []string
	Spider  spider.Scanner
	Storage storage.ReaderWriter
	channel chan<- int
	index   *index.Index
}

// New создает объект Service
func New(urls []string, sp spider.Scanner, st storage.ReaderWriter) (*Service, error) {

	var s Service
	s.Urls = urls
	s.Spider = sp
	s.Storage = st

	var err error
	// Создаем индекс для сканирования
	s.index, err = index.New(s.Storage)
	if err != nil {
		return &s, err
	}

	return &s, nil
}

// Channel инициализирует канал для сообщений о готовности данных
func (s *Service) Channel(ch chan<- int) {
	s.channel = ch
}

// Build сканирует страницы, строит индекс
func (s *Service) Build() {
	for _, url := range s.Urls {
		fmt.Printf("\n[builder] Сканируем  %s...", url)
		data, err := s.Spider.Scan(url, 2)
		if err != nil {
			fmt.Printf("\n[builder] ошибка при сканировании сайта %s: %v", url, err)
			continue
		}
		fmt.Printf("\n[builder]  ...найдено %d документов, индексируем...", len(data))
		// Строим индекс по списку просканированных документов
		_, err = s.index.Build(data)
		if err != nil {
			fmt.Printf("[builder] %s", err)
			continue
		}
		// Сохраняем индекс в файл
		err = s.index.SaveData()
		if err != nil {
			fmt.Printf("[builder] %s", err)
			continue
		}
		if s.channel != nil {
			// Посылаем сигнал в канал о необходимости обновить текущий поисковый индекс
			s.channel <- 1
		}
		fmt.Printf("\n[builder] Проиндексировано %d страниц", s.index.Records.Count)
	}
}
