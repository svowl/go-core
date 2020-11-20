package engine

import (
	"go-core/7-lesson/pkg/crawler"
	"go-core/7-lesson/pkg/index"
)

// Пакет engine производит поиск строки в индексе и возвращает список документов

// Service это структура для хранения состояния поискового движка
type Service struct {
	index *index.Index
}

// New создает объект Service
func New() *Service {
	s := Service{}
	return &s
}

// Init инициализирует объект Service
func (s *Service) Init(i *index.Index) {
	s.index = i
}

// Search ищет фразу в индексе и возвращает список документов
func (s *Service) Search(phrase string) ([]crawler.Document, bool) {
	var result []crawler.Document
	var found bool = false
	var ids []int
	ids, found = s.index.Search(phrase)
	if found == false {
		// Фраза не найдена
		return nil, false
	}
	// Фраза найдена в хеше, ids содержит индексы документов (Record.ID) в массиве Records
	for _, id := range ids {
		// Поиск записей в Records по id (Record.ID)
		record := s.index.Records.Search(id)
		if record == nil {
			continue
		}
		result = append(result, record.Value.(crawler.Document))
		found = true
	}
	return result, found
}
